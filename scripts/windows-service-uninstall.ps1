$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

param(
  [string]$ServiceName = "gogemini",
  [string]$NssmPath = "nssm"
)

$null = Get-Command $NssmPath -ErrorAction Stop

& $NssmPath stop $ServiceName 2>$null
& $NssmPath remove $ServiceName confirm

Write-Host "Service '$ServiceName' removed."
