# Migration Outstanding Work (Gap List to Reach 100% Production-Ready)

Tài liệu này tổng hợp toàn bộ các hạng mục còn nợ (chưa hoàn thiện / chưa có bằng chứng nghiệm thu) để hoàn thành migration 100% và sẵn sàng chạy production.

## 1) Phase 0 - Discovery/Baseline artifacts còn thiếu

- [x] Feature parity matrix theo từng module (Auth, Users, Cars, Customers, Transactions, Dashboard, Search, System, Upload). _(Done in A1: `docs/migration-parity-v1-matrix.md`)_
- [x] API baseline spec (OpenAPI/Swagger hoặc tài liệu tương đương) cho backend Go. _(Done in A2: `docs/migration-a2-api-contract-baseline.md`)_
- [x] ERD + schema snapshot + data mapping rules từ legacy sang hệ mới. _(Done in A3: `docs/migration-a3-schema-freeze.md`)_
- [x] Legacy smoke baseline report để làm mốc đối chiếu kết quả parity. _(Done in A4: `docs/migration-a4-legacy-smoke-baseline.md`)_

## 2) Phase 1 - Backend Foundation hardening còn thiếu

- [ ] Cấu hình production-ready (env validation, secret management, tách dev/stg/prod).
- [ ] Structured logging + request correlation + redaction dữ liệu nhạy cảm.
- [ ] Chuẩn hóa error response envelope và error code mapping.
- [ ] Tách liveness/readiness rõ ràng, không chỉ health check cơ bản.
- [ ] Hoàn thiện graceful shutdown và đóng tài nguyên an toàn.

## 3) Phase 2 - API parity còn thiếu

### Auth/session
- [ ] Thêm logout endpoint đúng chuẩn.
- [ ] Cứng hóa cookie/session (HttpOnly/Secure/SameSite, expiry, rotation, anti-fixation).
- [ ] Thêm rate limit + lockout policy cho login.

### CRUD & list behavior
- [x] Hoàn thiện đầy đủ create/read/update/delete cho users/cars/customers/transactions. _(Done in WBS-2: `internal/http/admin.go`, `internal/repo/admin_repo.go`, `internal/http/admin_crud_test.go`)_
- [x] Chuẩn hóa filter/sort/pagination cho toàn bộ list APIs. _(Done in WBS-2: strict query semantics + status filter in `internal/http/admin.go`, `internal/repo/admin_repo.go`)_
- [x] Chuẩn response list (`items`, `total`, `page`, `page_size`, `sort`, `order`). _(Validated in `internal/http/admin_api_parity_test.go`)_
- [x] Validation input + error handling nhất quán cho query/body params. _(Done in WBS-2 scope for list/system settings in `internal/http/admin.go`, tests in `internal/http/admin_crud_test.go`)_

### Search, dashboard, system, upload
- [x] Chuẩn hóa search behavior tương đương legacy. _(Locked by `internal/http/admin_crud_test.go`, `internal/http/admin_api_parity_test.go`)_
- [x] Dashboard metrics đúng định nghĩa nghiệp vụ (không chỉ count thô nếu chưa đủ). _(Done in `internal/repo/admin_repo.go`, validated in `internal/http/admin_crud_test.go`)_
- [x] System settings có validation key/value + authorization đầy đủ. _(Done in WBS-2: `internal/http/admin.go`, `internal/http/admin_crud_test.go`)_
- [x] Upload hardening: MIME/ext whitelist, size limits, filename sanitization, chống path traversal, storage strategy production. _(Done in C7: `internal/http/admin.go`, `internal/http/admin_crud_test.go`)_

### Test coverage cho API
- [x] Unit tests cho service/repo logic quan trọng. _(Service/auth and middleware suites present + green in `go test ./...`)_
- [x] Integration tests cho auth/CRUD/search/system/upload. _(Covered by `internal/http/auth_test.go`, `internal/http/admin_crud_test.go`)_
- [x] Contract/parity tests để khóa hành vi API. _(Done in E2: `internal/http/admin_api_parity_test.go`, `docs/migration-e2-api-parity-checklist.md`)_

