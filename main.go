package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/service"
	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

var debugLevel = 0 // 全局调试级别

// 添加格式化字节的函数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func init() {
	// 处理命令行参数
	for _, arg := range os.Args[1:] {
		switch arg {
		case "debug":
			debugLevel = 1
		case "debug-1":
			debugLevel = 1
		case "debug-2":
			debugLevel = 2
		}
	}
}

func main() {
	// 设置调试级别
	service.SetDebugLevel(debugLevel)

	reader := bufio.NewReader(os.Stdin)
	commonService := service.NewCommon()
	authService := service.NewAuth()
	dashboardService := service.NewDashboard()
	nodeService := service.NewNode()

	for {
		commonService.ClearScreen()
		fmt.Println(utils.ColorText(utils.Bold+utils.Cyan, "\n欢迎使用OpenBMCLAPI系统!"))
		fmt.Println(utils.ColorText(utils.Yellow, "0. GitHub登录"))
		fmt.Println(utils.ColorText(utils.Green, "1. 查看用户信息"))
		fmt.Println(utils.ColorText(utils.Green, "2. 查看系统状态"))
		fmt.Println(utils.ColorText(utils.Green, "3. 查看节点列表"))
		fmt.Println(utils.ColorText(utils.Green, "4. 查看节点排行榜"))
		fmt.Println(utils.ColorText(utils.Green, "5. 打开管理面板"))
		fmt.Println(utils.ColorText(utils.Red, "6. 退出程序"))
		fmt.Print(utils.ColorText(utils.Purple, "请选择操作 (0-6): "))

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "0":
			authURL, err := authService.GetGithubAuthURL()
			if err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "获取授权URL失败: %v\n"), err)
				continue
			}

			if err := authService.OpenBrowser(authURL); err != nil {
				fmt.Printf(utils.ColorText(utils.Yellow, "无法自动打开浏览器，请手动访问以下链接：\n%s\n"), authURL)
			}

			fmt.Print(utils.ColorText(utils.Cyan, "\n请将授权完成后的回调URL粘贴到这里: "))
			callbackURL, _ := reader.ReadString('\n')
			callbackURL = strings.TrimSpace(callbackURL)

			if code := authService.ExtractCode(callbackURL); code != "" {
				fmt.Printf(utils.ColorText(utils.Green, "授权码: %s\n"), code)

				if err := authService.VerifyCode(code); err != nil {
					fmt.Println(utils.ColorText(utils.Red, fmt.Sprintf("❌ 验证失败: %v", err)))
				}
			} else {
				fmt.Println(utils.ColorText(utils.Red, "❌ 无法获取授权码，请重试"))
			}

			fmt.Print(utils.ColorText(utils.Yellow, "\n按回车键继续..."))
			reader.ReadString('\n')
		case "1":
			profile, err := authService.GetUserProfile()
			if err != nil {
				fmt.Println(utils.ColorText(utils.Red, fmt.Sprintf("获取用户信息失败: %v", err)))
			} else {
				fmt.Println(utils.ColorText(utils.Green, "\n✓ 获取用户信息成功"))

				// 尝试显示ASCII头像
				if ascii, err := utils.ImageToAscii(profile.Avatar, 40); err == nil {
					fmt.Println(utils.ColorText(utils.Blue, "\n头像预览:"))
					fmt.Println(utils.ColorText(utils.Blue, ascii))
				}

				fmt.Printf(utils.ColorText(utils.Cyan, "用户名: %s\n"), profile.Name)
				fmt.Printf(utils.ColorText(utils.Cyan, "GitHub ID: %s\n"), profile.Username)
				fmt.Printf(utils.ColorText(utils.Cyan, "头像URL: %s\n"), profile.Avatar)
				if profile.RawProfile.Bio != "" {
					fmt.Printf(utils.ColorText(utils.Cyan, "简介: %s\n"), profile.RawProfile.Bio)
				}
				if profile.RawProfile.Blog != "" {
					fmt.Printf(utils.ColorText(utils.Cyan, "博客: %s\n"), profile.RawProfile.Blog)
				}
			}
			fmt.Print(utils.ColorText(utils.Yellow, "\n按回车键继续..."))
			reader.ReadString('\n')
		case "2":
			dashboard, err := dashboardService.GetDashboard()
			if err != nil {
				fmt.Println(utils.ColorText(utils.Red, fmt.Sprintf("获取面板数据失败: %v", err)))
			} else {
				dashboardService.DisplayDashboard(dashboard)
			}
			fmt.Print(utils.ColorText(utils.Yellow, "\n按回车键继续..."))
			reader.ReadString('\n')
		case "3":
			nodes, err := nodeService.GetNodeList()
			if err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "获取节点列表失败: %v\n"), err)
				commonService.WaitForEnter()
				continue
			}
			nodeService.DisplayAndSelectNode(nodes)
		case "4":
			ranks, err := nodeService.GetNodeMetricRank(context.Background())
			if err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "获取排行榜失败: %v\n"), err)
				commonService.WaitForEnter()
				continue
			}
			showNodeRank(ranks)
			commonService.WaitForEnter()
		case "5":
			webService := service.NewWeb(8080)
			if err := webService.StartServer(); err != nil {
				fmt.Printf(utils.ColorText(utils.Red, "启动 Web 服务器失败: %v\n"), err)
				commonService.WaitForEnter()
				continue
			}
			fmt.Println(utils.ColorText(utils.Yellow, "\n按回车键关闭服务器..."))
			reader.ReadString('\n')
		case "6":
			fmt.Println(utils.ColorText(utils.Green, "感谢使用，再见！"))
			return
		default:
			fmt.Println(utils.ColorText(utils.Red, "无效的选择，请重试"))
			commonService.WaitForEnter()
		}
	}
}

