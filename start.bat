@echo off
chcp 65001 >nul
echo ============================================================
echo   TEXT WoW - 一键启动脚本
echo ============================================================
echo.

REM 检查并关闭已运行的服务
echo [检查] 检查已运行的服务...
tasklist /FI "WINDOWTITLE eq Text WoW Server*" 2>NUL | find /I /N "cmd.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [关闭] 发现已运行的后端服务，正在关闭...
    taskkill /FI "WINDOWTITLE eq Text WoW Server*" /F >NUL 2>&1
    timeout /t 1 /nobreak >NUL
)

tasklist /FI "WINDOWTITLE eq Text WoW Client*" 2>NUL | find /I /N "cmd.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [关闭] 发现已运行的前端服务，正在关闭...
    taskkill /FI "WINDOWTITLE eq Text WoW Client*" /F >NUL 2>&1
    timeout /t 1 /nobreak >NUL
)

REM 检查端口占用
netstat -ano | findstr ":8080" >NUL
if "%ERRORLEVEL%"=="0" (
    echo [警告] 端口 8080 已被占用，可能影响后端服务启动
)

netstat -ano | findstr ":5173" >NUL
if "%ERRORLEVEL%"=="0" (
    echo [警告] 端口 5173 已被占用，可能影响前端服务启动
)

echo.
echo [1/2] 启动后端服务器 (端口 8080)...
cd /d %~dp0server
start "Text WoW Server" cmd /k "title Text WoW Server && go run main.go"
timeout /t 2 /nobreak >NUL

echo [2/2] 启动前端开发服务器 (端口 5173)...
cd /d %~dp0client
start "Text WoW Client" cmd /k "title Text WoW Client && npm run dev"
timeout /t 2 /nobreak >NUL

echo.
echo ============================================================
echo   启动完成！
echo ============================================================
echo   后端服务: http://localhost:8080
echo   前端服务: http://localhost:5173
echo.
echo   请在浏览器中访问 http://localhost:5173 开始游戏
echo.
echo   提示: 关闭服务请直接关闭对应的命令窗口
echo ============================================================
echo.
pause
