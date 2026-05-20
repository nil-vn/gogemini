# Task Breakdown triển khai Migration (Go + Svelte 5)

Tài liệu này tách chi tiết các task từ `docs/migration-plan-go-svelte5.md` thành backlog có thể giao việc, estimate theo ngày công (ideal), phụ thuộc, và tiêu chí hoàn thành (DoD).

## 1) Nguyên tắc lập kế hoạch

- Ưu tiên **feature parity** trước tối ưu.
- Triển khai theo hướng **Strangler Fig** (song song Flask cũ và Go/Svelte mới).
- Mỗi task phải có:
  - Owner chính (BE/FE/QA/DevOps).
  - Dependency rõ ràng.
  - Definition of Done (DoD) kiểm chứng được.

---

## 2) Work Breakdown Structure (WBS)

## EPIC A — Baseline & Discovery (Phase 0)

### A1. Khóa phạm vi v1 parity
- **Owner**: PM/Tech Lead + BE + FE
- **Estimate**: 0.5 ngày
- **Công việc**:
  - Chốt danh sách flow bắt buộc: auth, dashboard, CRUD 4 module, search, system settings, upload.
  - Chốt out-of-scope cho v1.
- **DoD**:
  - Có checklist parity v1 dạng table (feature, route hiện tại, expected behavior).

### A2. Audit endpoint & behavior Flask hiện tại
- **Owner**: BE
- **Estimate**: 1 ngày
- **Công việc**:
  - Inventory tất cả route quan trọng, payload, response, auth requirement.
  - Ghi nhận đặc thù business rules (validation, sorting default, pagination).
- **DoD**:
  - Có API contract baseline (OpenAPI draft hoặc Markdown spec).

### A3. Freeze schema DB + ERD
- **Owner**: BE/DBA
- **Estimate**: 0.5 ngày
- **Công việc**:
  - Chốt version migration hiện tại.
  - Export ERD và data dictionary.
- **DoD**:
  - Có file snapshot schema + ERD đính kèm docs.

### A4. Baseline smoke tests (hệ cũ)
- **Owner**: QA + BE
- **Estimate**: 1 ngày
- **Công việc**:
  - Viết test smoke cho login + 1 CRUD flow/module + search + upload.
- **DoD**:
  - Smoke tests chạy pass ổn định trên branch baseline.

---

## EPIC B — Backend Foundation (Phase 1)

### B1. Scaffold project Go
- **Owner**: BE
- **Estimate**: 1 ngày
- **Dependency**: A2, A3
- **Công việc**:
  - Tạo cấu trúc `cmd/`, `internal/`, `migrations/`, `scripts/`.
  - Thiết lập module Go, convention package.
- **DoD**:
  - `go test ./...` chạy được (kể cả khi test rỗng).

### B2. Config & environment strategy (Windows-first)
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: B1
- **Công việc**:
  - Chuẩn hóa env vars (DB, auth secret, upload dir, CORS).
  - Hỗ trợ `.env` + override qua environment.
- **DoD**:
  - App boot được với env tối thiểu; có tài liệu sample env.

### B3. Middleware core
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: B1
- **Công việc**:
  - request-id, recover panic, logging, CORS.
- **DoD**:
  - Mọi request có request-id; panic được xử lý không crash process.

### B4. DB connect + healthcheck
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: B2
- **Công việc**:
  - DB pool config.
  - `/healthz` check app + DB.
- **DoD**:
  - Healthcheck pass trong local Windows/macOS/Linux.

---

## EPIC C — Data Layer + Auth + Core API (Phase 2)

### C1. Domain model mapping
- **Owner**: BE
- **Estimate**: 1 ngày
- **Dependency**: A3, B4
- **Công việc**:
  - Map bảng hiện tại sang entity/domain structs.
- **DoD**:
  - Entity đầy đủ trường cần cho CRUD + dashboard.

### C2. Repository CRUD (Users/Cars/Customers/Transactions)
- **Owner**: BE
- **Estimate**: 2 ngày
- **Dependency**: C1
- **Công việc**:
  - List/detail/create/update/delete.
  - Pagination/sort/filter chuẩn hóa.
- **DoD**:
  - Integration test pass cho CRUD mỗi module.

### C3. Search API
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: C2
- **Công việc**:
  - Search đa module theo behavior cũ.
- **DoD**:
  - Kết quả search khớp baseline test cases.

### C4. Dashboard metrics API
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: C2
- **Công việc**:
  - Tổng hợp metric chính (count/revenue/... theo hệ cũ).
- **DoD**:
  - API trả đúng số liệu trên snapshot dataset.

### C5. Auth migration
- **Owner**: BE
- **Estimate**: 1 ngày
- **Dependency**: A2, C1
- **Công việc**:
  - Verify password hash tương thích dữ liệu Flask.
  - Login/logout + session/JWT cookie + middleware bảo vệ admin.
- **DoD**:
  - User cũ đăng nhập được, route admin bị chặn khi chưa auth.

### C6. System settings API
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: C2
- **Công việc**:
  - Read/update config (currency/theme/language).
- **DoD**:
  - Có endpoint GET/PUT settings + validation.

### C7. Upload API
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: B2, C5
- **Công việc**:
  - Upload ảnh car/customer, validate mime/size, normalize path.
