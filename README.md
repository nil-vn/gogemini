# GoGemini - Hướng dẫn chạy local và build/release

Tài liệu này hướng dẫn **từng bước** để:

1. Chạy app ở môi trường **dev** với **backend server** và **frontend server tách biệt**.
2. Build và launch app ở môi trường gần production trên **Linux** và **Windows**.

---

## 1) Tổng quan kiến trúc runtime

- **Backend (Go/Gin)** chạy mặc định tại `http://localhost:8080`.
- **Frontend (Svelte + Vite)** khi chạy dev dùng Vite server (mặc định `http://localhost:5173`).
- Health endpoints của backend:
  - `GET /healthz`
  - `GET /livez`
  - `GET /readyz`

---

## 2) Prerequisites

### Linux
- Go (khuyến nghị bản stable mới)
- Node.js 20+
- npm
- Git

### Windows
- Go
- Node.js 20+
- npm
- Git
- PowerShell 5.1+ hoặc PowerShell 7+

> Gợi ý kiểm tra nhanh:
>
> - `go version`
> - `node -v`
> - `npm -v`

---

## 3) Cấu hình môi trường

App tự đọc `.env` tại repo root (nếu có). OS env var có độ ưu tiên cao hơn `.env`.

Tạo file `.env` (hoặc copy từ `.env.example` nếu có) với cấu hình tối thiểu:

```env
APP_ENV=development
SERVER_ADDR=:8080
DB_DRIVER=sqlite
DB_DSN=file:app.db?cache=shared
AUTH_SECRET=dev-change-me
CORS_ORIGIN=http://localhost:5173
UPLOAD_DIR=static/uploads
```

---

## 3.1) Tự động khởi tạo database + tables + default admin

Backend Go **tự tạo schema** khi boot:
- Nếu DB file chưa tồn tại (ví dụ `app.db`) thì SQLite sẽ tạo file.
- App sẽ tự chạy `CREATE TABLE IF NOT EXISTS` cho toàn bộ bảng hệ thống (`users`, `car`, `customer`, `transaction`, `transaction_item`, `transaction_car`, `car_image`, `customer_image`, `config`).
- Không cần chạy script migration thủ công.

### Cấu hình default admin trong `.env` (tuỳ chọn)

Để app tự tạo/cập nhật admin mặc định ngay khi khởi động, set các biến sau:

```env
DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_PASSWORD=your-strong-password
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_ROLE=admin
DEFAULT_ADMIN_STATUS=active
```

Quy tắc:
- `DEFAULT_ADMIN_USERNAME` + `DEFAULT_ADMIN_PASSWORD`: **bắt buộc** nếu muốn bootstrap admin.
- `DEFAULT_ADMIN_EMAIL`: optional, mặc định `<username>@local`.
- `DEFAULT_ADMIN_ROLE`: optional, mặc định `admin`.
- `DEFAULT_ADMIN_STATUS`: optional, mặc định `active`.
- Nếu user đã tồn tại theo `username` hoặc `email`, app sẽ update password/role/status theo giá trị mới.


## 3.2) Cấu hình môi trường cho frontend (Vite)

Frontend sử dụng biến môi trường riêng của Vite, quan trọng nhất là `VITE_API_BASE`.

- File mẫu đã có sẵn tại `frontend/.env.example`:

```env
VITE_API_BASE=http://localhost:8080
```

- Bạn cần tạo `frontend/.env` (hoặc copy từ file mẫu) để tránh nhầm host/port khi gọi API:

```bash
cp frontend/.env.example frontend/.env
```

- Nếu backend không chạy ở `localhost:8080`, bắt buộc sửa `VITE_API_BASE` cho đúng endpoint backend thực tế.
- Sau khi chỉnh `frontend/.env`, cần restart Vite dev server để nhận cấu hình mới.

## 4) DEV mode (backend & frontend tách biệt)

## 4.1 Linux

### Bước 1: Cài dependencies frontend

```bash
cd frontend
npm install
cd ..
```

