package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/models"
	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

type NodeService struct{}

func NewNode() *NodeService {
	return &NodeService{}
}

// GetNodeList 获取节点列表
func (s *NodeService) GetNodeList() ([]models.Node, error) {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return nil, fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return nil, fmt.Errorf("解析 cookie 失败: %v", err)
	}

	client := utils.NewHTTPClient()
	respBody, err := client.DoGet("https://bd.bangbang93.com/openbmclapi/mgmt/cluster/my", cookies)
	if err != nil {
		return nil, err
	}

	var nodes []models.Node
	if err := json.Unmarshal(respBody.Body, &nodes); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return nodes, nil
}

// GetNodeDetail 获取节点详情
func (s *NodeService) GetNodeDetail(nodeID string) (*models.Node, error) {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return nil, fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return nil, fmt.Errorf("解析 cookie 失败: %v", err)
	}

	url := fmt.Sprintf("https://bd.bangbang93.com/openbmclapi/mgmt/cluster/%s", nodeID)
	client := utils.NewHTTPClient()
	respBody, err := client.DoGet(url, cookies)
	if err != nil {
		return nil, err
	}

	var node models.Node
	if err := json.Unmarshal(respBody.Body, &node); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return &node, nil
}

// DisplayAndSelectNode 显示节点列表并处理选择
func (s *NodeService) DisplayAndSelectNode(nodes []models.Node) {
	reader := bufio.NewReader(os.Stdin)
	commonService := NewCommon()

	for {
		commonService.ClearScreen() // 每次显示列表前清屏
		fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📡 节点列表"))
		fmt.Println(strings.Repeat("─", 50))

		for i, node := range nodes {
			statusColor := utils.Green
			status := "在线"
			if !node.IsEnabled {
				statusColor = utils.Red
				status = "离线"
			}

			fmt.Printf("%d. %s [%s] (ID: %s)\n",
				i+1,
				utils.ColorText(utils.Cyan, node.Name),
				utils.ColorText(statusColor, status),
				node.ID)
		}

		fmt.Println(strings.Repeat("─", 50))
		fmt.Print(utils.ColorText(utils.Yellow, "请选择节点编号 (输入 q 返回): "))

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			return
		}

		var index int
		if _, err := fmt.Sscanf(input, "%d", &index); err != nil || index < 1 || index > len(nodes) {
			fmt.Println(utils.ColorText(utils.Red, "无效的选择，请重试"))
			commonService.WaitForEnter() // 使用通用的等待函数
			continue
		}

		selectedNode := nodes[index-1]
		fmt.Printf("选择的节点 ID: %s\n", selectedNode.ID)
		nodeDetail, err := s.GetNodeDetail(selectedNode.ID)
		if err != nil {
			fmt.Printf(utils.ColorText(utils.Red, "获取节点详情失败: %v\n"), err)
			commonService.WaitForEnter() // 使用通用的等待函数
			continue
		}

		s.DisplayNodeDetail(nodeDetail)
	}
}

// NodeUpdateInfo 节点更新信息
type NodeUpdateInfo struct {
	Name      string `json:"name"`
	Bandwidth int    `json:"bandwidth"`
}

// UpdateNode 更新节点信息
func (s *NodeService) UpdateNode(nodeID string, info NodeUpdateInfo) error {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return fmt.Errorf("解析 cookie 失败: %v", err)
	}

	url := fmt.Sprintf("https://bd.bangbang93.com/openbmclapi/mgmt/cluster/%s", nodeID)
	client := utils.NewHTTPClient()
	_, err = client.DoPatch(url, info, cookies)
	return err
}

// NodeSponsorUpdate 赞助商更新信息
type NodeSponsorUpdate struct {
	Sponsor models.NodeSponsor `json:"sponsor"`
}

// UpdateNodeSponsor 更新节点赞助商信息
func (s *NodeService) UpdateNodeSponsor(nodeID string, sponsor models.NodeSponsor) error {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return fmt.Errorf("解析 cookie 失败: %v", err)
	}

	updateInfo := NodeSponsorUpdate{
		Sponsor: sponsor,
	}

	url := fmt.Sprintf("https://bd.bangbang93.com/openbmclapi/mgmt/cluster/%s", nodeID)
	client := utils.NewHTTPClient()
	_, err = client.DoPatch(url, updateInfo, cookies)
	return err
}

// ResetNodeSecret 重置节点密钥
func (s *NodeService) ResetNodeSecret(nodeId string) (string, error) {
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return "", fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return "", fmt.Errorf("解析 cookie 失败: %v", err)
	}

	client := utils.NewHTTPClient()
	url := fmt.Sprintf("https://bd.bangbang93.com/openbmclapi/mgmt/cluster/%s/reset-secret", nodeId)
	respBody, err := client.DoPatch(url, nil, cookies)
	if err != nil {
		return "", err
	}

	var result struct {
		Secret string `json:"secret"`
	}

	if err := json.Unmarshal(respBody.Body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Secret, nil
}

