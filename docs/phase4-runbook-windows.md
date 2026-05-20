# Phase 4 Runbook (Windows)

## 1) Build RC artifacts

```powershell
./scripts/build.ps1
```

Expected output:
- `bin/server.exe`
- `frontend/dist/*`

## 2) Run application

```powershell
./scripts/run.ps1
```

Health check:

```powershell
curl http://localhost:8080/healthz
```

## 3) Run database migrations

```powershell
$env:DB_URL="sqlite3://app.db"
./scripts/migrate.ps1
```

## 4) Smoke / E2E checks

From repo root:

```powershell
cd frontend
npm install
npx playwright install --with-deps chromium
npx playwright test
```

## 5) Rollback

- Stop `server.exe`
- Restore previous DB backup
- Deploy previous binary + previous `frontend/dist`
