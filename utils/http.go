package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/models"
)

// RequestStatus 定义请求状态
type RequestStatus int

const (
	Preparing RequestStatus = iota
	Requesting
	Overtime
	Timeout
)

func (s RequestStatus) String() string {
	switch s {
	case Preparing:
		return "准备请求中"
	case Requesting:
		return "正在请求中"
	case Overtime:
		return "请求超过预期"
	case Timeout:
		return "请求超时"
	default:
		return "未知状态"
	}
}

// HTTPResponse 封装 HTTP 响应
type HTTPResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// HTTPClient 封装 HTTP 请求工具
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

// showProgress 显示请求进度
func showProgress(done chan bool, requestInfo string) {
	status := Preparing
	startTime := time.Now()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	fmt.Printf("\n%s\n", ColorText(Blue, requestInfo))

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			duration := time.Since(startTime)
			fmt.Print("\r")

			// 更新状态
			if duration.Seconds() < 1 {
				status = Preparing
			} else if duration.Seconds() < 20 {
				status = Requesting
			} else {
				status = Overtime
			}

			// 显示状态和耗时
			statusColor := Yellow
			switch status {
			case Preparing:
				statusColor = Cyan
			case Requesting:
				statusColor = Yellow
			case Overtime, Timeout:
				statusColor = Red
			}

			fmt.Print(ColorText(statusColor, fmt.Sprintf("[%s] ", status)))
			fmt.Printf(ColorText(Blue, "请求耗时：%.1f秒"), duration.Seconds())

			// 只在正常请求阶段显示预估时间
			if status == Requesting {
				remaining := 20.0 - duration.Seconds()
				if remaining > 0 {
					fmt.Printf(ColorText(Purple, " 预估剩余：%.1f秒"), remaining)
				}
			}

			// 根据状态显示不同的附加信息
			switch status {
			case Preparing:
				fmt.Print(ColorText(Cyan, " 正在初始化..."))
			case Overtime:
				fmt.Print(ColorText(Red, " 请求时间已超过预期，但仍在继续..."))
			case Timeout:
				fmt.Print(ColorText(Red, " 请求已超时"))
			}
		}
	}
}

// doRequest 执行 HTTP 请求
func (c *HTTPClient) doRequest(method, url string, body interface{}, cookies []models.Cookie) (*HTTPResponse, error) {
	// 准备请求体
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %v", err)
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "OpenBMCLAPI-Client/1.0")

	// 添加 cookie
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
			if cookie.Name == "XSRF-TOKEN" {
				req.Header.Set("X-XSRF-TOKEN", cookie.Value)
			}
		}
	}

	// 显示请求信息
	DebugLog(1, "[HTTP] %s %s", method, getEndpointDescription(url))
	if body != nil && DebugLevel >= 2 {
		DebugLog(2, "[HTTP] 请求数据:\n%s", JsonPretty(body))
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// 发送请求
	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("请求超时")
		}
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 计算请求耗时
	duration := time.Since(start)

	// 显示响应信息
	statusColor := Green
	if resp.StatusCode >= 400 {
		statusColor = Red
	} else if resp.StatusCode >= 300 {
		statusColor = Yellow
	}

	DebugLog(1, "[HTTP] %s %s %s (耗时: %v)",
		method,
		getEndpointDescription(url),
		ColorText(statusColor, fmt.Sprintf("[%d]", resp.StatusCode)),
		duration)

	// 如果是错误响应，尝试解析错误信息
	if resp.StatusCode >= 400 {
		var errResp struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Msg != "" {
			DebugLog(1, "[HTTP] 错误信息: %s", errResp.Msg)
		}
	}

	// 在调试级别2时显示详细响应
	if DebugLevel >= 2 {
		DebugLog(2, "[HTTP] 响应头:\n%s", formatHeaders(resp.Header))
		if len(respBody) > 0 {
			DebugLog(2, "[HTTP] 响应数据:\n%s", JsonPretty(json.RawMessage(respBody)))
		}
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}, nil
}

// formatHeaders 格式化 HTTP 头
func formatHeaders(headers http.Header) string {
	var sb strings.Builder
	for key, values := range headers {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}
	return sb.String()
}

// getEndpointDescription 获取接口描述
func getEndpointDescription(url string) string {
	switch {
	case strings.Contains(url, "/user"):
		return "获取用户信息"
	case strings.Contains(url, "/metric/dashboard"):
		return "获取仪表盘数据"
	case strings.Contains(url, "/mgmt/cluster/my"):
		return "获取节点列表"
	case strings.Contains(url, "/reset-secret"):
		return "重置节点密钥"
	case strings.Contains(url, "/mgmt/cluster/"):
		if strings.Contains(url, "/sponsor") {
			return "更新节点赞助商信息"
		}
		return "节点管理"
	default:
		return url
	}
}

// DoGet 执行 GET 请求
func (c *HTTPClient) DoGet(url string, cookies []models.Cookie) (*HTTPResponse, error) {
	return c.doRequest("GET", url, nil, cookies)
}

// DoPost 执行 POST 请求
func (c *HTTPClient) DoPost(url string, body interface{}, cookies []models.Cookie) (*HTTPResponse, error) {
	return c.doRequest("POST", url, body, cookies)
}

// DoPatch 执行 PATCH 请求
func (c *HTTPClient) DoPatch(url string, body interface{}, cookies []models.Cookie) (*HTTPResponse, error) {
	return c.doRequest("PATCH", url, body, cookies)
}

// JsonPretty 格式化 JSON 输出
func JsonPretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(b)
}
