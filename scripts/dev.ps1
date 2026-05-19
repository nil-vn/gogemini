$env:SERVER_ADDR = ":8080"
$env:DB_DRIVER = "sqlite"
$env:DB_DSN = "file:app.db?cache=shared"
$env:CORS_ORIGIN = "*"

go run ./cmd/server
