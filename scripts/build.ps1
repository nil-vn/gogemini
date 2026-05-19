$ErrorActionPreference = "Stop"
Set-Location "$PSScriptRoot/.."

Write-Host "[1/2] Build Go server"
go build -o bin/server.exe ./cmd/server

Write-Host "[2/2] Build frontend"
Set-Location "frontend"
npm install
npm run build

Write-Host "Build completed: bin/server.exe and frontend/dist"
