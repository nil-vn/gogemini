# Legacy Flask vs Go+Svelte Migration Gap Matrix (Production Closure)
Date: 2026-05-21

## Scope
Đối chiếu toàn bộ flow backend + frontend hiện có giữa:
- Legacy Flask (source of truth nghiệp vụ): `app/admin/controllers/*.py`, `templates/admin/*.html`.
- Migration target Go + Svelte5: `internal/http/admin.go`, `internal/repo/admin_repo.go`, `frontend/src/App.svelte` và components liên quan.

## Phương pháp
1. Rà route-level parity (list/new/detail/delete/action).  
2. Rà data-contract parity (field coverage, validation, relation loading).  
3. Rà interaction parity (JS behavior: tabs, sticky action, dynamic rows, plugins).  
4. Rà operational parity (auth, error states, i18n, upload, observability).

## Gap Matrix

## A) Users
- Legacy:
  - `/users`, `/user/new`, `/user/<id>`, `/user/<id>/delete` + validation duplicate username/email + password hashing on update if changed.
- Go+Svelte hiện tại:
  - Có CRUD API cơ bản và form frontend.
- Gap còn thiếu:
  1. Validation duplicate username/email chuẩn nghiệp vụ trước khi create/update (hiện chủ yếu DB/logic đơn giản).  
  2. Backend nhận `password_hash` trực tiếp thay vì raw password + server-side hash policy thống nhất.  
  3. UX parity tab Admin/Staff + badge/status behavior cần bám template cũ 100%.

## B) Cars
- Legacy:
  - Segmentation phong phú (available/awaiting_delivery/sold + situation groups).  
  - Upload multi-image, delete image API, purchase shortcut từ car detail.
- Go+Svelte hiện tại:
  - Có list/detail/create cơ bản, có upload endpoint theo module.
- Gap còn thiếu:
  1. Image lifecycle parity chưa đủ: delete image theo image_id + cleanup file + DB relation parity như Flask.
  2. Segmentation/UI cards chưa fully parity theo thông tin cũ (KPI cards + visual grouping).  
  3. Car detail action panel/sticky action parity chưa hoàn tất.

## C) Customers
- Legacy:
  - Active customers subset theo transaction relation.  
  - Upload multi-image, delete image API.  
  - Purchase shortcut từ customer detail.
- Go+Svelte hiện tại:
  - Có CRUD nền tảng.
- Gap còn thiếu:
  1. Active customer semantics theo relation thật chưa đầy đủ ở API list.  
  2. Customer image delete parity chưa bám legacy endpoint behavior.
  3. Detail layout/action parity (button group, status labels) chưa full.

## D) Transactions (critical)
- Legacy:
  - List theo All/Deposited/Paid + summary revenue/deposit.  
  - Create/detail với dynamic item rows (`item_name[]`, `item_price[]`).  
  - Car/customer purchase shortcut integration.
- Go+Svelte hiện tại:
  - Đã bổ sung persist/load `items` backend + dynamic editor frontend.
- Gap còn thiếu:
  1. Query/filter status semantics cần chuẩn hóa toàn bộ value mapping (PAID/DEPOSITED/wait2pay) xuyên suốt FE/API/DB.  
  2. Financial summary cần parity công thức legacy trên mọi page state.
  3. JS micro-parity (sticky action behavior + input formatting) chưa đủ.

## E) Search
- Legacy:
  - Search grouped results + deep-links entity detail.
- Go+Svelte hiện tại:
  - Có grouped results + links.
- Gap còn thiếu:
  1. Highlight/preview fields theo domain (không chỉ JSON slice).  
  2. Result ranking và empty-state copy parity.

## F) System settings
- Legacy:
  - Currency/theme/language update form.
- Go+Svelte hiện tại:
  - Có GET/PUT + admin role check.
- Gap còn thiếu:
  1. Persist/validation matrix parity đầy đủ cho i18n display ở mọi page.
  2. Theme behavior parity trên toàn bộ admin views.

## G) 404/Error/Auth UX
- Legacy:
  - 404 template rõ ràng, flash messaging.
- Go+Svelte hiện tại:
  - Có not-found route và notice stack.
- Gap còn thiếu:
  1. Error taxonomy consistency (API code/message -> UI render mapping).  
  2. Unauthorized/session-expired redirect UX parity.

## H) Frontend architecture parity risk
- Hiện vẫn là single large generic screen trong `App.svelte` cho nhiều route semantics.  
- Backlog P0 yêu cầu dedicated pages cho list/new/detail tất cả modules.  
- Kết luận: chưa đạt condition production parity nếu chưa tách page model theo legacy contract.

## I) Deployment readiness blockers (must close)
1. Route/page contract parity cho Users/Cars/Customers/Transactions/Search/404.  
2. Domain validation parity ở backend (không chỉ shape validation).  
3. Interaction parity (sticky actions, tabs, charts blocks, status badges).  
4. Test evidence parity: API + Playwright critical paths + visual baseline.

## Exit Criteria
Chỉ được coi là completed khi:
- 100% route contract map với legacy page contract.
- 100% critical flows pass automation + UAT operator sign-off.
- Không còn generic fallback làm primary path.
