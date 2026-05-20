# B4 - DB connect + healthcheck

## Scope
- Cấu hình DB connection pool qua env để kiểm soát số kết nối và vòng đời connection.
- Hoàn thiện endpoint `/healthz` kiểm tra cả app process và DB connectivity.
- Bổ sung test để xác minh trạng thái healthy/unhealthy.

## DB pool config
Biến env mới:
- `DB_MAX_OPEN_CONNS` (default `10`, bắt buộc > 0)
- `DB_MAX_IDLE_CONNS` (default `5`, bắt buộc > 0 và <= `DB_MAX_OPEN_CONNS`)
- `DB_CONN_MAX_LIFETIME_MS` (default `300000`, bắt buộc >= 0)

## Healthcheck contract
- `GET /healthz`
  - 200: `{"status":"ok","app":"up","db":"up"}` khi DB ping thành công.
  - 503: `{"status":"unhealthy","app":"up","db":"down","error":"..."}` khi DB ping lỗi.

## Verification
- `go test ./...`
- `go run ./cmd/server` + `curl http://localhost:8080/healthz`