- **DoD**:
  - Upload thành công và trả URL/path đúng trên Windows.

---

## EPIC D — Frontend Svelte 5 (Phase 3)

### D1. FE scaffold & toolchain
- **Owner**: FE
- **Estimate**: 0.5 ngày
- **Dependency**: B1
- **Công việc**:
  - Svelte 5 + Vite + TypeScript + lint/format.
- **DoD**:
  - Build dev/prod chạy ổn.

### D2. App shell + routing + auth guard
- **Owner**: FE
- **Estimate**: 1 ngày
- **Dependency**: D1, C5
- **Công việc**:
  - Layout admin, route grouping, intercept 401.
- **DoD**:
  - Route admin yêu cầu auth; logout redirect đúng.

### D3. Module pages parity (Dashboard, Users, Cars, Customers, Transactions)
- **Owner**: FE
- **Estimate**: 3 ngày
- **Dependency**: C2, C4, D2
- **Công việc**:
  - Bảng danh sách + form create/update + detail cơ bản.
- **DoD**:
  - Hoàn tất end-to-end UI flow cho 5 module.

### D4. Search + System settings + i18n
- **Owner**: FE
- **Estimate**: 1 ngày
- **Dependency**: C3, C6, D2
- **Công việc**:
  - Ô search global/module.
  - Màn settings + language/theme switch.
- **DoD**:
  - Text hiển thị qua key i18n, không hardcode nội dung chính.

### D5. Upload UI + image preview
- **Owner**: FE
- **Estimate**: 0.5 ngày
- **Dependency**: C7, D3
- **Công việc**:
  - Chọn file, preview, handle lỗi upload.
- **DoD**:
  - Người dùng upload ảnh trong form car/customer thành công.

---

## EPIC E — QA, Packaging, Cutover (Phase 4-5)

### E1. E2E Playwright critical flows
- **Owner**: QA + FE
- **Estimate**: 1 ngày
- **Dependency**: D5
- **Công việc**:
  - Kịch bản login, CRUD, search, upload.
- **DoD**:
  - Test critical pass trên CI/local.

### E2. API parity test suite
- **Owner**: QA + BE
- **Estimate**: 0.5 ngày
- **Dependency**: C7
- **Công việc**:
  - So sánh response/behavior Flask vs Go cho case quan trọng.
- **DoD**:
  - Có checklist parity ký xác nhận.

### E3. Performance baseline
- **Owner**: BE
- **Estimate**: 0.5 ngày
- **Dependency**: C7
- **Công việc**:
  - Benchmark p95 endpoint chính, memory footprint.
- **DoD**:
  - Có report trước release candidate.

### E4. Windows packaging & runbook
- **Owner**: DevOps/BE
- **Estimate**: 1 ngày
- **Dependency**: D5
- **Công việc**:
  - Build `server.exe`, FE `dist`, scripts PowerShell build/run/migrate.
  - Cấu hình chạy nền (NSSM/WinSW).
- **DoD**:
  - Môi trường Windows mới có thể setup và chạy theo runbook.

### E5. Cutover & hypercare
- **Owner**: Tech Lead + BE + QA
- **Estimate**: 1 ngày + theo dõi 2–3 ngày
- **Dependency**: E1-E4
- **Công việc**:
  - Final DB migrate, switch traffic, theo dõi log/error.
  - Chuẩn bị rollback procedure.
- **DoD**:
  - Không có incident P1/P2 trong 24h đầu sau cutover.

---

## 3) Dependency graph (rút gọn)

- `A1/A2/A3` -> `B1` -> `B2/B3` -> `B4`
- `B4 + A3` -> `C1` -> `C2` -> `C3/C4/C6`
- `A2 + C1` -> `C5` -> `D2`
- `B2 + C5` -> `C7` -> `D5`
- `D1` -> `D2` -> `D3` -> `D5`
- `C3/C6 + D2` -> `D4`
- `D5` -> `E1/E4` ; `C7` -> `E2/E3`; tất cả -> `E5`

---

## 4) Sprint gợi ý (2 tuần/sprint)

### Sprint 1
- A1-A4, B1-B4, C1, C5 (mốc: login + healthcheck + DB running).

### Sprint 2
- C2, C3, C4, C6, C7, D1, D2 (mốc: API core hoàn chỉnh + FE shell).

### Sprint 3
- D3, D4, D5, E1, E2 (mốc: feature parity + E2E pass).

### Sprint 4
- E3, E4, E5 (mốc: RC trên Windows + cutover).

---

## 5) Phân công theo vai trò

- **BE Go**: B1-B4, C1-C7, E3.
- **FE Svelte**: D1-D5, phối hợp E1.
- **QA**: A4, E1, E2, hỗ trợ UAT.
- **DevOps/Infra**: E4, hỗ trợ E5.
- **Tech Lead/PM**: A1, điều phối dependency, quyết định cutover.

---

## 6) Acceptance criteria toàn chương trình migration

- 100% chức năng trong parity checklist v1 đạt.
- E2E critical flows pass.
- User hiện tại đăng nhập được với dữ liệu cũ.
- Upload ảnh hoạt động ổn định trên Windows.
- Có runbook rollback, đã diễn tập tối thiểu 1 lần.
