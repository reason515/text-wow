# 启动和停止脚本使用说明

## 快速开始

### 在 PowerShell 中运行

```powershell
# 启动服务
.\start.bat

# 或者使用 PowerShell 脚本
.\start.ps1

# 停止服务
.\stop.bat

# 或者使用 PowerShell 脚本
.\stop.ps1
```

### 在 CMD 中运行

```cmd
# 启动服务
start.bat

# 停止服务
stop.bat
```

### 直接双击运行

- 双击 `start.bat` 启动服务
- 双击 `stop.bat` 停止服务

## 脚本说明

### start.bat / start.ps1
- ✅ 自动检测并关闭已运行的服务
- ✅ 检查端口占用情况
- ✅ 在新窗口中启动后端和前端服务
- ✅ 显示服务地址和状态

### stop.bat / stop.ps1
- ✅ 一键停止所有运行中的服务
- ✅ 清理端口占用

## 服务地址

- **后端服务**: http://localhost:8080
- **前端服务**: http://localhost:5173

## 注意事项

1. 首次运行前端服务可能需要安装依赖：`cd client && npm install`
2. 确保已安装 Go 和 Node.js
3. 如果端口被占用，脚本会显示警告，但不会强制关闭占用端口的进程
4. 关闭服务可以直接关闭对应的命令窗口，或运行 `stop.bat`




































