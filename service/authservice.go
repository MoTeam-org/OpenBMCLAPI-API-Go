package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/models"
	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

type AuthService struct{}

func NewAuth() *AuthService {
	return &AuthService{}
}

// 定义请求状态
type requestStatus int

const (
	preparing requestStatus = iota
	requesting
	overtime
	timeout
)

func (s requestStatus) String() string {
	switch s {
	case preparing:
		return "准备请求中"
	case requesting:
		return "正在请求中"
	case overtime:
		return "请求超过预期"
	case timeout:
		return "请求超时"
	default:
		return "未知状态"
	}
}

func (s *AuthService) GetGithubAuthURL() (string, error) {
	client := utils.NewHTTPClient()
	respBody, err := client.DoGet("https://bd.bangbang93.com/openbmclapi/user/auth/github", nil)
	if err != nil {
		return "", err
	}

	location := respBody.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf(utils.ColorText(utils.Red, "未找到重定向URL"))
	}

	// 检查并补全URL
	if !strings.Contains(location, "client_id") {
		return "", fmt.Errorf(utils.ColorText(utils.Red, "获取到的授权URL不完整"))
	}

	fmt.Println()
	fmt.Println(utils.ColorText(utils.Green, "✓ 获取授权地址成功"))
	fmt.Printf(utils.ColorText(utils.Blue, "授权地址: %s\n"), location)
	fmt.Println(utils.ColorText(utils.Yellow, "正在准备 GitHub 授权页面..."))
	return location, nil
}

func (s *AuthService) OpenBrowser(url string) error {
	// 验证URL是否包含必要的参数
	if !strings.Contains(url, "client_id") || !strings.Contains(url, "redirect_uri") {
		return fmt.Errorf(utils.ColorText(utils.Red, "无效的授权URL"))
	}

	fmt.Println(utils.ColorText(utils.Green, "✓ 正在打开浏览器"))
	fmt.Println(utils.ColorText(utils.Yellow, "请在浏览器中完成 GitHub 授权..."))

	var err error
	switch runtime.GOOS {
	case "windows":
		// 使用 rundll32 来打开 URL
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		fmt.Printf(utils.ColorText(utils.Yellow, "无法自动打开浏览器，请手动访问以下链接：\n%s\n"), url)
	}
	return err
}

func (s *AuthService) ExtractCode(callbackURL string) string {
	code := s.parseCode(callbackURL)
	if code != "" {
		fmt.Println(utils.ColorText(utils.Green, "✓ 成功获取授权码"))
	}
	return code
}

// 新增一个辅助函数来解析code
func (s *AuthService) parseCode(callbackURL string) string {
	const prefix = "code="
	if idx := strings.Index(callbackURL, prefix); idx != -1 {
		code := callbackURL[idx+len(prefix):]
		if endIdx := strings.Index(code, "&"); endIdx != -1 {
			code = code[:endIdx]
		}
		return code
	}
	return ""
}

func (s *AuthService) VerifyCode(code string) error {
	url := fmt.Sprintf("https://bd.bangbang93.com/openbmclapi/user/auth/github?code=%s", code)
	client := utils.NewHTTPClient()
	respBody, err := client.DoGet(url, nil)
	if err != nil {
		return err
	}

	// 获取所有 Set-Cookie 头
	cookies := respBody.Header["Set-Cookie"]
	if len(cookies) == 0 {
		return fmt.Errorf("未找到 Cookie")
	}

	// 解析并存储 cookies
	var cookieList []models.Cookie
	for _, cookieStr := range cookies {
		cookie := parseCookie(cookieStr)
		if cookie != nil {
			cookieList = append(cookieList, *cookie)
		}
	}

	// 确保目录存在
	dir := filepath.Dir("cookie.json")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 将 cookies 写入文件
	cookieData, err := json.MarshalIndent(cookieList, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 Cookie 失败: %v", err)
	}

	if err := ioutil.WriteFile("cookie.json", cookieData, 0644); err != nil {
		return fmt.Errorf("保存 Cookie 失败: %v", err)
	}

	fmt.Println() // 清除进度显示的行
	fmt.Println(utils.ColorText(utils.Green, "✓ Cookie 已保存"))
	return nil
}

// 辅助函数：解析 Cookie 字符串
func parseCookie(cookieStr string) *models.Cookie {
	parts := strings.Split(cookieStr, ";")
	if len(parts) == 0 {
		return nil
	}

	// 解析主要的 name=value 部分
	nameValue := strings.SplitN(strings.TrimSpace(parts[0]), "=", 2)
	if len(nameValue) != 2 {
		return nil
	}

	cookie := &models.Cookie{
		Name:  nameValue[0],
		Value: nameValue[1],
	}

	// 解析其他属性
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		switch {
		case strings.EqualFold(part, "HttpOnly"):
			cookie.HttpOnly = true
		case strings.EqualFold(part, "Secure"):
			cookie.Secure = true
		case strings.HasPrefix(strings.ToLower(part), "path="):
			cookie.Path = strings.TrimPrefix(strings.ToLower(part), "path=")
		case strings.HasPrefix(strings.ToLower(part), "domain="):
			cookie.Domain = strings.TrimPrefix(strings.ToLower(part), "domain=")
		}
	}

	return cookie
}

func (s *AuthService) GetUserProfile() (*models.UserProfile, error) {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return nil, fmt.Errorf("读取Cookie失败，请先登录: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return nil, fmt.Errorf("解析Cookie失败: %v", err)
	}

	client := utils.NewHTTPClient()
	respBody, err := client.DoGet("https://bd.bangbang93.com/openbmclapi/user", cookies)
	if err != nil {
		return nil, err
	}

	var profile models.UserProfile
	if err := json.Unmarshal(respBody.Body, &profile); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &profile, nil
}
