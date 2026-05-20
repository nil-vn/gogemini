param(
  [switch]$SkipBuild,
  [switch]$SkipMigrate,
  [int]$Port = 8080
)

$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$logDir = "logs/cutover"
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$logFile = Join-Path $logDir "cutover-$timestamp.log"

function Write-Log([string]$message) {
  $line = "[$(Get-Date -Format s)] $message"
  Write-Host $line
  Add-Content -Path $logFile -Value $line
}

Write-Log "== Phase 5 Cutover =="

if (-not $SkipBuild) {
  if (!(Test-Path "bin/server.exe") -or !(Test-Path "frontend/dist/index.html")) {
    Write-Log "Build artifacts missing. Running build script..."
    & "$PSScriptRoot/build.ps1"
  }
}

Write-Log "1) Backup DB"
if (Test-Path "app.db") {
  $backupFile = "app.db.backup.$timestamp"
  Copy-Item "app.db" $backupFile
  Write-Log "Created backup: $backupFile"
} else {
  Write-Log "app.db not found; skip backup"
}

if (-not $SkipMigrate) {
  Write-Log "2) Apply migrations"
  & "$PSScriptRoot/migrate.ps1"
}

Write-Log "3) Start new service"
$env:SERVER_ADDR = ":$Port"
$env:DB_DRIVER = "sqlite"
$env:DB_DSN = "file:app.db?cache=shared"
$env:CORS_ORIGIN = "*"
$env:UPLOAD_ROOT = "static"

Start-Process -FilePath "./bin/server.exe" -NoNewWindow
Start-Sleep -Seconds 2

Write-Log "4) Smoke checks"
$base = "http://localhost:$Port"
$health = curl "$base/healthz"
$ready = curl "$base/readyz"
$live = curl "$base/livez"
Write-Log "healthz=$health"
Write-Log "readyz=$ready"
Write-Log "livez=$live"

Write-Log "Cutover completed. Continue 24-72h hypercare monitoring per runbook."
Write-Log "Evidence log: $logFile"