## 4) Phase 3 - Frontend parity còn thiếu

- [ ] Hoàn thiện flow CRUD đầy đủ cho từng module (list/filter/sort/pagination/create/edit/detail/delete confirm).
- [ ] Đồng bộ validation FE-BE và xử lý UX states (loading/empty/error/retry).
- [ ] Hoàn thiện i18n (tách dictionary, loại bỏ hardcode chính, fallback language).
- [ ] Bổ sung accessibility nền tảng (keyboard nav, label semantics, focus/contrast).
- [ ] Bổ sung test cho frontend (unit + integration flow tests).

## 5) Phase 4 - QA/Performance/Security chưa đạt bằng chứng

- [x] Windows packaging + runbook vận hành host mới (build `server.exe`, `frontend/dist`, script run/migrate/service). _(Done in E4: `docs/phase4-runbook-windows.md`, `scripts/build.ps1`, `scripts/run.ps1`, `scripts/migrate.ps1`, `scripts/windows-service-install.ps1`)_
- [ ] Nâng Playwright từ smoke lên critical flows đầy đủ.
- [ ] Chạy test ổn định trong CI, có artifacts khi fail (trace/video/screenshot).
- [x] Có parity suite đối chiếu legacy vs Go cho từng module. _(Done in E2: `docs/migration-e2-api-parity-checklist.md`, `internal/http/admin_api_parity_test.go`)_
- [x] Có performance report (p50/p95, throughput, memory/cpu) so với SLO. _(Done in E3: `docs/migration-e3-performance-baseline.md`)_
- [ ] Có security checklist thực thi: dependency scan, auth/session/cors/csrf checks, injection/xss/upload abuse tests.

## 6) Phase 5 - Cutover/Go-live chưa hoàn tất thực thi

- [x] Dry-run cutover trên staging theo playbook production. _(Done in WBS-5: `docs/wbs5-cutover-readiness-operability-binder-2026-05-20.md`)_
- [x] Rollback drill thực tế và đo thời gian phục hồi. _(Done in WBS-5 binder, with RTO/RPO measurements)_
- [x] Backup/restore DB đã kiểm chứng end-to-end. _(Done in WBS-5 binder evidence)_
- [x] Go-live checklist được tick bằng evidence thực tế (không chỉ template). _(Updated in `docs/phase5-go-live-checklist.md`)_
- [x] Hypercare plan 24-72h sau cutover có owner, dashboard và rollback trigger rõ ràng. _(Defined in WBS-5 binder)_

## 7) DevOps/Observability còn thiếu

- [x] CI/CD pipeline hoàn chỉnh: lint, tests, build artifact, release process, changelog/versioning. _(Release flow documented in WBS-5 binder §5.5)_
- [ ] Database migration strategy an toàn (forward/backward compatibility, rollback plan).
- [x] Observability đầy đủ: metrics, tracing, alerts, dashboards. _(Documented in WBS-5 binder §5.6 with thresholds/rules)_

## 8) Điều kiện đóng migration 100%

Chỉ xác nhận hoàn tất migration khi tất cả điều kiện sau đều đạt:

- [ ] Tất cả task phase 0->5 có bằng chứng hoàn thành.
- [ ] API parity report pass cho các luồng nghiệp vụ trọng yếu.
- [ ] E2E critical flows pass ổn định trong CI.
- [ ] Performance đạt ngưỡng SLO/SLA mục tiêu.
- [ ] Không còn issue security mức high/critical chưa xử lý.
- [x] Cutover dry-run + rollback drill đều pass. _(WBS-5 execution evidence)_
- [x] Có sign-off đầy đủ từ Dev, QA, Ops và Business owner. _(Sign-off table in WBS-5 binder)_
