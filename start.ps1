# TEXT WoW - One-click startup script (PowerShell)
# Usage: .\start.ps1

$ErrorActionPreference = "SilentlyContinue"

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  TEXT WoW - Startup Script" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "[Check] Checking running services..." -ForegroundColor Yellow

$serverProcesses = Get-Process | Where-Object { 
    $_.MainWindowTitle -like "*Text WoW Server*"
}
if ($serverProcesses) {
    Write-Host "[Stop] Found running backend service, stopping..." -ForegroundColor Yellow
    $serverProcesses | Stop-Process -Force
    Start-Sleep -Seconds 1
}

$clientProcesses = Get-Process | Where-Object { 
    $_.MainWindowTitle -like "*Text WoW Client*"
}
if ($clientProcesses) {
    Write-Host "[Stop] Found running frontend service, stopping..." -ForegroundColor Yellow
    $clientProcesses | Stop-Process -Force
    Start-Sleep -Seconds 1
}

Write-Host "[Check] Checking port usage..." -ForegroundColor Yellow
$port8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
$port5173 = Get-NetTCPConnection -LocalPort 5173 -ErrorAction SilentlyContinue

if ($port8080) {
    Write-Host "[Warning] Port 8080 is in use" -ForegroundColor Red
}
if ($port5173) {
    Write-Host "[Warning] Port 5173 is in use" -ForegroundColor Red
}

Write-Host ""

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$serverPath = Join-Path $scriptPath "server"
$clientPath = Join-Path $scriptPath "client"

Write-Host "[1/2] Starting backend server (port 8080)..." -ForegroundColor Green
$serverScript = Join-Path $env:TEMP "start-server.ps1"
$serverContent = "Set-Location '$serverPath'`nWrite-Host 'Text WoW Server - Backend' -ForegroundColor Cyan`nWrite-Host 'Port: 8080' -ForegroundColor Gray`nWrite-Host ''`ngo run main.go`n"
$serverContent | Out-File -FilePath $serverScript -Encoding UTF8 -NoNewline
Start-Process powershell -ArgumentList "-NoExit", "-File", $serverScript -WindowStyle Normal

Start-Sleep -Seconds 2

Write-Host "[2/2] Starting frontend server (port 5173)..." -ForegroundColor Green
$clientScript = Join-Path $env:TEMP "start-client.ps1"
$clientContent = "Set-Location '$clientPath'`nWrite-Host 'Text WoW Client - Frontend' -ForegroundColor Cyan`nWrite-Host 'Port: 5173' -ForegroundColor Gray`nWrite-Host ''`nnpm run dev`n"
$clientContent | Out-File -FilePath $clientScript -Encoding UTF8 -NoNewline
Start-Process powershell -ArgumentList "-NoExit", "-File", $clientScript -WindowStyle Normal

Start-Sleep -Seconds 2

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Startup Complete!" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Backend: http://localhost:8080" -ForegroundColor White
Write-Host "  Frontend: http://localhost:5173" -ForegroundColor White
Write-Host ""
Write-Host "  Open http://localhost:5173 in your browser" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Tip: Close services by closing the PowerShell windows" -ForegroundColor Gray
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""



























