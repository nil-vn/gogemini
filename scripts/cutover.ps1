$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

Write-Host "== Phase 5 Cutover =="

if (!(Test-Path "bin/server.exe")) {
  Write-Host "Build artifacts missing. Running build script..."
  & "$PSScriptRoot/build.ps1"
}

if (!(Test-Path "frontend/dist/index.html")) {
  Write-Host "frontend/dist missing. Running build script..."
  & "$PSScriptRoot/build.ps1"
}

Write-Host "1) Backup DB"
if (Test-Path "app.db") {
  $ts = Get-Date -Format "yyyyMMdd-HHmmss"
  Copy-Item "app.db" "app.db.backup.$ts"
  Write-Host "Created backup: app.db.backup.$ts"
} else {
  Write-Host "app.db not found; skip backup"
}

Write-Host "2) Apply migrations"
try {
  & "$PSScriptRoot/migrate.ps1"
} catch {
  Write-Host "Migration failed; aborting cutover"
  exit 1
}

Write-Host "3) Start new service"
$env:SERVER_ADDR = ":8080"
$env:DB_DRIVER = "sqlite"
$env:DB_DSN = "file:app.db?cache=shared"
$env:CORS_ORIGIN = "*"
$env:UPLOAD_ROOT = "static"

Start-Process -FilePath "./bin/server.exe" -NoNewWindow
Start-Sleep -Seconds 2

Write-Host "4) Smoke check /healthz"
try {
  $health = curl http://localhost:8080/healthz
  Write-Host $health
} catch {
  Write-Host "Health check failed after cutover"
  exit 1
}

Write-Host "Cutover completed. Monitor logs for 24-72h per runbook."
