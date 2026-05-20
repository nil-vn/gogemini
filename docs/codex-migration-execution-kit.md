# Codex Migration Execution Kit (A1 → E5)

Tài liệu này là **single source of truth** để điều phối Codex implement migration tuần tự, không bỏ sót requirement, và luôn lưu dấu vết (traceability) sau mỗi lần implement/fix.

---

## 1) Quy tắc vận hành bắt buộc

1. **Chỉ implement theo Task ID** trong phạm vi A1→E5, không làm lan scope.
2. Mỗi lần giao việc phải dùng prompt template trong mục 3.
3. Mỗi task phải bám:
   - `docs/migration-task-breakdown.md` (owner/dependency/DoD)
   - `docs/migration-outstanding-production-checklist.md` (gap production-ready)
4. Không được chuyển phase nếu chưa pass phase gate ở mục 4.
5. Sau mỗi lần implement/fix: **bắt buộc ghi traceability log** theo mẫu ở mục 6.

---

## 2) Danh sách task chuẩn A1→E5

### Phase 0
- A1. Khóa phạm vi v1 parity
- A2. Audit endpoint & behavior Flask hiện tại
- A3. Freeze schema DB + ERD
- A4. Baseline smoke tests (hệ cũ)

### Phase 1
- B1. Scaffold project Go
- B2. Config & environment strategy (Windows-first)
- B3. Middleware core
- B4. DB connect + healthcheck

### Phase 2
- C1. Domain model mapping
- C2. Repository CRUD (Users/Cars/Customers/Transactions)
- C3. Search API
- C4. Dashboard metrics API
- C5. Auth migration
- C6. System settings API
- C7. Upload API

### Phase 3
- D1. FE scaffold & toolchain
- D2. App shell + routing + auth guard
- D3. Module pages parity
- D4. Search + System settings + i18n
- D5. Upload UI + image preview

### Phase 4–5
- E1. E2E Playwright critical flows
- E2. API parity test suite
- E3. Performance baseline
- E4. Windows packaging & runbook
- E5. Cutover & hypercare

---

## 3) Prompt Pack (copy-paste dùng liền)

> Cách dùng: thay `<TASK_ID>` và `<TASK_TITLE>` tương ứng. Với mỗi task, dán prompt vào Codex, không rút gọn.

### 3.1 Prompt khung chuẩn cho mọi task

```text
Bạn là coding agent implement migration theo task tuần tự.

Task được giao:
- Task ID: <TASK_ID>
- Task title: <TASK_TITLE>

Nguồn yêu cầu bắt buộc (source of truth):
1) docs/codex-migration-execution-kit.md
2) docs/migration-task-breakdown.md
3) docs/migration-outstanding-production-checklist.md
4) docs/migration-plan-go-svelte5.md
5) docs/phase4-runbook-windows.md (nếu liên quan)
6) docs/phase5-cutover-plan.md + docs/phase5-go-live-checklist.md (nếu liên quan)

Quy tắc thực hiện:
- Chỉ làm đúng phạm vi task này, không tự mở rộng sang task khác.
- Trước khi code: tóm tắt plan ngắn + assumptions + dependencies.
- Khi code xong: map kết quả vào DoD của task.
- Bắt buộc chạy test/check phù hợp và báo kết quả.
- Nếu gặp blocker: dừng, nêu rõ blocker, đề xuất hướng xử lý.

Output bắt buộc:
1) Plan
2) Files changed
3) Implementation details
4) Test/check commands + kết quả
5) DoD mapping (đạt/chưa đạt từng mục)
6) Cập nhật traceability log vào docs/migration-traceability-log.md
```

### 3.2 Prompt implement task mới

```text
Thực hiện TASK <TASK_ID> - <TASK_TITLE>.

Yêu cầu:
- Bám đúng DoD của task trong docs/migration-task-breakdown.md.
- Kiểm tra các mục gap liên quan trong docs/migration-outstanding-production-checklist.md và xử lý phần thuộc phạm vi task.
- Không sửa ngoài phạm vi nếu không thật sự cần thiết.

Sau khi hoàn thành:
- Commit với message: "feat(migration): complete <TASK_ID> <short-title>"
- Ghi một dòng traceability vào docs/migration-traceability-log.md theo template.
- Báo cáo rõ phần nào đã done, phần nào còn mở.
```

### 3.3 Prompt fix/rework cho task đã làm

```text
Thực hiện FIX cho TASK <TASK_ID> - <TASK_TITLE>.

Context:
- Task này đã implement trước đó, cần rework theo feedback mới.

Yêu cầu:
- Đọc traceability log cũ của task trong docs/migration-traceability-log.md.
- Chỉ chỉnh phần liên quan feedback, tránh regression.
- Cập nhật test để chứng minh fix hoạt động.

Sau khi hoàn thành:
- Commit với message: "fix(migration): <TASK_ID> <short-fix-title>"
- Append log mới vào docs/migration-traceability-log.md (không ghi đè log cũ).
```

### 3.4 Prompt phase gate review

