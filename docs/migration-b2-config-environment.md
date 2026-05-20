# B2 - Config & Environment Strategy (Windows-first)

Tài liệu này chốt chiến lược config/env cho task **B2**.

## Env variables chuẩn hóa (B2 scope)

| Key | Required | Default | Mô tả |
|---|---|---|---|
| `APP_ENV` | No | `development` | Runtime environment (`development`, `staging`, `production`). |
| `SERVER_ADDR` | No | `:8080` | Địa chỉ bind server Go. |
| `DB_DRIVER` | Yes | `sqlite` | DB driver hiện hỗ trợ (`sqlite` baseline). |
| `DB_DSN` | Yes | `file:app.db?cache=shared` | Chuỗi kết nối DB. |
| `AUTH_SECRET` | Yes | `dev-change-me` | Secret dùng cho auth/session signing ở các phase sau. |
| `CORS_ORIGIN` | No | `*` | CORS allow-origin baseline. |
| `UPLOAD_DIR` | Yes | `static/uploads` | Thư mục lưu upload; hỗ trợ path kiểu Windows (`C:\\...`). |

## Quy tắc load config

1. App tự đọc file `.env` nếu tồn tại ở repo root.
2. Environment variables của OS **ưu tiên cao hơn** `.env` (không bị ghi đè).
3. Nếu biến required bị rỗng (`DB_DRIVER`, `DB_DSN`, `AUTH_SECRET`, `UPLOAD_DIR`) thì app fail-fast khi boot.

## Minimal boot config

App có thể boot với cấu hình tối thiểu:

- Không cần tạo `.env` (dùng defaults đã định nghĩa).
- Hoặc copy `.env.example` -> `.env` và chỉnh lại cho môi trường thật.

## Windows-first notes

- Ưu tiên khai báo bằng PowerShell `$env:KEY="value"` để override nhanh theo môi trường.
- `UPLOAD_DIR` chấp nhận cả slash (`static/uploads`) và backslash path tuyệt đối Windows.
- Script `scripts/dev.ps1` và `scripts/run.ps1` đã thống nhất key `UPLOAD_DIR`.
