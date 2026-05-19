# Kế hoạch migrate dự án sang Golang + Svelte 5 (target Windows)

## 1) Hiện trạng codebase (branch hiện tại)

- Backend hiện tại là Flask (Python), app factory trong `app/__init__.py`, entrypoint tại `main.py`.
- Chức năng chính:
  - Admin dashboard + CRUD: users, cars, customers, transactions.
  - Homepage route đơn giản (`/`).
  - Auth session-based qua Flask-Login.
  - ORM SQLAlchemy + Flask-Migrate (Alembic).
  - i18n qua Flask-Babel (đang có `vi`, `ja`).
  - Upload ảnh local dưới `static/uploads/...`.
- Frontend hiện tại là server-rendered Jinja templates + assets tĩnh JS/CSS.

## 2) Mục tiêu kiến trúc sau migrate

### 2.1 Kiến trúc tổng thể

- **Backend**: Golang (khuyến nghị `Gin` hoặc `Chi` + `sqlc`/`gorm`) cung cấp REST API.
- **Frontend**: Svelte 5 (Vite), chạy như SPA cho admin + homepage (có thể chia route public/admin).
- **DB**: giữ schema hiện tại, migrate bằng tool Go (khuyến nghị `golang-migrate`).
- **Auth**: chuyển sang cookie session hoặc JWT cookie-based (ưu tiên HttpOnly + Secure + SameSite).
- **i18n**: frontend i18n (vd. `svelte-i18n`), backend trả message code.
- **Deploy trên Windows**:
  - Local dev: backend exe + `npm run dev` frontend.
  - Production: 1 binary Go + static build Svelte (`dist/`) được serve bởi Go hoặc reverse proxy IIS/Nginx for Windows.

### 2.2 Mapping module cũ → mới

- `app/admin/controllers/*` → `internal/http/admin/*` (handler Go).
- `app/admin/models/*`, `app/homepage/models/*` → `internal/domain/*` + `internal/repo/*`.
- `app/admin/services/*` → `internal/service/*`.
- `templates/*` + `static/admin/*` → `frontend/src/*` components/pages.
- `app/utils/settings.py` (config hệ thống) → bảng config + service Go + endpoint `/api/system/settings`.

## 3) Lộ trình migration theo phase

## Phase 0 — Discovery & freeze baseline (2–3 ngày)

1. Chốt phạm vi feature parity v1:
   - Login/logout.
   - Dashboard metrics.
   - CRUD: User/Car/Customer/Transaction.
   - Search.
   - System settings (currency/theme/language).
   - Upload ảnh customer/car.
2. Freeze schema DB hiện tại + export ERD.
3. Thêm test baseline smoke cho luồng quan trọng (để so sánh trước/sau migrate).

**Deliverable**: danh sách endpoint/behavior hiện tại + bộ test smoke tối thiểu.

## Phase 1 — Dựng skeleton Golang backend (3–5 ngày)

1. Tạo cấu trúc dự án Go:
   - `cmd/server/main.go`
   - `internal/config`, `internal/http`, `internal/service`, `internal/repo`, `internal/domain`
   - `migrations/` cho `golang-migrate`
2. Đọc config từ env + file `.env` phù hợp Windows.
3. Setup logger, middleware (recover, request-id, CORS, auth).
4. Kết nối DB và healthcheck `/healthz`.

**Deliverable**: backend Go chạy được trên Windows, kết nối DB OK.

## Phase 2 — Data layer + auth + API lõi (5–8 ngày)

1. Implement model/domain tương đương bảng hiện tại.
2. Implement repository CRUD + filter/search.
3. Implement auth:
   - Login verify password hash tương thích dữ liệu cũ.
   - Session/JWT cookie + middleware bảo vệ route admin.
4. Implement API endpoints (ưu tiên read rồi write):
   - `/api/admin/dashboard`
   - `/api/admin/users`
   - `/api/admin/cars`
   - `/api/admin/customers`
   - `/api/admin/transactions`
   - `/api/admin/search`
   - `/api/admin/system`

**Deliverable**: API parity cho toàn bộ admin use-case.

## Phase 3 — Svelte 5 frontend (6–10 ngày)

1. Scaffold Svelte 5 + Vite + TypeScript.
2. Setup router + layout admin.
3. Tách page theo module: dashboard, users, cars, customers, transactions, system.
4. Form handling + validation client + mapping lỗi từ API.
5. i18n frontend + theme switch.
6. Upload file + preview ảnh.

**Deliverable**: UI Svelte 5 thay thế toàn bộ template Jinja admin.

## Phase 4 — Tích hợp, QA và tối ưu Windows runtime (3–5 ngày)

1. E2E flows quan trọng (Playwright): login, CRUD, search, upload.
2. Benchmark cơ bản (p95 response, memory).
3. Windows packaging:
   - Build `server.exe`.
   - Script PowerShell cho run/migrate.
   - Cấu hình service (NSSM/WinSW) nếu chạy nền.
