$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

Write-Host "== Rollback Procedure =="

Write-Host "1) Stop Go service process"
Get-Process server -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue

Write-Host "2) Restore latest DB backup (if present)"
$backup = Get-ChildItem -Path . -Filter "app.db.backup.*" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if ($backup) {
  Copy-Item $backup.FullName "app.db" -Force
  Write-Host "Restored DB from $($backup.Name)"
} else {
  Write-Host "No backup found; skipped restore"
}

Write-Host "3) Revert traffic to previous stack manually (IIS/Nginx/old service)"
Write-Host "Rollback script finished."
