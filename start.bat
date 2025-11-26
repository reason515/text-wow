@echo off
echo ═══════════════════════════════════════════════════════════
echo   TEXT WoW - 放置类文字RPG 启动脚本
echo ═══════════════════════════════════════════════════════════
echo.

echo [1/2] 启动后端服务器...
cd /d %~dp0server
start "Text WoW Server" cmd /k "go run main.go"

echo [2/2] 启动前端开发服务器...
cd /d %~dp0client
start "Text WoW Client" cmd /k "npm run dev"

echo.
echo ═══════════════════════════════════════════════════════════
echo   启动完成！
echo   后端: http://localhost:8080
echo   前端: http://localhost:5173
echo ═══════════════════════════════════════════════════════════
echo.
echo 请在浏览器中访问 http://localhost:5173 开始游戏
echo.
pause


