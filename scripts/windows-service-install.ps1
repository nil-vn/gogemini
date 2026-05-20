$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

param(
  [string]$ServiceName = "gogemini",
  [string]$NssmPath = "nssm",
  [string]$ServerExePath = "bin/server.exe",
  [string]$RunScriptPath = "scripts/run.ps1",
  [string]$AppDirectory = (Get-Location).Path
)

if (!(Test-Path $ServerExePath)) {
  throw "Missing $ServerExePath. Run ./scripts/build.ps1 first."
}

if (!(Test-Path $RunScriptPath)) {
  throw "Missing $RunScriptPath."
}

$null = Get-Command $NssmPath -ErrorAction Stop

$existing = & $NssmPath status $ServiceName 2>$null
if ($LASTEXITCODE -eq 0 -and $existing) {
  Write-Host "Service '$ServiceName' already exists. Updating configuration..."
} else {
  & $NssmPath install $ServiceName "powershell.exe" "-ExecutionPolicy Bypass -File \"$AppDirectory\\$RunScriptPath\""
}

& $NssmPath set $ServiceName AppDirectory $AppDirectory
& $NssmPath set $ServiceName Start SERVICE_AUTO_START
& $NssmPath set $ServiceName AppStdout "$AppDirectory\\logs\\server.stdout.log"
& $NssmPath set $ServiceName AppStderr "$AppDirectory\\logs\\server.stderr.log"
& $NssmPath set $ServiceName AppRotateFiles 1
& $NssmPath set $ServiceName AppRotateOnline 1
& $NssmPath set $ServiceName AppRotateSeconds 86400

Write-Host "Service '$ServiceName' configured successfully."
Write-Host "Start service with: $NssmPath start $ServiceName"
