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
- [ ] Hoàn thiện đầy đủ create/read/update/delete cho users/cars/customers/transactions.
- [ ] Chuẩn hóa filter/sort/pagination cho toàn bộ list APIs.
- [ ] Chuẩn response list (`items`, `total`, `page`, `page_size`, `sort`, `order`).
- [ ] Validation input + error handling nhất quán cho query/body params.

### Search, dashboard, system, upload
- [ ] Chuẩn hóa search behavior tương đương legacy.
- [ ] Dashboard metrics đúng định nghĩa nghiệp vụ (không chỉ count thô nếu chưa đủ).
- [ ] System settings có validation key/value + authorization đầy đủ.
- [ ] Upload hardening: MIME/ext whitelist, size limits, filename sanitization, chống path traversal, storage strategy production.

### Test coverage cho API
- [ ] Unit tests cho service/repo logic quan trọng.
- [ ] Integration tests cho auth/CRUD/search/system/upload.
- [ ] Contract/parity tests để khóa hành vi API.

## 4) Phase 3 - Frontend parity còn thiếu

- [ ] Hoàn thiện flow CRUD đầy đủ cho từng module (list/filter/sort/pagination/create/edit/detail/delete confirm).
- [ ] Đồng bộ validation FE-BE và xử lý UX states (loading/empty/error/retry).
- [ ] Hoàn thiện i18n (tách dictionary, loại bỏ hardcode chính, fallback language).
- [ ] Bổ sung accessibility nền tảng (keyboard nav, label semantics, focus/contrast).
- [ ] Bổ sung test cho frontend (unit + integration flow tests).

## 5) Phase 4 - QA/Performance/Security chưa đạt bằng chứng

- [ ] Nâng Playwright từ smoke lên critical flows đầy đủ.
- [ ] Chạy test ổn định trong CI, có artifacts khi fail (trace/video/screenshot).
- [ ] Có parity suite đối chiếu legacy vs Go cho từng module.
- [ ] Có performance report (p50/p95, throughput, memory/cpu) so với SLO.
- [ ] Có security checklist thực thi: dependency scan, auth/session/cors/csrf checks, injection/xss/upload abuse tests.

## 6) Phase 5 - Cutover/Go-live chưa hoàn tất thực thi

- [ ] Dry-run cutover trên staging theo playbook production.
- [ ] Rollback drill thực tế và đo thời gian phục hồi.
- [ ] Backup/restore DB đã kiểm chứng end-to-end.
- [ ] Go-live checklist được tick bằng evidence thực tế (không chỉ template).
- [ ] Hypercare plan 24-72h sau cutover có owner, dashboard và rollback trigger rõ ràng.

## 7) DevOps/Observability còn thiếu

- [ ] CI/CD pipeline hoàn chỉnh: lint, tests, build artifact, release process, changelog/versioning.
- [ ] Database migration strategy an toàn (forward/backward compatibility, rollback plan).
- [ ] Observability đầy đủ: metrics, tracing, alerts, dashboards.

## 8) Điều kiện đóng migration 100%

Chỉ xác nhận hoàn tất migration khi tất cả điều kiện sau đều đạt:

- [ ] Tất cả task phase 0->5 có bằng chứng hoàn thành.
- [ ] API parity report pass cho các luồng nghiệp vụ trọng yếu.
- [ ] E2E critical flows pass ổn định trong CI.
- [ ] Performance đạt ngưỡng SLO/SLA mục tiêu.
- [ ] Không còn issue security mức high/critical chưa xử lý.
- [ ] Cutover dry-run + rollback drill đều pass.
- [ ] Có sign-off đầy đủ từ Dev, QA, Ops và Business owner.
