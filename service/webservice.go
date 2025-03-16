package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/models"
	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

//go:embed web/dist
var webContent embed.FS

type WebService struct {
	port int
}

func NewWeb(port int) *WebService {
	return &WebService{
		port: port,
	}
}

// StartServer 启动 Web 服务器
func (s *WebService) StartServer() error {
	// API 路由
	http.HandleFunc("/api/nodes", s.handleGetNodes)
	http.HandleFunc("/api/dashboard", s.handleGetDashboard)
	http.HandleFunc("/api/user", s.handleGetUser)
	http.HandleFunc("/api/nodes/rank", s.handleGetNodeRank)
	http.HandleFunc("/api/nodes/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/reset-secret"):
			s.handleResetSecret(w, r)
		case strings.HasSuffix(r.URL.Path, "/sponsor"):
			s.handleUpdateNode(w, r)
		default:
			s.handleUpdateNode(w, r)
		}
	})

	// 静态文件服务
	fsys, err := fs.Sub(webContent, "web/dist")
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(http.FS(fsys)))

	// 启动服务器
	serverURL := fmt.Sprintf("http://localhost:%d", s.port)
	fmt.Printf("Web 服务器已启动: %s\n", serverURL)

	// 在新的 goroutine 中启动服务器
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
			fmt.Printf("服务器启动失败: %v\n", err)
		}
	}()

	// 自动打开浏览器
	time.Sleep(100 * time.Millisecond)
	if err := s.openBrowser(serverURL); err != nil {
		fmt.Printf("无法自动打开浏览器，请手动访问: %s\n", serverURL)
	}

	return nil
}

// openBrowser 打开浏览器
func (s *WebService) openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	return err
}

// 修改 wrapResponse 函数
func wrapResponse(w http.ResponseWriter, code int, msg string, data interface{}) {
	var resp models.APIResponse
	if code == 200 {
		resp = models.ResponseSuccess(data)
	} else {
		resp = models.ResponseError(code, msg)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// API 处理函数
func (s *WebService) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog(1, "[Web API] GET /api/nodes - 获取节点列表")
	start := time.Now()

	nodeService := NewNode()
	nodes, err := nodeService.GetNodeList()
	if err != nil {
		utils.DebugLog(1, "[Web API] 获取节点列表失败: %v", err)
		wrapResponse(w, 500, err.Error(), nil)
		return
	}

	duration := time.Since(start)
	utils.DebugLog(1, "[Web API] 成功获取 %d 个节点，耗时: %v", len(nodes), duration)
	wrapResponse(w, 200, "success", nodes)
}

func (s *WebService) handleGetDashboard(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog(1, "[Web API] GET /api/dashboard - 获取仪表盘数据")
	start := time.Now()

	dashboardService := NewDashboard()
	dashboard, err := dashboardService.GetDashboard()
	if err != nil {
		utils.DebugLog(1, "[Web API] 获取仪表盘数据失败: %v", err)
		wrapResponse(w, 500, err.Error(), nil)
		return
	}

	duration := time.Since(start)
	utils.DebugLog(1, "[Web API] 成功获取仪表盘数据，耗时: %v", duration)
	wrapResponse(w, 200, "success", dashboard)
}

func (s *WebService) handleGetUser(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog(1, "[Web API] GET /api/user - 获取用户信息")
	start := time.Now()

	authService := NewAuth()
	user, err := authService.GetUserProfile()
	if err != nil {
		utils.DebugLog(1, "[Web API] 获取用户信息失败: %v", err)
		wrapResponse(w, 500, err.Error(), nil)
		return
	}

	duration := time.Since(start)
	utils.DebugLog(1, "[Web API] 成功获取用户信息，耗时: %v", duration)
	wrapResponse(w, 200, "success", user)
}

// 添加新的 API 处理函数
func (s *WebService) handleUpdateNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		wrapResponse(w, 405, "Method not allowed", nil)
		return
	}

	nodeID := strings.TrimPrefix(r.URL.Path, "/api/nodes/")
	nodeID = strings.TrimSuffix(nodeID, "/sponsor")

	var updateData struct {
		Name      string             `json:"name,omitempty"`
		Bandwidth int                `json:"bandwidth,omitempty"`
		Sponsor   models.NodeSponsor `json:"sponsor,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		wrapResponse(w, 400, "Invalid request data", nil)
		return
	}

	nodeService := NewNode()
	if strings.HasSuffix(r.URL.Path, "/sponsor") {
		// 更新赞助商信息
		err := nodeService.UpdateNodeSponsor(nodeID, updateData.Sponsor)
		if err != nil {
			wrapResponse(w, 500, err.Error(), nil)
			return
		}
	} else {
		// 更新节点基本信息
		err := nodeService.UpdateNode(nodeID, NodeUpdateInfo{
			Name:      updateData.Name,
			Bandwidth: updateData.Bandwidth,
		})
		if err != nil {
			wrapResponse(w, 500, err.Error(), nil)
			return
		}
	}

	wrapResponse(w, 200, "success", nil)
}

func (s *WebService) handleResetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		wrapResponse(w, 405, "Method not allowed", nil)
		return
	}

	nodeID := strings.TrimPrefix(r.URL.Path, "/api/nodes/")
	nodeID = strings.TrimSuffix(nodeID, "/reset-secret")

	nodeService := NewNode()
	secret, err := nodeService.ResetNodeSecret(nodeID)
	if err != nil {
		wrapResponse(w, 500, err.Error(), nil)
		return
	}

	wrapResponse(w, 200, "success", map[string]string{"secret": secret})
}

// 添加排行榜处理函数
func (s *WebService) handleGetNodeRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		wrapResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	nodeService := NewNode()
	ranks, err := nodeService.GetNodeMetricRank(r.Context())
	if err != nil {
		wrapResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	wrapResponse(w, http.StatusOK, "success", ranks)
}
