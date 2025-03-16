package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/models"
	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

type DashboardService struct{}

func NewDashboard() *DashboardService {
	return &DashboardService{}
}

func (s *DashboardService) GetDashboard() (*models.Dashboard, error) {
	client := utils.NewHTTPClient()
	respBody, err := client.DoGet("https://bd.bangbang93.com/openbmclapi/metric/dashboard", nil)
	if err != nil {
		return nil, err
	}

	var dashboard models.Dashboard
	if err := json.Unmarshal(respBody.Body, &dashboard); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return &dashboard, nil
}

// 修改带宽格式化函数
func formatBandwidth(bandwidth float64) string {
	// bandwidth 输入单位是 Mbps
	if bandwidth < 1000 {
		return fmt.Sprintf("%.2f Mbps", bandwidth)
	}
	return fmt.Sprintf("%.2f Gbps", bandwidth/1000)
}

func (s *DashboardService) DisplayDashboard(dashboard *models.Dashboard) {
	fmt.Println(utils.ColorText(utils.Bold+utils.Green, "\n=== OpenBMCLAPI 系统状态面板 ==="))

	boxWidth := 24

	// 第一行：关键指标
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📊 关键指标"))
	fmt.Printf("┌%s┐ ┌%s┐ ┌%s┐\n",
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth))

	// 标题行
	title1 := utils.PadString("当前在线节点数", boxWidth-2, true)
	title2 := utils.PadString("当前出网带宽", boxWidth-2, true)
	title3 := utils.PadString("系统负载", boxWidth-2, true)
	fmt.Printf("│ %s │ │ %s │ │ %s │\n",
		utils.ColorText(utils.Yellow, title1),
		utils.ColorText(utils.Yellow, title2),
		utils.ColorText(utils.Yellow, title3))

	// 数值行
	value1 := utils.PadString(fmt.Sprintf("%d 个", dashboard.CurrentNodes), boxWidth-2, false)
	value2 := utils.PadString(formatBandwidth(dashboard.CurrentBandwidth), boxWidth-2, false)
	value3 := utils.PadString(fmt.Sprintf("%.2f%%", dashboard.Load*100), boxWidth-2, false)
	fmt.Printf("│ %s │ │ %s │ │ %s │\n",
		utils.ColorText(utils.Cyan, value1),
		utils.ColorText(utils.Cyan, value2),
		utils.ColorText(utils.Cyan, value3))

	fmt.Printf("└%s┘ └%s┘ └%s┘\n",
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth))

	// 第二行：累计数据
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📈 累计数据"))
	fmt.Printf("┌%s┐ ┌%s┐ ┌%s┐\n",
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth))

	// 标题行
	title1 = utils.PadString("当日总流量", boxWidth-2, true)
	title2 = utils.PadString("当日请求次数", boxWidth-2, true)
	title3 = utils.PadString("带宽上限", boxWidth-2, true)
	fmt.Printf("│ %s │ │ %s │ │ %s │\n",
		utils.ColorText(utils.Yellow, title1),
		utils.ColorText(utils.Yellow, title2),
		utils.ColorText(utils.Yellow, title3))

	// 数值行
	value1 = utils.PadString(models.FormatBytes(dashboard.Bytes), boxWidth-2, false)
	value2 = utils.PadString(fmt.Sprintf("%d 次", dashboard.Hits), boxWidth-2, false)
	value3 = utils.PadString(formatBandwidth(dashboard.Bandwidth), boxWidth-2, false)
	fmt.Printf("│ %s │ │ %s │ │ %s │\n",
		utils.ColorText(utils.Cyan, value1),
		utils.ColorText(utils.Cyan, value2),
		utils.ColorText(utils.Cyan, value3))

	fmt.Printf("└%s┘ └%s┘ └%s┘\n",
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth),
		strings.Repeat("─", boxWidth))

	// 第三行：历史数据
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📅 最近24小时数据趋势"))
	s.displayHourlyChart(dashboard.Hourly)

	// 修改图表显示部分
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📊 数据趋势"))

	// 节点数趋势图
	fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, "每小时在线节点数 (个)"))
	s.displayLineChart(dashboard.Hourly, func(h models.HourlyMetric) float64 {
		return float64(h.Nodes)
	}, 120)

	// 带宽趋势图
	fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, "平均每小时出网带宽 (Gbps)"))
	s.displayLineChart(dashboard.Hourly, func(h models.HourlyMetric) float64 {
		return h.Bandwidth / 1000 // 转换为 Gbps
	}, 15)

	// 流量趋势图
	fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, "每小时流量分布 (TB)"))
	s.displayLineChart(dashboard.Hourly, func(h models.HourlyMetric) float64 {
		return float64(h.Bytes) / (1024 * 1024 * 1024 * 1024) // 转换为 TB
	}, 6)

	// 请求数趋势图
	fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, "每小时请求次数 (万次)"))
	s.displayLineChart(dashboard.Hourly, func(h models.HourlyMetric) float64 {
		return float64(h.Hits) / 10000
	}, 180)

	// 修改详细数据表格
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📋 最近6小时详细数据"))
	fmt.Printf("┌──────┬──────────┬────────────┬──────────────┬────────────┐\n")
	fmt.Printf("│ %-4s │ %-8s │ %-10s │ %-12s │ %-10s │\n",
		utils.ColorText(utils.Yellow, "时段"),
		utils.ColorText(utils.Yellow, "节点数"),
		utils.ColorText(utils.Yellow, "流量"),
		utils.ColorText(utils.Yellow, "请求次数"),
		utils.ColorText(utils.Yellow, "平均带宽"))
	fmt.Printf("├──────┼──────────┼────────────┼──────────────┼────────────┤\n")

	recentHours := dashboard.Hourly[len(dashboard.Hourly)-6:]
	// 反转顺序，从小时数小的开始显示
	for i := 0; i < len(recentHours); i++ {
		hour := recentHours[i]
		fmt.Printf("│ %2d时  │ %4d台   │ %8s │ %10d次  │ %8s │\n",
			hour.ID,
			hour.Nodes,
			models.FormatBytes(hour.Bytes),
			hour.Hits,
			formatBandwidth(hour.Bandwidth))
	}
	fmt.Printf("└──────┴──────────┴────────────┴──────────────┴────────────┘\n")
}

