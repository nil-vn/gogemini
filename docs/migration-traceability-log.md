# Migration Traceability Log

Tất cả lượt implement/fix/review/go-live-check của migration phải append vào file này.

## Phase Gate Summary
- Phase 0: NOT STARTED
- Phase 1: NOT STARTED
- Phase 2: NOT STARTED
- Phase 3: NOT STARTED
- Phase 4: NOT STARTED
- Phase 5: NOT STARTED

---

## 2026-05-20 | INIT | REVIEW
- Owner: Codex
- Branch/Commit: <fill-after-first-task>
- Related PR: <fill-after-first-task>
- Scope summary:
  - Khởi tạo execution kit + reviewer checklist + traceability template.
  - Thiết lập chuẩn prompt để giao task A1→E5 tuần tự.
- Files changed:
  - docs/codex-migration-execution-kit.md
  - docs/migration-traceability-log.md
- Tests/Checks:
  - `markdown lint/manual review` => PASS (manual)
- DoD Mapping:
  - [x] Có prompt pack A1→E5 copy-paste dùng liền.
  - [x] Có checklist review cho reviewer.
  - [x] Có bảng traceability mẫu để append sau mỗi implement/fix.
- Related Outstanding Checklist Items:
  - Governance/Execution process support for phase 0→5 completion.
- Risks/Blockers:
  - Cần team tuân thủ kỷ luật append log sau mỗi task.
- Next Action:
  - Bắt đầu giao task A1 bằng prompt chuẩn trong execution kit.

## 2026-05-20 | A1 | IMPLEMENT
- Owner: Codex
- Branch/Commit: work @ (see commit history)
- Related PR: PR created via make_pr (pending link/id)
- Scope summary:
  - Tạo parity matrix v1 dạng table gồm feature, route hiện tại, expected behavior.
  - Khóa rõ phạm vi in-scope và out-of-scope cho v1 parity.
  - Đánh dấu item outstanding checklist tương ứng đã hoàn tất trong phạm vi A1.
- Files changed:
  - docs/migration-parity-v1-matrix.md
  - docs/migration-outstanding-production-checklist.md
  - docs/migration-traceability-log.md
- Tests/Checks:
  - `rg -n "A1|parity matrix|Feature parity matrix" docs/*.md` => PASS
  - `git diff -- docs/migration-parity-v1-matrix.md docs/migration-outstanding-production-checklist.md docs/migration-traceability-log.md` => PASS
- DoD Mapping:
  - [x] Có checklist parity v1 dạng table (feature, route hiện tại, expected behavior).
- Related Outstanding Checklist Items:
  - Phase 0 / Feature parity matrix theo từng module.
- Risks/Blockers:
  - Route legacy chi tiết từng endpoint CRUD/API cần A2 audit để khóa chính xác 100% theo code thực tế.
- Next Action:
  - Chuyển sang A2 để audit endpoint/payload/behavior và gắn lại route cụ thể mức API contract.

## 2026-05-20 | A2 | IMPLEMENT
- Owner: Codex
- Branch/Commit: work @ 34f43b4
- Related PR: PR created via make_pr (pending link/id)
- Scope summary:
  - Audit toàn bộ endpoint Flask hiện tại cho scope parity v1 (auth/dashboard/users/cars/customers/transactions/search/system/upload).
  - Chuẩn hóa baseline contract ở dạng Markdown spec (endpoint, method, auth, payload, response, behavior).
  - Đánh dấu outstanding item “API baseline spec” đã hoàn tất trong phạm vi A2.
- Files changed:
  - docs/migration-a2-api-contract-baseline.md
  - docs/migration-outstanding-production-checklist.md
  - docs/migration-traceability-log.md
- Tests/Checks:
  - `rg -n "@routes.route|login_required" app/admin/controllers app/homepage/controllers` => PASS
  - `rg -n "A2|API baseline spec|migration-a2-api-contract-baseline" docs/*.md` => PASS
- DoD Mapping:
  - [x] Có API contract baseline (OpenAPI draft hoặc Markdown spec).
- Related Outstanding Checklist Items:
  - Phase 0 / API baseline spec (OpenAPI/Swagger hoặc tài liệu tương đương) cho backend Go.
- Risks/Blockers:
  - Baseline hiện tại là server-rendered HTML flow, nên khi chuyển REST cần task phase 2 map lại response contract JSON mà không lệch nghiệp vụ.
- Next Action:
  - Tiếp tục A3 để freeze schema + ERD + data mapping rules.

## 2026-05-20 | A3 | IMPLEMENT
- Owner: Codex
- Branch/Commit: work @ (see commit history)
- Related PR: PR created via make_pr (pending link/id)
- Scope summary:
  - Freeze schema revision head và chuỗi migration làm baseline cho migration Go.
  - Tạo schema snapshot dạng SQL + ERD mermaid + data mapping rules từ legacy sang Go domain.
  - Đánh dấu checklist outstanding của A3 đã hoàn tất.
- Files changed:
  - docs/migration-a3-schema-freeze.md
  - docs/migration-outstanding-production-checklist.md
  - docs/migration-traceability-log.md
- Tests/Checks:
  - `rg -n "revision =|down_revision =" migrations/versions/*.py` => PASS
  - `rg -n "A3|schema|ERD|data mapping|migration-a3-schema-freeze" docs/*.md` => PASS
- DoD Mapping:
  - [x] Có file snapshot schema + ERD đính kèm docs.
- Related Outstanding Checklist Items:
  - Phase 0 / ERD + schema snapshot + data mapping rules từ legacy sang hệ mới.
- Risks/Blockers:
  - Kiểu dữ liệu ngày giờ hiện dùng string ở legacy, có rủi ro parse/validate ở Phase 2 nếu không chuẩn hóa format input.
- Next Action:
  - Chuyển sang A4 để tạo legacy smoke baseline report đối chiếu parity.

## 2026-05-20 | A4 | IMPLEMENT
- Owner: Codex
- Branch/Commit: work @ (pending commit)
- Related PR: PR created via make_pr (pending link/id)
- Scope summary:
  - Tạo legacy smoke test baseline cho các flow trọng yếu: login, CRUD module chính, search, upload.
  - Ghi baseline report để làm mốc đối chiếu parity các phase sau.
  - Đánh dấu outstanding checklist item Phase 0 liên quan legacy smoke baseline đã hoàn tất.
- Files changed:
  - tests/smoke/test_legacy_smoke.py
  - docs/migration-a4-legacy-smoke-baseline.md
  - docs/migration-outstanding-production-checklist.md
  - docs/migration-traceability-log.md
- Tests/Checks:
  - `PYTHONPATH=. pytest -q tests/smoke/test_legacy_smoke.py` => PASS
- DoD Mapping:
  - [x] Smoke tests chạy pass ổn định trên branch baseline.
- Related Outstanding Checklist Items:
  - Phase 0 / Legacy smoke baseline report để làm mốc đối chiếu kết quả parity.
- Risks/Blockers:
  - Còn warnings deprecated từ stack legacy (không chặn A4, sẽ xử lý ở phase hardening).
- Next Action:
  - Dùng baseline này làm mốc cho E2 API parity và E1/E2E critical flows.
