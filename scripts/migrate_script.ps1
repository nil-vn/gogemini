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
  Write-Host "migrate CLI not found, trying to install via 'go install'..."
  if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go is required to auto-install migrate CLI. Please install Go or install migrate manually."
    exit 1
  }

  go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  if ($LASTEXITCODE -ne 0) {
    Write-Host "Auto-install migrate failed. Please install manually: https://github.com/golang-migrate/migrate"
    exit 1
  }

  $goBin = (go env GOPATH).Trim() + "\bin"
  $env:Path = $env:Path + ";" + $goBin
}

if (!(Get-Command migrate -ErrorAction SilentlyContinue)) {
  Write-Host "Auto-install migrate failed. Please install manually: https://github.com/golang-migrate/migrate"
  exit 1
}

if ([string]::IsNullOrWhiteSpace($Email)) {
  $Email = "$Username@local"
}

Write-Host "[1/2] Applying migrations to $DbUrl"
$migrationsPath = Join-Path (Get-Location) "migrations"
migrate -path $migrationsPath -database $DbUrl up

Write-Host "[2/2] Ensuring default admin user '$Username' exists"
go run ./cmd/bootstrap-admin --db-driver=$DbDriver --db-dsn=$DbDsn --username=$Username --password=$Password --role=$Role --status=$Status --email=$Email

Write-Host "Done"
