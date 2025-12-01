# TEXT WoW - Stop services script (PowerShell)
# Usage: .\stop.ps1

$ErrorActionPreference = "SilentlyContinue"

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  TEXT WoW - Stop Services Script" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "[Stop] Stopping backend service..." -ForegroundColor Yellow
$serverProcesses = Get-Process | Where-Object { 
    $_.MainWindowTitle -like "*Text WoW Server*"
}
if ($serverProcesses) {
    $serverProcesses | Stop-Process -Force
    Write-Host "[Success] Backend service stopped" -ForegroundColor Green
} else {
    Write-Host "[Info] No running backend service found" -ForegroundColor Gray
}

Write-Host "[Stop] Stopping frontend service..." -ForegroundColor Yellow
$clientProcesses = Get-Process | Where-Object { 
    $_.MainWindowTitle -like "*Text WoW Client*"
}
if ($clientProcesses) {
    $clientProcesses | Stop-Process -Force
    Write-Host "[Success] Frontend service stopped" -ForegroundColor Green
} else {
    Write-Host "[Info] No running frontend service found" -ForegroundColor Gray
}

Write-Host "[Check] Checking port usage..." -ForegroundColor Yellow
$port8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
if ($port8080) {
    $pid = $port8080.OwningProcess
    Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
    Write-Host "[Clean] Cleared port 8080" -ForegroundColor Yellow
}

$port5173 = Get-NetTCPConnection -LocalPort 5173 -ErrorAction SilentlyContinue
if ($port5173) {
    $pid = $port5173.OwningProcess
    Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
    Write-Host "[Clean] Cleared port 5173" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  All services stopped" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""