func showNodeRank(ranks []service.NodeMetricRank) {
	fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📊 节点排行榜"))
	fmt.Println(strings.Repeat("─", 100))

	// 分页显示
	pageSize := 10 // 每页显示的数量
	totalPages := (len(ranks) + pageSize - 1) / pageSize
	currentPage := 0

	for {
		commonService := service.NewCommon()
		commonService.ClearScreen()
		fmt.Printf("\n%s\n", utils.ColorText(utils.Bold+utils.Blue, "📊 节点排行榜"))
		fmt.Println(strings.Repeat("─", 100))

		// 显示表头
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, utils.ColorText(utils.Bold, "序号\t节点名称\t请求数\t流量\t状态\t赞助商"))
		fmt.Fprintln(w, strings.Repeat("─", 100))

		// 计算当前页的起始和结束索引
		start := currentPage * pageSize
		end := start + pageSize
		if end > len(ranks) {
			end = len(ranks)
		}

		// 显示当前页的数据
		for i, rank := range ranks[start:end] {
			status := utils.ColorText(utils.Green, "在线")
			if !rank.IsEnabled {
				status = utils.ColorText(utils.Red, "离线")
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				start+i+1,
				utils.ColorText(utils.Cyan, rank.Name),
				utils.ColorText(utils.Yellow, fmt.Sprintf("%d", rank.Metric.Hits)),
				utils.ColorText(utils.Purple, formatBytes(rank.Metric.Bytes)),
				status,
				rank.Sponsor.Name,
			)
		}
		w.Flush()

		// 显示分页信息和操作提示
		fmt.Printf("\n%s\n", utils.ColorText(utils.Yellow, fmt.Sprintf("第 %d/%d 页 (共 %d 条记录)", currentPage+1, totalPages, len(ranks))))
		fmt.Println("\n操作说明:")
		fmt.Println(utils.ColorText(utils.Green, "n") + ": 下一页")
		fmt.Println(utils.ColorText(utils.Green, "p") + ": 上一页")
		fmt.Println(utils.ColorText(utils.Green, "q") + ": 返回主菜单")
		fmt.Print("\n请输入操作: ")

		var input string
		fmt.Scanln(&input)

		switch input {
		case "n":
			if currentPage < totalPages-1 {
				currentPage++
			}
		case "p":
			if currentPage > 0 {
				currentPage--
			}
		case "q":
			return
		}
	}
}
