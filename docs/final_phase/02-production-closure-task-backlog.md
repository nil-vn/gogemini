# Production Closure Backlog (Detailed Execution Spec)
Date: 2026-05-21

## Goal
Đưa hệ Go + Svelte 5 đạt migration parity 100% để deploy production, cả backend logic lẫn frontend UI/interaction.

## Tracking convention
- Priority: P0 blocker trước cutover, P1 trước GA, P2 hậu GA.
- Mỗi task có: Scope, Subtasks, Constraints, DoD, Evidence log.

---

## P0-01 — Dedicated page model parity (replace generic screen)
### Scope
Tách `frontend/src/App.svelte` thành page-level modules:
- `/admin/users`, `/admin/user/new`, `/admin/user/:id`
- `/admin/cars`, `/admin/car/new`, `/admin/car/:id`
- `/admin/customers`, `/admin/customer/new`, `/admin/customer/:id`
- `/admin/transactions`, `/admin/transaction/new`, `/admin/transaction/:id`
- `/admin/search`, `/404`

### Subtasks
1. Tạo router/page components riêng theo module.
2. Tạo shared layout + shared form primitives, nhưng giữ domain-specific sections riêng.
3. Map breadcrumbs/sidebar active states chuẩn legacy IA.
4. Xóa default generic JSON path khỏi operator flow.

### Constraints
- Không downgrade UX thành generic CRUD table.
- Tương thích deep-link và query prefill (`car_id`, `customer_id`).

### DoD
- 100% routes render page chuyên biệt.
- Không còn một page generic gánh mọi module primary flow.

### Evidence log template
- [ ] Screenshot từng page list/new/detail.
- [ ] Route matrix pass.
- [ ] Reviewer sign-off FE lead.

---

## P0-02 — Users parity closure
### Subtasks
1. Backend uniqueness validation username/email (create/update).
2. Password handling policy: nhận raw password, hash server-side, không yêu cầu client gửi hash.
3. FE form parity (role/status/password confirm/edit policy).
4. Admin/Staff segmentation counters + tabs parity.

### DoD
- Tạo/sửa user không thể trùng username/email.
- Password update behavior parity legacy.

### Evidence log
- [ ] API tests create/update conflict.
- [ ] UI tests add/edit user.

---

## P0-03 — Cars parity closure
### Subtasks
1. Complete segmentation filters theo status/situation legacy.
2. KPI/summary cards parity on cars list.
3. Image lifecycle parity: multi-upload + delete image endpoint + UI remove action.
4. Purchase shortcut flow from car detail -> transaction/new prefilled.

### Constraints
- Xử lý file upload an toàn (mime/ext/size).

### DoD
- Cars workflows list/new/detail/delete/upload/delete-image/purchase đều pass.

### Evidence log
- [ ] API test image delete.
- [ ] Playwright flow car -> add purchase.

---

## P0-04 — Customers parity closure
### Subtasks
1. Active customer grouping by transaction relation.
2. Image upload/delete parity như legacy.
3. Purchase shortcut from customer detail.
4. Detail page actions + notices parity.

### DoD
- Customer list/new/detail/delete + active segment + media actions pass.

---

## P0-05 — Transactions parity closure (critical)
### Subtasks
1. Hoàn thiện status normalization matrix (DB/API/UI).
2. Financial KPI parity (total revenue, paid revenue, deposited amount).
3. Dynamic line-item editor: add/remove, parse/format number, update replacement semantics.
4. Sticky action parity & form behavior parity với legacy.
5. Car/customer prefill flows + status quick actions.

### DoD
- End-to-end flows pass:
  - create tx with N items
  - update tx replace items
  - list segmented + correct summary

---

## P0-06 — Search parity closure
### Subtasks
1. Grouped results domain cards parity.
2. Result row template per entity (not JSON blob).
3. Clickthrough detail route parity.
4. Empty-state/no-result copy parity.

### DoD
- Search usable như legacy cho operator lookup.

---

## P0-07 — Error/404/Auth parity
### Subtasks
1. 401/403/session-expired UX mapping thống nhất.
2. 404 page parity visual + navigation actions.
3. Error codes to user-friendly messages/i18n.

### DoD
- Không còn dead-end/error thô từ API hiển thị trực tiếp.

---

## P0-08 — CSS/JS interaction parity pack
### Subtasks
1. Sticky action bars, tab switching behavior, badges.
2. Datepicker/choices/select behavior parity.
3. Table pagination/sort affordances gần legacy.
4. Responsive breakpoints parity các trang chính.

### DoD
- QA checklist interaction pass 100% trên browser matrix.

---

## P0-09 — Test & audit evidence gate
### Subtasks
1. API contract tests module-by-module.
2. Playwright critical flows:
   - login/logout
   - users CRUD
   - cars upload + purchase shortcut
   - customers upload + purchase shortcut
   - transactions items
   - search deep-link
3. Visual baseline capture per critical page.
4. UAT operator replay scripts + sign-off artifacts.

### DoD
- Gate chỉ PASS khi toàn bộ checklist green + evidence attached.

---

## Daily Reporting Log (to fill during execution)
| Date | Task ID | Owner | Status | Risks | Evidence Link | Next Action |
|---|---|---|---|---|---|---|
| 2026-05-21 | P0-05 | TBD | In progress | Status mapping divergence | TBD | Align enum mapping |

