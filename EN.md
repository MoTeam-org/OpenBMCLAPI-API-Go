# OBA-BD (OpenBMCLAPI Backend Dashboard)

<div align="center">
    <img src="logo.svg" alt="OBA-BD Logo" width="200" height="200">
</div>

OBA-BD is a command-line tool for managing and monitoring OpenBMCLAPI backend services. It provides a user-friendly interactive interface that allows administrators to conveniently view and manage various OpenBMCLAPI services.

## ✨ Features

- 🔐 GitHub Authentication
  - Browser-based OAuth login
  - Direct cookie paste login
- 👤 User Management
  - View user profile information
  - Display GitHub avatar (ASCII art)
  - User role management
- 📊 System Monitoring
  - Real-time system status
  - Performance metrics tracking
  - Resource usage monitoring
  - Alert configuration
- 🖥️ Node Management
  - Comprehensive node list view
  - Detailed node information
  - Node performance rankings
  - Node health monitoring
- 🌐 Web Dashboard
  - Modern web interface
  - Real-time data visualization
  - Interactive charts and graphs
  - Responsive design

## 🛠️ System Requirements

- Windows OS (Windows 10 recommended)
- Network connection
- 4GB RAM minimum
- 100MB free disk space

## 📦 Installation

1. Download the latest release `OBA-BD-V1.0.1.exe`
2. Run the executable by double-clicking or via command line
3. Complete GitHub authorization on first use

## 🚀 Quick Start

1. Launch the program
2. Select from the main menu:
   ```
   0. GitHub Login
   1. View User Profile
   2. View System Status
   3. View Node List
   4. View Node Rankings
   5. Open Web Dashboard
   6. Exit
   ```
3. Follow on-screen prompts for each module

## 🎨 Web Interface

The web dashboard provides a modern, responsive interface with:

- Vue 3 based modern UI
- Ant Design Vue components
- Dark/Light theme support
- Real-time data visualization with ECharts
- Responsive and mobile-friendly design

## 🔧 Debug Mode

Launch with debug flags:
```bash
# Basic debugging
./OBA-BD-V1.0.1.exe debug

# Advanced debugging
./OBA-BD-V1.0.1.exe debug-2
```

## 💻 Tech Stack

### Backend
- Go 1.19+
- Gin Web Framework
- GORM
- JWT Authentication

### Frontend
- Vue 3.4+
- TypeScript 5.2+
- Ant Design Vue 4.1+
- ECharts 5.5
- Axios
- Pinia (State Management)
- Vue Router 4

### Development Tools
- Vite 5
- TypeScript
- Vue TSC
- Node.js

### Infrastructure
- GitHub OAuth
- WebSocket
- SQLite
- Docker support

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 🔒 Security

### Vulnerability Scanning
- The project uses GitHub's Dependabot for automated security updates
- Regular security scans for dependencies
- Immediate patches for critical vulnerabilities
- Security advisories through GitHub Security tab

### Security Best Practices
- Regular dependency updates
- Code scanning enabled
- Secure authentication implementation
- Protected API endpoints

### Reporting Security Issues
If you discover a security vulnerability, please:
1. **DO NOT** open a public issue
2. Follow our [Security Policy](SECURITY.md)
3. Report it through GitHub's Security Advisory feature
4. Or email us directly at [security contact]

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details 