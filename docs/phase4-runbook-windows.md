# Phase 4 Runbook (Windows) — E4 Packaging & Runtime

Runbook này là hướng dẫn chuẩn để dựng runtime mới trên **Windows host sạch** cho release candidate migration.

## 0) Preconditions

- Windows PowerShell 5.1+ hoặc PowerShell 7+.
- Đã cài: Go, Node.js 20+, npm.
- (Nếu chạy service nền) đã cài NSSM (`nssm.exe`) và thêm vào `PATH`.
- Có file `.env` theo mẫu `.env.example`.

## 1) Build RC artifacts

```powershell
./scripts/build.ps1
```

Expected output:
- `bin/server.exe`
- `frontend/dist/*`

## 2) Run DB migrations

```powershell
$env:DB_URL="sqlite3://app.db"
./scripts/migrate.ps1
```

Expected output:
- Migration chạy `up` thành công, không lỗi schema.

## 3) Run application trực tiếp (foreground)

```powershell
./scripts/run.ps1
```

Health checks:

```powershell
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/livez
```

## 4) Cài chạy nền bằng NSSM (production-like)

Cài service:

```powershell
./scripts/windows-service-install.ps1 -ServiceName gogemini
```

Start/verify service:

```powershell
nssm start gogemini
nssm status gogemini
```

Stop & uninstall service:

```powershell
./scripts/windows-service-uninstall.ps1 -ServiceName gogemini
```

## 5) Smoke / E2E checks

Từ repo root:

```powershell
cd frontend
npm install
npx playwright install --with-deps chromium
npx playwright test
```

## 6) Runtime checklist cho host mới

- Tạo sẵn thư mục `logs/` và đảm bảo quyền ghi cho user chạy service.
- Tạo sẵn thư mục upload theo `UPLOAD_DIR` và đảm bảo quyền ghi.
- Kiểm tra firewall mở port app (`8080` hoặc port cấu hình).
- Xác minh có thể đọc env secrets từ môi trường hệ thống.

## 7) Rollback

- Stop service hoặc stop `server.exe`.
- Restore DB backup gần nhất.
- Deploy lại binary cũ + `frontend/dist` cũ.
- Chạy lại smoke tối thiểu (`healthz`, login, list users/cars).

## 8) Security hardening operations (WBS-1)

- Verify readiness/liveness split:
  - `curl http://localhost:8080/livez`
  - `curl http://localhost:8080/readyz`
- Verify structured logs are JSON and include `request_id`.
- Confirm login lockout behavior after 5 failures (expect `AUTH_LOCKED`).
- Confirm auth/logout envelopes and cookie invalidation on logout.