// DisplayNodeDetail 显示节点详情
func (s *NodeService) DisplayNodeDetail(node *models.Node) {
	commonService := NewCommon()

	for {
		commonService.ClearScreen()
		s.showNodeDetail(node)

		fmt.Println(strings.Repeat("─", 50))
		fmt.Println(utils.ColorText(utils.Yellow, "操作选项:"))
		fmt.Println("1. 修改节点信息")
		fmt.Println("2. 修改赞助商信息")
		fmt.Println("3. 重置节点密钥")
		fmt.Println("4. 刷新节点信息")
		fmt.Println("q. 返回上级菜单")
		fmt.Print(utils.ColorText(utils.Purple, "请选择操作: "))

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if err := s.editNodeInfo(node); err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "修改失败: %v\n"), err)
				commonService.WaitForEnter()
			} else {
				fmt.Println(utils.ColorText(utils.Green, "修改成功!"))
				// 刷新节点信息
				updatedNode, err := s.GetNodeDetail(node.ID)
				if err == nil {
					node = updatedNode
				}
				commonService.WaitForEnter()
			}
		case "2":
			if err := s.editSponsorInfo(node); err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "修改失败: %v\n"), err)
				commonService.WaitForEnter()
			} else {
				fmt.Println(utils.ColorText(utils.Green, "修改成功!"))
				// 刷新节点信息
				updatedNode, err := s.GetNodeDetail(node.ID)
				if err == nil {
					node = updatedNode
				}
				commonService.WaitForEnter()
			}
		case "3":
			commonService.ClearScreen() // 重置密钥前清屏
			fmt.Print(utils.ColorText(utils.Red, "\n⚠️ 警告: 重置密钥将导致节点需要重新配置!\n"))
			fmt.Print(utils.ColorText(utils.Yellow, "确认重置? (输入 'RESET' 确认): "))
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)

			if confirm == "RESET" {
				if secret, err := s.ResetNodeSecret(node.ID); err != nil {
					fmt.Printf(utils.ColorText(utils.Red, "重置失败: %v\n"), err)
				} else {
					fmt.Printf(utils.ColorText(utils.Green, "重置成功!\n"))
					fmt.Printf(utils.ColorText(utils.Yellow, "新密钥: %s\n"), secret)
				}
			} else {
				fmt.Println(utils.ColorText(utils.Yellow, "已取消重置"))
			}
			commonService.WaitForEnter()
		case "4":
			updatedNode, err := s.GetNodeDetail(node.ID)
			if err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "刷新失败: %v\n"), err)
			} else {
				node = updatedNode
				fmt.Println(utils.ColorText(utils.Green, "刷新成功!"))
			}
			commonService.WaitForEnter()
		case "q":
			return
		default:
			fmt.Println(utils.ColorText(utils.Red, "无效的选择"))
			commonService.WaitForEnter()
		}
	}
}

// 编辑节点信息
func (s *NodeService) editNodeInfo(node *models.Node) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📝 编辑节点信息"))
	fmt.Println(strings.Repeat("─", 50))

	// 显示当前值并获取新值
	fmt.Printf("节点名称 (当前: %s): ", node.Name)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)
	if newName == "" {
		newName = node.Name
	}

	fmt.Printf("带宽限制 (当前: %d Mbps): ", node.Bandwidth)
	bandwidthStr, _ := reader.ReadString('\n')
	bandwidthStr = strings.TrimSpace(bandwidthStr)
	newBandwidth := node.Bandwidth
	if bandwidthStr != "" {
		if bw, err := strconv.Atoi(bandwidthStr); err == nil {
			newBandwidth = bw
		}
	}

	// 确认修改
	fmt.Print(utils.ColorText(utils.Yellow, "\n确认修改? (y/N): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm == "y" {
		updateInfo := NodeUpdateInfo{
			Name:      newName,
			Bandwidth: newBandwidth,
		}
		return s.UpdateNode(node.ID, updateInfo)
	}

	return nil
}

// 编辑赞助商信息
func (s *NodeService) editSponsorInfo(node *models.Node) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📝 编辑赞助商信息"))
	fmt.Println(strings.Repeat("─", 50))

	// 显示当前值并获取新值
	fmt.Printf("赞助商名称 (当前: %s): ", node.Sponsor.Name)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)
	if newName == "" {
		newName = node.Sponsor.Name
	}

	fmt.Printf("赞助商网址 (当前: %s): ", node.Sponsor.URL)
	newURL, _ := reader.ReadString('\n')
	newURL = strings.TrimSpace(newURL)
	if newURL == "" {
		newURL = node.Sponsor.URL
	}

	fmt.Printf("赞助商图片 (当前: %s): ", node.Sponsor.Banner)
	newBanner, _ := reader.ReadString('\n')
	newBanner = strings.TrimSpace(newBanner)
	if newBanner == "" {
		newBanner = node.Sponsor.Banner
	}

	// 确认修改
	fmt.Print(utils.ColorText(utils.Yellow, "\n确认修改? (y/N): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm == "y" {
		sponsor := models.NodeSponsor{
			Name:   newName,
			URL:    newURL,
			Banner: newBanner,
		}
		if err := s.UpdateNodeSponsor(node.ID, sponsor); err != nil {
			return err
		}
		fmt.Println(utils.ColorText(utils.Yellow, "\n提示: 赞助商信息的修改需要管理员审核后才会生效"))
	}

	return nil
}