```text
Thực hiện Phase Gate Review cho <PHASE_NAME>.

Yêu cầu:
- Tổng hợp tất cả task trong phase này từ docs/migration-task-breakdown.md.
- Đối chiếu gap checklist liên quan trong docs/migration-outstanding-production-checklist.md.
- Đưa kết luận PASS/FAIL phase gate với lý do cụ thể.
- Nếu FAIL, liệt kê action items theo mức ưu tiên P0/P1.
- Cập nhật mục "Phase Gate Summary" trong docs/migration-traceability-log.md.
```

### 3.5 Prompt pre-go-live lock

```text
Thực hiện pre-go-live verification cho migration.

Yêu cầu:
- Đối chiếu đầy đủ docs/phase5-go-live-checklist.md.
- Đối chiếu điều kiện đóng migration trong docs/migration-outstanding-production-checklist.md (mục "Điều kiện đóng migration 100%").
- Kết luận READY / NOT READY.
- Nếu NOT READY: liệt kê chính xác checklist item chưa đạt và owner đề xuất.
- Cập nhật traceability log với nhãn "GO-LIVE-CHECK".
```

---

## 4) Phase Gate tiêu chuẩn (không được bỏ qua)

- **Gate Phase 0**: Có parity matrix + API baseline spec + ERD/schema snapshot + legacy smoke baseline evidence.
- **Gate Phase 1**: Backend foundation có env validation, logging, readiness/liveness, graceful shutdown.
- **Gate Phase 2**: API parity trọng yếu + auth/upload hardening + unit/integration/contract tests.
- **Gate Phase 3**: FE parity đầy đủ luồng CRUD/search/settings/upload + i18n + UX states + test.
- **Gate Phase 4**: CI critical flows ổn định + perf report + security checklist + packaging Windows.
- **Gate Phase 5**: Dry-run cutover pass + rollback drill pass + hypercare plan + go-live checklist evidence.

Nếu gate FAIL, **không chuyển phase**.

---

## 5) Checklist review cho reviewer

## A. Scope & Dependency
- [ ] Task ID khớp backlog A1→E5.
- [ ] Không làm lan scope ngoài task.
- [ ] Dependency của task đã thỏa.

## B. Implementation Quality
- [ ] Code/change đúng kiến trúc migration mục tiêu (Go API + Svelte 5).
- [ ] Không phá behavior legacy cần parity.
- [ ] Có xử lý error path và edge cases.

## C. Test & Evidence
- [ ] Có test/check commands liên quan.
- [ ] Kết quả test có thể tái lập.
- [ ] Có bằng chứng cho DoD (log/report/screenshot nếu cần).

## D. Security/Performance/Operations (nếu áp dụng)
- [ ] Auth/session/cookie policy đúng.
- [ ] Upload validation/hardening đúng.
- [ ] Có cân nhắc perf impact (p95/memory với task liên quan).
- [ ] Có cập nhật runbook/script khi chạm deploy/runtime.

## E. Traceability & Documentation
- [ ] Đã append log vào docs/migration-traceability-log.md.
- [ ] Có link commit/PR và file thay đổi rõ ràng.
- [ ] Có trạng thái DoD mapping: done/partial/block.

## F. Gate Decision
- [ ] PASS task.
- [ ] NEEDS FIX (ghi rõ action items).
- [ ] BLOCKED (ghi blocker + owner).

---

## 6) Bảng traceability mẫu (append sau mỗi implement/fix)

> Ghi vào file: `docs/migration-traceability-log.md`.

Template record:

```markdown
## <YYYY-MM-DD> | <TASK_ID> | <IMPLEMENT/FIX/REVIEW/GO-LIVE-CHECK>
- Owner: <name/team>
- Branch/Commit: <branch> @ <commit-sha>
- Related PR: <pr-link-or-id>
- Scope summary: <1-3 bullets>
- Files changed:
  - <file1>
  - <file2>
- Tests/Checks:
  - `<command 1>` => PASS/FAIL
  - `<command 2>` => PASS/FAIL
- DoD Mapping:
  - [x] <DoD item 1>
  - [ ] <DoD item 2> (nêu lý do nếu chưa đạt)
- Related Outstanding Checklist Items:
  - <section/item refs>
- Risks/Blockers:
  - <none or details>
- Next Action:
  - <what next>
```

---

## 7) Cách yêu cầu Codex từ bây giờ (chuẩn vận hành)

Khi yêu cầu task mới, luôn dùng cú pháp:

```text
Dùng docs/codex-migration-execution-kit.md làm chuẩn thực thi.
Hãy implement TASK <TASK_ID> - <TASK_TITLE> theo Prompt khung chuẩn + Prompt implement task mới.
Bắt buộc append traceability log vào docs/migration-traceability-log.md sau khi làm xong.
```

Với yêu cầu fix:

```text
Dùng docs/codex-migration-execution-kit.md làm chuẩn thực thi.
Hãy fix TASK <TASK_ID> - <TASK_TITLE> theo Prompt fix/rework.
Bắt buộc append traceability log mới vào docs/migration-traceability-log.md.
```

