@echo off
chcp 65001 >nul
echo ═══════════════════════════════════════════════════════════
echo           TEXT-WOW 自动化测试套件
echo ═══════════════════════════════════════════════════════════
echo.

set SCRIPT_DIR=%~dp0

:: 参数处理
if "%1"=="" goto :run_all
if "%1"=="backend" goto :run_backend
if "%1"=="frontend" goto :run_frontend
if "%1"=="watch" goto :run_watch
if "%1"=="coverage" goto :run_coverage
if "%1"=="help" goto :show_help

:show_help
echo 用法: test.bat [选项]
echo.
echo 选项:
echo   (无参数)    运行所有测试
echo   backend     仅运行后端 Go 测试
echo   frontend    仅运行前端 Vue 测试
echo   watch       以监视模式运行前端测试
echo   coverage    运行测试并生成覆盖率报告
echo   help        显示此帮助信息
echo.
goto :eof

:run_all
echo [1/2] 运行后端测试...
echo ───────────────────────────────────────────────────────────
cd /d "%SCRIPT_DIR%server"
go test ./... -v
if errorlevel 1 (
    echo.
    echo ❌ 后端测试失败
    exit /b 1
)
echo.
echo ✓ 后端测试通过
echo.

echo [2/2] 运行前端测试...
echo ───────────────────────────────────────────────────────────
cd /d "%SCRIPT_DIR%client"
call npm run test:run
if errorlevel 1 (
    echo.
    echo ❌ 前端测试失败
    exit /b 1
)
echo.
echo ═══════════════════════════════════════════════════════════
echo ✓ 所有测试通过！
echo ═══════════════════════════════════════════════════════════
goto :eof

:run_backend
echo 运行后端 Go 测试...
echo ───────────────────────────────────────────────────────────
cd /d "%SCRIPT_DIR%server"
go test ./... -v
goto :eof

:run_frontend
echo 运行前端 Vue 测试...
echo ───────────────────────────────────────────────────────────
cd /d "%SCRIPT_DIR%client"
call npm run test:run
goto :eof

:run_watch
echo 以监视模式运行前端测试...
echo (按 Ctrl+C 退出)
echo ───────────────────────────────────────────────────────────
cd /d "%SCRIPT_DIR%client"
call npm run test
goto :eof

:run_coverage
echo 运行测试并生成覆盖率报告...
echo ───────────────────────────────────────────────────────────
echo.
echo [后端覆盖率]
cd /d "%SCRIPT_DIR%server"
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
echo 后端覆盖率报告已生成: server/coverage.html
echo.
echo [前端覆盖率]
cd /d "%SCRIPT_DIR%client"
call npm run test:coverage
echo 前端覆盖率报告已生成: client/coverage/index.html
echo.
echo ═══════════════════════════════════════════════════════════
echo 覆盖率报告生成完成！
echo ═══════════════════════════════════════════════════════════
goto :eof