4. Logging file path và quyền thư mục upload trên Windows.

**Deliverable**: bản release candidate chạy ổn định trên Windows.

## Phase 5 — Cutover (1–2 ngày)

1. Đồng bộ DB migration cuối cùng.
2. Chạy dual-run/blue-green ngắn hạn nếu có điều kiện.
3. Switch traffic sang app Go + Svelte.
4. Theo dõi log/error 24–72 giờ, có plan rollback.

## 4) Thiết kế kỹ thuật đề xuất

## 4.1 Backend Go (khuyến nghị)

- Web framework: `Gin` (nhanh triển khai) hoặc `Chi` (nhẹ, modular).
- ORM/query:
  - Nếu cần tốc độ dev: `GORM`.
  - Nếu cần truy vấn chặt chẽ + hiệu năng: `sqlc` + SQL thuần.
- Migration: `golang-migrate`.
- Validation: `go-playground/validator`.
- Auth:
  - Option A: session store (Redis/in-memory) + cookie.
  - Option B: JWT access + refresh cookie.
- File upload: lưu local path chuẩn hóa (`uploads/cars`, `uploads/customers`), validate mime/size.

## 4.2 Frontend Svelte 5

- Svelte 5 + TypeScript + Vite.
- State: store tối giản (Svelte store), không cần framework state lớn giai đoạn đầu.
- API client: `fetch` wrapper có interceptor xử lý 401/refresh.
- UI strategy:
  - Có thể tái sử dụng CSS hiện có từ `static/admin/css/*` trong giai đoạn chuyển tiếp.
  - Refactor dần sang component hóa.

## 4.3 Hạ tầng Windows

- Build backend: `go build -o bin/server.exe ./cmd/server`.
- FE build: `npm run build` tạo `dist`.
- Run production:
  - Mode 1: Go serve static `dist` + API cùng domain.
  - Mode 2: IIS reverse proxy `/api` -> Go, `/` -> static site.
- Dùng PowerShell scripts:
  - `scripts/dev.ps1`, `scripts/build.ps1`, `scripts/migrate.ps1`, `scripts/run.ps1`.

## 5) Kế hoạch dữ liệu & tương thích

1. Giữ nguyên schema trước, migrate logic trước.
2. Chuyển dần tên cột/chuẩn hóa schema ở phase sau khi cutover ổn định.
3. Password hash:
   - Cần xác nhận format hash hiện tại (Werkzeug `generate_password_hash`) để verify đúng trong Go.
4. Timezone/date parsing:
   - Chuẩn hóa UTC trong DB, format theo locale ở frontend.
5. Seed dữ liệu test + snapshot để regression test.

## 6) Testing strategy

- Unit tests (Go): service, repo query, validation.
- API integration tests: auth + CRUD + search + settings + upload.
- Frontend tests:
  - Component tests cơ bản.
  - E2E Playwright cho critical flows.
- So sánh parity:
  - Dùng checklist response/behavior giữa Flask và Go.

## 7) Rủi ro chính & giảm thiểu

1. **Sai lệch behavior business logic** giữa Flask và Go.
   - Giảm thiểu: tạo test parity trước khi migrate.
2. **Auth/session không tương thích** gây logout bất thường.
   - Giảm thiểu: chọn 1 chiến lược auth rõ ràng, test đa trình duyệt.
3. **Upload file trên Windows path/permission issues**.
   - Giảm thiểu: normalize path, tách config upload dir, test với user quyền hạn chế.
4. **i18n regressions**.
   - Giảm thiểu: map key i18n ngay từ đầu, không hardcode text.
5. **Độ trễ migrate UI lớn**.
   - Giảm thiểu: ưu tiên module admin quan trọng trước, rollout incremental.

## 8) Kế hoạch nhân sự/ước lượng

- Team tối thiểu: 1 BE Go, 1 FE Svelte, 1 QA (part-time).
- Ước lượng tổng: **4–6 tuần** cho feature parity + hardening cơ bản.
- Nếu 1 người fullstack: **6–8 tuần** thực tế hơn.

## 9) Backlog implementation gợi ý (ưu tiên)

1. [P0] Khởi tạo backend Go + healthcheck + DB connect.
2. [P0] Auth login/logout + middleware.
3. [P0] Users/Cars/Customers/Transactions read APIs.
4. [P0] Write APIs + validation.
5. [P1] Dashboard metrics + search.
6. [P1] System settings + i18n.
7. [P1] Upload images.
8. [P2] Frontend polish + accessibility + performance tuning.

## 10) Đề xuất triển khai thực tế

- Bắt đầu theo hướng **Strangler Fig**:
  - Dựng Go API song song app Flask hiện tại.
  - Svelte 5 consume Go API cho từng module admin mới.
  - Module nào hoàn thành thì chuyển route/module đó trước.
- Tránh big-bang rewrite một lần để giảm rủi ro cutover.
