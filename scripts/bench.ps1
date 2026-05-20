$ErrorActionPreference = "Stop"

Write-Host "Simple benchmark against /healthz (10 requests)"
1..10 | ForEach-Object {
  Measure-Command { curl http://localhost:8080/healthz | Out-Null } | Select-Object TotalMilliseconds
}