### Bước 2: Chạy backend server (Terminal 1)

```bash
export SERVER_ADDR=":8080"
export DB_DRIVER="sqlite"
export DB_DSN="file:app.db?cache=shared"
export CORS_ORIGIN="http://localhost:5173"
export UPLOAD_DIR="static/uploads"
go run ./cmd/server
```

### Bước 3: Chạy frontend dev server (Terminal 2)

```bash
cd frontend
npm run dev -- --host 0.0.0.0 --port 5173
```

### Bước 4: Verify

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Mở frontend: `http://localhost:5173`.

---

## 4.2 Windows (PowerShell)

### Bước 1: Cài dependencies frontend

```powershell
cd frontend
npm install
cd ..
```

### Bước 2: Chạy backend server (PowerShell 1)

```powershell
$env:SERVER_ADDR=":8080"
$env:DB_DRIVER="sqlite"
$env:DB_DSN="file:app.db?cache=shared"
$env:CORS_ORIGIN="*"
$env:UPLOAD_DIR="static/uploads"
go run ./cmd/server
```

### Bước 3: Chạy frontend dev server (PowerShell 2)

```powershell
cd frontend
npm run dev -- --host 0.0.0.0 --port 5173
```

### Bước 4: Verify

```powershell
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Mở frontend: `http://localhost:5173`.

---

## 5) Build + Launch (Linux)

## 5.1 Build

### Bước 1: Build backend binary

```bash
mkdir -p bin
go build -o bin/server ./cmd/server
```

### Bước 2: Build frontend static assets

```bash
cd frontend
npm install
npm run build
cd ..
```

Kết quả mong đợi:
- `bin/server`
- `frontend/dist/*`

## 5.2 Launch

```bash
export SERVER_ADDR=":8080"
export DB_DRIVER="sqlite"
export DB_DSN="file:app.db?cache=shared"
export CORS_ORIGIN="*"
export UPLOAD_DIR="static/uploads"
./bin/server
```

## 5.3 Smoke check

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/livez
curl http://localhost:8080/readyz
```

---

## 6) Build + Launch (Windows)

## 6.1 Build

Dùng script chuẩn:

```powershell
./scripts/build.ps1
```

Kết quả mong đợi:
- `bin/server.exe`
- `frontend/dist/*`

## 6.2 Apply migrations

```powershell
$env:DB_URL="sqlite3://app.db"
./scripts/migrate.ps1
```

## 6.3 Launch

Dùng script run:

```powershell
./scripts/run.ps1
```

Hoặc chạy trực tiếp:

```powershell
$env:SERVER_ADDR=":8080"
$env:DB_DRIVER="sqlite"
$env:DB_DSN="file:app.db?cache=shared"
$env:CORS_ORIGIN="*"
$env:UPLOAD_DIR="static/uploads"
./bin/server.exe
```

## 6.4 Smoke check

```powershell
curl http://localhost:8080/healthz
curl http://localhost:8080/livez
curl http://localhost:8080/readyz
```

---

## 7) Troubleshooting nhanh

- `readyz` fail: kiểm tra `DB_DSN`, quyền ghi file SQLite, path thư mục upload.
- Frontend gọi API bị CORS: kiểm tra `CORS_ORIGIN` ở backend.
- Port conflict: đổi `SERVER_ADDR` (vd `:8081`) hoặc đổi port Vite (`--port 5174`).
- Windows path lỗi: thử dùng slash `/` trong `UPLOAD_DIR` (`static/uploads`) hoặc path tuyệt đối Windows.

---

## 8) Lệnh nhanh (cheat sheet)

### Dev tách backend/frontend

- Backend: `go run ./cmd/server`
- Frontend: `cd frontend && npm run dev`

### Build release

- Linux backend: `go build -o bin/server ./cmd/server`
- Windows full build: `./scripts/build.ps1`
- Frontend: `cd frontend && npm run build`

### Health checks

- `curl http://localhost:8080/healthz`
- `curl http://localhost:8080/livez`
- `curl http://localhost:8080/readyz`
