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
