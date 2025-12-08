@echo off
chcp 65001 >nul
echo ============================================================
echo   TEXT WoW - 停止服务脚本
echo ============================================================
echo.

echo [停止] 正在停止后端服务...
taskkill /FI "WINDOWTITLE eq Text WoW Server*" /F >NUL 2>&1
if "%ERRORLEVEL%"=="0" (
    echo [成功] 后端服务已停止
) else (
    echo [信息] 未发现运行中的后端服务
)

echo [停止] 正在停止前端服务...
taskkill /FI "WINDOWTITLE eq Text WoW Client*" /F >NUL 2>&1
if "%ERRORLEVEL%"=="0" (
    echo [成功] 前端服务已停止
) else (
    echo [信息] 未发现运行中的前端服务
)

echo.
echo ============================================================
echo   所有服务已停止
echo ============================================================
echo.
pause
























