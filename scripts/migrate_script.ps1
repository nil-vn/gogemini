$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

param(
  [Parameter(Mandatory=$true)][string]$Username,
  [Parameter(Mandatory=$true)][string]$Password,
  [string]$DbUrl = "sqlite3://app.db",
  [string]$DbDriver = "sqlite",
  [string]$DbDsn = "file:app.db?cache=shared",
  [string]$Role = "admin",
  [string]$Status = "active",
  [string]$Email = ""
)

if (!(Get-Command migrate -ErrorAction SilentlyContinue)) {
  Write-Host "golang-migrate CLI is required: https://github.com/golang-migrate/migrate"
  exit 1
}

if ([string]::IsNullOrWhiteSpace($Email)) {
  $Email = "$Username@local"
}

Write-Host "[1/2] Applying migrations to $DbUrl"
migrate -path migrations -database $DbUrl up

Write-Host "[2/2] Ensuring default admin user '$Username' exists"
go run ./cmd/bootstrap-admin --db-driver=$DbDriver --db-dsn=$DbDsn --username=$Username --password=$Password --role=$Role --status=$Status --email=$Email

Write-Host "Done"
