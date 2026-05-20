$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$logDir = "logs/cutover"
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$logFile = Join-Path $logDir "rollback-$timestamp.log"

function Write-Log([string]$message) {
  $line = "[$(Get-Date -Format s)] $message"
  Write-Host $line
  Add-Content -Path $logFile -Value $line
}

Write-Log "== Rollback Procedure =="
$start = Get-Date

Write-Log "1) Stop Go service process"
Get-Process server -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue

Write-Log "2) Restore latest DB backup (if present)"
$backup = Get-ChildItem -Path . -Filter "app.db.backup.*" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if ($backup) {
  Copy-Item $backup.FullName "app.db" -Force
  Write-Log "Restored DB from $($backup.Name)"
} else {
  Write-Log "No backup found; skipped restore"
}

$duration = [math]::Round(((Get-Date) - $start).TotalSeconds, 2)
Write-Log "3) Revert traffic to previous stack manually (IIS/Nginx/old service)"
Write-Log "Rollback drill duration (seconds): $duration"
Write-Log "Rollback script finished. Evidence log: $logFile"
