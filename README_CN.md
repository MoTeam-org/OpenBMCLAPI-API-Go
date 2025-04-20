# OBA-BD (OpenBMCLAPI 后端管理工具)

<div align="center">
    <img src="logo.svg" alt="OBA-BD Logo" width="200" height="200">
</div>

OBA-BD 是一个用于管理和监控 OpenBMCLAPI 后端服务的命令行工具。该工具提供了友好的交互式界面，让管理员能够方便地查看和管理 OpenBMCLAPI 的各项服务。

## ✨ 功能特性

- 🔐 GitHub 账号登录
  - 支持浏览器授权登录
  - 支持直接粘贴 Cookie 登录
  - 多重身份验证支持
- 👤 用户管理
  - 查看用户基本信息
  - 显示用户 GitHub 头像（ASCII 艺术）
  - 用户角色权限管理
- 📊 系统监控
  - 实时系统状态展示
  - 性能指标追踪
  - 资源使用监控
  - 告警配置管理
- 🖥️ 节点管理
  - 完整节点列表视图
  - 节点详细信息展示
  - 节点性能排行榜
  - 节点健康监控
- 🌐 Web 管理面板
  - 现代化 Web 界面
  - 实时数据可视化
  - 交互式图表
  - 响应式设计

## 🛠️ 系统要求

- Windows 操作系统（推荐 Windows 10）
- 网络连接
- 最小 4GB 内存
- 100MB 可用磁盘空间

## 📦 安装说明

1. 下载最新版本可执行文件 `OBA-BD-V1.0.1.exe`
2. 双击运行程序或通过命令行启动
3. 首次使用需完成 GitHub 授权

## 🚀 快速开始

1. 启动程序
2. 从主菜单选择功能：
   ```
   0. GitHub 登录
   1. 查看用户信息
   2. 查看系统状态
   3. 查看节点列表
   4. 查看节点排行
   5. 打开管理面板
   6. 退出程序
   ```
3. 按照屏幕提示操作

## 🎨 Web 界面

Web 管理面板提供现代化的响应式界面：

- 基于 Vue 3 的现代化界面
- Ant Design Vue 组件库
- 深色/浅色主题支持
- 基于 ECharts 的实时数据可视化
- 响应式设计，支持移动端

## 🔧 调试模式

通过启动参数开启调试：
```bash
# 基础调试
./OBA-BD-V1.0.1.exe debug

# 高级调试
./OBA-BD-V1.0.1.exe debug-2
```

## 💻 技术栈

### 后端
- Go 1.19+
- Gin Web 框架
- GORM
- JWT 认证

### 前端
- Vue 3.4+
- TypeScript 5.2+
- Ant Design Vue 4.1+
- ECharts 5.5
- Axios
- Pinia (状态管理)
- Vue Router 4

### 开发工具
- Vite 5
- TypeScript
- Vue TSC
- Node.js

### 基础设施
- GitHub OAuth
- WebSocket
- SQLite
- Docker 支持

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 🔒 安全性

### 漏洞扫描
- 项目使用 GitHub Dependabot 进行自动安全更新
- 定期进行依赖项安全扫描
- 及时修补关键漏洞
- 通过 GitHub Security 页面发布安全公告

### 安全最佳实践
- 定期更新依赖项
- 启用代码扫描
- 安全的身份验证实现
- 受保护的 API 端点

### 报告安全问题
如果您发现安全漏洞，请：
1. **不要**创建公开的 issue
2. 遵循我们的[安全政策](SECURITY.md)
3. 通过 GitHub 的安全公告功能报告
4. 或直接发送邮件至[安全联系人]

## 📄 许可证

MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情 