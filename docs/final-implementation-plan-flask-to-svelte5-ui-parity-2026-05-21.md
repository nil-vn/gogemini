# Final Implementation Plan — Migrate 100% Legacy Flask UI to Svelte 5 (No Scaffolding)

Generated: 2026-05-21 (UTC)
Status: Final execution blueprint (must-deliver)

## 1) Non-negotiable outcome
Migrate toàn bộ UI admin legacy từ Flask/Jinja sang Svelte 5 với **full parity theo page + workflow + interaction**, không chấp nhận trạng thái generic CRUD scaffold.

**Done means:**
- 100% legacy routes có Svelte equivalent chạy production.
- 0 critical operator flow bị downgrade thành JSON editor.
- Visual/interaction parity pass theo checklist và UAT.
- Cutover Svelte có rollback gate rõ ràng.

## 2) Scope baseline (must migrate)

### Auth/System
- `/auth/login`
- `/admin/dashboard`
- `/admin/system`
- 404/forbidden/error states

### Users
- `/admin/users`
- `/admin/user/new`
- `/admin/user/:id`

### Cars
- `/admin/cars`
- `/admin/car/new`
- `/admin/car/:id`
- `/admin/car/:id/purchase/new`

### Customers
- `/admin/customers`
- `/admin/customer/new`
- `/admin/customer/:id`
- `/admin/customer/:id/purchase/new`

### Transactions
- `/admin/transactions`
- `/admin/transaction/new`
- `/admin/transaction/:id`

### Search
- `/admin/search`

## 3) Technical principles (forced)
1. **Route-first parity:** route map phải mirror legacy mental model.
2. **Domain-first UI:** mỗi module có page riêng, không generic renderer làm primary UI.
3. **Workflow-first acceptance:** nghiệm thu theo thao tác nghiệp vụ end-to-end.
4. **No JSON-primary:** JSON editor chỉ debug mode, ẩn sau feature flag.
5. **Evidence-first QA:** mỗi page có screenshot baseline + scripted flow.

## 4) Deliverables by workstream

### WS1 — Routing & shell parity
- Hash/history route mapping đầy đủ list/new/detail/search/system/login/404.
- Guard auth + redirect logic parity.
- Breadcrumb/header/sidebar/action zone parity.

### WS2 — Shared UI primitives parity
- Section cards, notices, table wrappers, sticky action bars, empty/loading/error states.
- Reusable form sections để tái dùng cho users/cars/customers/transactions.

### WS3 — Users full parity
- List page có segmentation admin/staff.
- Create form có validation parity + role/status behavior.
- Detail/edit page có password update/reset behavior parity.

### WS4 — Cars full parity
- Tabs: all/available/awaiting/sold/refurbished/not_refurbished/refurbished_pending_cleaning.
- KPI cards + status badges + branch logo + table columns parity.
- Create/detail/edit + image management + purchase shortcut parity.

### WS5 — Customers full parity
- Tabs/segmentation active vs all.
- Create/detail/edit + image upload parity.
- Purchase shortcut flow parity.

### WS6 — Transactions full parity
- Tabs all/deposited/paid.
- KPI cards revenue/paid/deposit.
- Create/detail with **dynamic line items** `item_name[]` / `item_price[]`.
- Status transitions + amount computations parity.

### WS7 — Search & cross-module navigation parity
- Search form + grouped results by entity.
- Deep-link tới detail pages đúng context.
- Empty/no-result states parity.

### WS8 — Quality gates & go-live
- Page-by-page parity checklist (must-pass).
- Playwright smoke + critical flow tests.
- Side-by-side visual regression (legacy vs Svelte).
- UAT sign-off từ operator owners.

## 5) Concrete backlog (implementation tickets)

### P0 (must complete before cutover)
- Remove generic JSON CRUD as default view.
- Build dedicated Svelte pages for all list/new/detail routes.
- Implement transaction dynamic line-item editor parity.
- Implement car/customer purchase shortcut flows.
- Implement search page parity + grouped results.
- Implement 404/error parity.

### P1 (must complete before GA)
- Add charts/KPI blocks matching legacy information hierarchy.
- Add interaction micro-parity (sticky actions, badges, tab behaviors).
- Harden i18n labels and currency formatting parity.

### P2 (post-GA optimization)
- Performance polish and minor visual pixel refinements.
- Non-blocking UX enhancements.

## 6) Definition of Done (strict acceptance matrix)
A page is “Done” only if all pass:
1. Route exists and navigates correctly.
2. All primary controls/fields/actions exist.
3. Validation + error handling equivalent.
4. Domain summaries/tabs/KPIs equivalent.
5. Deep-link flows to related modules work.
6. Visual baseline approved.
7. Playwright scenario green.
8. Operator UAT sign-off recorded.

If any one item fails => page is **Not Done**.

## 7) Execution timeline (10 working days)
- **Day 1:** Route matrix freeze + acceptance checklist + ownership assignment.
- **Day 2-3:** Shared primitives + Users module parity.
- **Day 4-5:** Cars module parity.
- **Day 6-7:** Customers + Transactions parity (line items mandatory).
- **Day 8:** Search + cross-module flows + 404/error.
- **Day 9:** Visual regression + bug burn-down.
- **Day 10:** UAT replay, go/no-go review, cutover readiness sign-off.

## 8) Ownership model
- FE Lead: architecture + parity governance.
- Module owners: Users/Cars/Customers/Transactions/Search.
- QA owner: parity matrix + automation + evidence pack.
- Product/Ops owner: acceptance sign-off.

No shared ambiguous ownership is allowed.

## 9) Risk controls
- **Risk:** scope trượt về scaffold generic.
  - **Control:** PR template bắt buộc mapping route/page parity evidence.
- **Risk:** declare done không có UAT.
  - **Control:** merge gate cần chữ ký QA + Ops.
- **Risk:** behavior mismatch hidden behind visual similarity.
  - **Control:** flow-based test scripts cho critical actions.

## 10) Final commitment statement
Kể từ plan này, mọi claim “migrated” chỉ hợp lệ khi đạt full parity theo Definition of Done ở mục 6. Bất kỳ route/workflow legacy nào thiếu hoặc downgrade xuống scaffold đều được tính là migration chưa hoàn tất.
