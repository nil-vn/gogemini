$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

if (!(Get-Command migrate -ErrorAction SilentlyContinue)) {
  Write-Host "golang-migrate CLI is required: https://github.com/golang-migrate/migrate"
  exit 1
}

if (-not $env:DB_URL) {
  $env:DB_URL = "sqlite3://app.db"
}

migrate -path migrations -database $env:DB_URL up
