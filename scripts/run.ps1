$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

if (!(Test-Path "bin/server.exe")) {
  Write-Host "server.exe not found, running build first..."
  & "$PSScriptRoot/build.ps1"
}

$env:SERVER_ADDR = ":8080"
$env:DB_DRIVER = "sqlite"
$env:DB_DSN = "file:app.db?cache=shared"
$env:CORS_ORIGIN = "*"
$env:UPLOAD_ROOT = "static"

./bin/server.exe