// 显示节点详情（原来的显示逻辑）
func (s *NodeService) showNodeDetail(node *models.Node) {
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📝 节点详情"))
	fmt.Println(strings.Repeat("─", 50))

	// 基本信息
	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "节点名称:"),
		utils.ColorText(utils.Cyan, node.Name))

	// 运行状态
	statusColor := utils.Green
	status := "在线"
	if !node.IsEnabled {
		statusColor = utils.Red
		status = "离线"
		if node.DownReason != "" {
			status = fmt.Sprintf("离线 (%s)", node.DownReason)
		}
	}
	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "运行状态:"),
		utils.ColorText(statusColor, status))

	// Ban 状态
	if node.IsBanned {
		fmt.Printf("%s %s\n",
			utils.ColorText(utils.Yellow, "封禁状态:"),
			utils.ColorText(utils.Red, fmt.Sprintf("已封禁 (%s)", node.BanReason)))
	}

	// 其他信息保持不变...
	fmt.Printf("%s %d Mbps\n",
		utils.ColorText(utils.Yellow, "设置带宽:"),
		node.Bandwidth)

	fmt.Printf("%s %d Mbps\n",
		utils.ColorText(utils.Yellow, "实测带宽:"),
		node.MeasureBandwidth)

	// 运行信息
	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "运行环境:"),
		node.Flavor.Runtime)

	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "存储类型:"),
		node.Flavor.Storage)

	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "程序版本:"),
		node.Version)

	// 节点状态
	trustColor := utils.Green
	if node.Trust < 0 {
		trustColor = utils.Red
	}
	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "信任度:"),
		utils.ColorText(trustColor, fmt.Sprintf("%d", node.Trust)))

	// 时间信息
	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "创建时间:"),
		node.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("%s %s\n",
		utils.ColorText(utils.Yellow, "最后活动:"),
		node.LastActivity.Format("2006-01-02 15:04:05"))

	if !node.Uptime.IsZero() {
		fmt.Printf("%s %s\n",
			utils.ColorText(utils.Yellow, "上线时间:"),
			node.Uptime.Format("2006-01-02 15:04:05"))
	}

	if !node.Downtime.IsZero() {
		fmt.Printf("%s %s\n",
			utils.ColorText(utils.Yellow, "下线时间:"),
			node.Downtime.Format("2006-01-02 15:04:05"))
	}

	// 赞助商信息
	if node.Sponsor.Name != "" {
		fmt.Println(strings.Repeat("─", 50))
		fmt.Printf("%s %s\n",
			utils.ColorText(utils.Yellow, "赞助商:"),
			node.Sponsor.Name)
		fmt.Printf("%s %s\n",
			utils.ColorText(utils.Yellow, "赞助商网站:"),
			node.Sponsor.URL)
	}

	// 节点地址
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("%s %s://%s:%d\n",
		utils.ColorText(utils.Yellow, "节点地址:"),
		node.Endpoint.Proto,
		node.Endpoint.Host,
		node.Endpoint.Port)
}

// NodeMetricRank 节点排行榜数据结构
type NodeMetricRank struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	FullSize  bool   `json:"fullSize,omitempty"`
	IsEnabled bool   `json:"isEnabled"`
	User      *struct {
		Name string `json:"name"`
	} `json:"user,omitempty"`
	Version      string    `json:"version,omitempty"`
	LastActivity time.Time `json:"lastActivity,omitempty"`
	DownReason   string    `json:"downReason,omitempty"`
	DownTime     time.Time `json:"downtime,omitempty"`
	Sponsor      struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Banner string `json:"banner"`
	} `json:"sponsor"`
	Metric struct {
		ID        string    `json:"_id"`
		ClusterID string    `json:"clusterId"`
		Date      time.Time `json:"date"`
		Version   int       `json:"__v"`
		Bytes     int64     `json:"bytes"`
		Hits      int64     `json:"hits"`
	} `json:"metric"`
}

// GetNodeMetricRank 获取节点排行榜
func (s *NodeService) GetNodeMetricRank(ctx context.Context) ([]NodeMetricRank, error) {
	// 发起 HTTP 请求获取数据
	url := "https://bd.bangbang93.com/openbmclapi/metric/rank"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 如果有登录态，添加认证头
	if token := ctx.Value("token"); token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var ranks []NodeMetricRank
	if err := json.NewDecoder(resp.Body).Decode(&ranks); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return ranks, nil
}

// 修改响应处理函数
func (s *NodeService) HandleGetNodeMetricRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		resp := models.ResponseError(http.StatusMethodNotAllowed, "Method not allowed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Code)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ranks, err := s.GetNodeMetricRank(r.Context())
	if err != nil {
		resp := models.ResponseError(http.StatusInternalServerError, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Code)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.ResponseSuccess(ranks)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 在 RegisterRoutes 方法中添加新的路由
func (s *NodeService) RegisterRoutes() {
	// ... 其他路由 ...
	http.HandleFunc("/api/nodes/rank", s.HandleGetNodeMetricRank)
}