// 修改图表显示函数
func (s *DashboardService) displayLineChart(hourly []models.HourlyMetric, getValue func(models.HourlyMetric) float64, maxY float64) {
	height := 10
	width := 24 // 固定24小时

	// 重新排序数据，从1点到24点
	sortedHourly := make([]models.HourlyMetric, len(hourly))
	copy(sortedHourly, hourly)
	for i, j := 0, len(sortedHourly)-1; i < j; i, j = i+1, j-1 {
		sortedHourly[i], sortedHourly[j] = sortedHourly[j], sortedHourly[i]
	}

	// 创建图表数据
	chart := make([][]string, height)
	for i := range chart {
		chart[i] = make([]string, width)
		for j := range chart[i] {
			chart[i][j] = " "
		}
	}

	// 填充图表数据
	for x, metric := range sortedHourly {
		value := getValue(metric)
		y := int((value / maxY) * float64(height-1))
		if y >= height {
			y = height - 1
		}
		for i := 0; i <= y; i++ {
			chart[height-1-i][x] = "█"
		}
	}

	// 显示Y轴刻度和图表
	for i := 0; i < height; i++ {
		value := maxY * float64(height-i) / float64(height)
		fmt.Printf("%6.0f┤", value)
		for j := 0; j < width; j++ {
			if j < len(sortedHourly) {
				fmt.Print(utils.ColorText(utils.Blue, chart[i][j]))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	// 显示X轴
	fmt.Print("      └")
	for i := 0; i < width; i++ {
		fmt.Print("─")
	}
	fmt.Println()

	// 修改时间刻度显示
	fmt.Print("        ")
	for i := 0; i < width; i += 6 {
		if i < len(sortedHourly) {
			fmt.Printf("%2d时   ", sortedHourly[i].ID)
		}
	}
	fmt.Println()
}

// 显示小时数据的趋势图
func (s *DashboardService) displayHourlyChart(hourly []models.HourlyMetric) {
	// 重新排序数据
	sortedHourly := make([]models.HourlyMetric, len(hourly))
	copy(sortedHourly, hourly)
	for i, j := 0, len(sortedHourly)-1; i < j; i, j = i+1, j-1 {
		sortedHourly[i], sortedHourly[j] = sortedHourly[j], sortedHourly[i]
	}

	// 找出节点数的最大值，用于计算比例
	maxNodes := 0
	for _, h := range sortedHourly {
		if h.Nodes > maxNodes {
			maxNodes = h.Nodes
		}
	}

	// 显示节点数趋势
	fmt.Println(utils.ColorText(utils.Yellow, "节点数变化趋势:"))
	for i := 0; i < 3; i++ {
		fmt.Printf("%3d┤", (maxNodes*(3-i))/3)
		for _, h := range sortedHourly {
			height := (h.Nodes * 3) / maxNodes
			if height >= (3 - i) {
				fmt.Print(utils.ColorText(utils.Blue, "█"))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Print("  0┤")
	for range sortedHourly {
		fmt.Print("─")
	}
	fmt.Println()
	fmt.Print("   ")
	for i := 0; i < len(sortedHourly); i += 6 {
		fmt.Printf("%5d时", sortedHourly[i].ID)
	}
	fmt.Println()

	// 修改带宽趋势显示
	maxBandwidth := 0.0
	for _, h := range sortedHourly {
		if h.Bandwidth/1000 > maxBandwidth { // 转换为 Gbps
			maxBandwidth = h.Bandwidth / 1000
		}
	}

	fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, "带宽使用趋势 (Gbps):"))
	for i := 0; i < 3; i++ {
		scale := float64(3-i) / 3.0
		fmt.Printf("%7.1f┤", maxBandwidth*scale)
		for _, h := range sortedHourly {
			height := int((h.Bandwidth / 1000 * 3) / maxBandwidth)
			if height >= (3 - i) {
				fmt.Print(utils.ColorText(utils.Cyan, "█"))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Print("   0┤")
	for range sortedHourly {
		fmt.Print("─")
	}
	fmt.Println()
	fmt.Print("   ")
	for i := 0; i < len(sortedHourly); i += 6 {
		fmt.Printf("%5d时", sortedHourly[i].ID)
	}
	fmt.Println()
}

// 获取节点列表
func (s *DashboardService) GetNodeList() ([]models.Node, error) {
	// 读取 cookie
	cookieData, err := ioutil.ReadFile("cookie.json")
	if err != nil {
		return nil, fmt.Errorf("读取 cookie 失败: %v", err)
	}

	var cookies []models.Cookie
	if err := json.Unmarshal(cookieData, &cookies); err != nil {
		return nil, fmt.Errorf("解析 cookie 失败: %v", err)
	}

	// 创建请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://bd.bangbang93.com/openbmclapi/mgmt/cluster/my", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加 cookie
	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var nodes []models.Node
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return nodes, nil
}

// 显示节点列表
func (s *DashboardService) DisplayNodeList(nodes []models.Node) {
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📡 节点列表"))
	fmt.Printf("┌──────────────────────┬──────────┬──────────┬──────────┬─────────┬──────────────┐\n")
	fmt.Printf("│ %-20s │ %-8s │ %-8s │ %-8s │ %-7s │ %-10s │\n",
		utils.ColorText(utils.Yellow, "节点名称"),
		utils.ColorText(utils.Yellow, "状态"),
		utils.ColorText(utils.Yellow, "带宽"),
		utils.ColorText(utils.Yellow, "实测带宽"),
		utils.ColorText(utils.Yellow, "信任度"),
		utils.ColorText(utils.Yellow, "最后活动"))
	fmt.Printf("├──────────────────────┼──────────┼──────────┼──────────┼─────────┼──────────────┤\n")

	for _, node := range nodes {
		// 状态颜色
		statusColor := utils.Green
		status := "在线"
		if !node.IsEnabled {
			statusColor = utils.Red
			status = "离线"
		}

		// 信任度颜色
		trustColor := utils.Green
		if node.Trust < 0 {
			trustColor = utils.Red
		}

		// 格式化节点名称
		nodeName := utils.PadString(truncateString(node.Name, 20), 20, true)

		// 格式化状态
		statusStr := utils.PadString(status, 8, true)

		// 格式化带宽和实测带宽
		bandwidthStr := fmt.Sprintf("%8d", node.Bandwidth)
		measureBandwidthStr := fmt.Sprintf("%8d", node.MeasureBandwidth)

		// 格式化信任度
		trustStr := fmt.Sprintf("%7d", node.Trust)

		// 格式化时间
		timeStr := utils.PadString(node.LastActivity.Format("01-02 15:04"), 10, true)

		fmt.Printf("│ %s │ %s │ %s │ %s │ %s │ %s │\n",
			utils.ColorText(utils.Cyan, nodeName),
			utils.ColorText(statusColor, statusStr),
			bandwidthStr,
			measureBandwidthStr,
			utils.ColorText(trustColor, trustStr),
			timeStr)
	}
	fmt.Printf("└──────────────────────┴──────────┴──────────┴──────────┴─────────┴──────────────┘\n")
}

// 辅助函数：截断字符串，考虑中文字符
func truncateString(s string, length int) string {
	r := []rune(s)
	if len(r) > length {
		return string(r[:length-3]) + "..."
	}
	return s
}
