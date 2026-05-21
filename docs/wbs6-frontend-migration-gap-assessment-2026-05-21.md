# WBS-6 Frontend Migration Gap Assessment (Flask Template ➜ Svelte 5)

Generated: 2026-05-21 (UTC)
Scope: verify UI parity between legacy Flask templates and new Svelte 5 frontend.

## 1) Verification method

Commands used:
- `rg --files -g '*.html'`
- `sed -n '1,260p' app/admin/controllers/main.py`
- `sed -n '1,260p' app/admin/controllers/user.py`
- `sed -n '1,260p' app/admin/controllers/car.py`
- `sed -n '1,260p' app/admin/controllers/customer.py`
- `sed -n '1,260p' app/admin/controllers/transaction.py`
- `sed -n '1,520p' frontend/src/App.svelte`

Comparison basis:
- Legacy route/page inventory = Flask controllers + Jinja templates.
- New route/page inventory = hash routes in `frontend/src/App.svelte` + loaded Svelte components.

## 2) Legacy page inventory (baseline)

### Admin/auth pages
- `/auth/login` → `templates/admin/login.html`
- `/admin/dashboard` → `templates/admin/dashboard.html`
- `/admin/system` → `templates/admin/system.html`

### User pages
- `/admin/users` → `templates/admin/users.html`
- `/admin/user/new` → `templates/admin/user_new.html`
- `/admin/user/<id>` → `templates/admin/user_detail.html`

### Car pages
- `/admin/cars` → `templates/admin/cars.html`
- `/admin/car/new` → `templates/admin/car_new.html`
- `/admin/car/<id>` → `templates/admin/car_detail.html`
- `/admin/car/<id>/purchase/new` (flow action from detail/new)

### Customer pages
- `/admin/customers` → `templates/admin/customers.html`
- `/admin/customer/new` → `templates/admin/customer_new.html`
- `/admin/customer/<id>` → `templates/admin/customer_detail.html`
- `/admin/customer/<id>/purchase/new` (flow action from detail/new)

### Transaction pages
- `/admin/transactions` → `templates/admin/transactions.html`
- `/admin/transaction/new` → `templates/admin/transaction_new.html`
- `/admin/transaction/<id>` → `templates/admin/transaction_detail.html`

### Other page
- `templates/admin/search.html` (legacy dedicated search page)

## 3) Current Svelte 5 page inventory

Current hash routes in `App.svelte`:
- `#/auth/login` (LoginForm)
- `#/admin/dashboard` (dashboard cards)
- `#/admin/system` (SettingsForm)
- `#/admin/users|cars|customers|transactions` (single generic CRUD screen with JSON editor + table)

Current shared components:
- `AdminLayout`, `AdminHeader`, `AdminSidebar`
- `ModuleTable`, `UploadForm`

## 4) Parity verdict

## Verdict: **Chưa migrate hoàn chỉnh về UI/UX parity**.

### What is already migrated (functional base)
- Core authentication screen exists.
- Dashboard route exists.
- System settings route exists.
- CRUD list/detail by API fetch exists for all 4 modules (users/cars/customers/transactions).
- Basic upload UI exists for cars/customers.

### Gaps vs legacy UI completeness (must close)
1. **Missing dedicated create pages** (`user_new`, `car_new`, `customer_new`, `transaction_new`) with field-first forms.
2. **Missing dedicated detail/edit pages** (`user_detail`, `car_detail`, `customer_detail`, `transaction_detail`) and their page-level context.
3. **No legacy-style advanced domain sections**:
   - Car page KPIs (available/awaiting/sold + refurbishment states).
   - Transaction revenue summaries (total/paid/deposit).
   - Customers segmentation (active vs all).
   - Users segmentation (admin vs staff).
4. **No transaction-items editor UI parity** (dynamic `item_name[]` / `item_price[]` behavior from legacy).
5. **No purchase flow UI from car/customer context** (`/car/<id>/purchase/new`, `/customer/<id>/purchase/new`).
6. **No dedicated search page parity** (`search.html` equivalent).
7. **Generic JSON editor UX differs strongly from legacy operator workflow**, increasing risk for production operations.
8. **Template-level reusable UI blocks parity not proven** (`_header`, `_sidebar`, flash/toast patterns, breadcrumbs, action shortcuts).
9. **404/error page parity not implemented in Svelte routing experience.**

## 5) Implementation WBS to close UI gaps

## WBS-FE-PARITY-01 — Information architecture & routing parity
- Define target route map matching legacy mental model (list/new/detail/search per module).
- Add explicit Svelte routes/screens:
  - `/admin/users`, `/admin/user/new`, `/admin/user/:id`
  - `/admin/cars`, `/admin/car/new`, `/admin/car/:id`
  - `/admin/customers`, `/admin/customer/new`, `/admin/customer/:id`
  - `/admin/transactions`, `/admin/transaction/new`, `/admin/transaction/:id`
  - `/admin/search`
- Add route guards, 404 route, and fallback UX.
- Deliverable: route parity matrix signed off by product/ops.

## WBS-FE-PARITY-02 — Shared UI system parity
- Rebuild legacy layout atoms in Svelte: sidebar groups, header actions, breadcrumbs, notices/toasts, form sections.
- Standardize tables/cards/forms to match legacy visual + interaction density.
- Add loading/empty/error variants for all page types.
- Deliverable: shared component library used by all parity pages.

## WBS-FE-PARITY-03 — Users module parity
- Implement Users List with admin/staff segmentation cards and filtered tables.
- Implement User Create form with validation parity.
- Implement User Detail/Edit form including password reset/update workflow parity.
- Deliverable: users UI flows equivalent to `users.html`, `user_new.html`, `user_detail.html`.

## WBS-FE-PARITY-04 — Cars module parity
- Implement Cars List with status/situation KPI blocks.
- Implement Car Create page with full form semantics and multi-image upload UX.
- Implement Car Detail page with edit, image management, and purchase shortcut actions.
- Deliverable: cars UI parity with `cars.html`, `car_new.html`, `car_detail.html`.

## WBS-FE-PARITY-05 — Customers module parity
- Implement Customers List with active-customer segmentation.
- Implement Customer Create page with multi-image upload.
- Implement Customer Detail with edit + purchase shortcut.
- Deliverable: customers UI parity with `customers.html`, `customer_new.html`, `customer_detail.html`.

## WBS-FE-PARITY-06 — Transactions module parity
- Implement Transactions List with revenue/deposit summary blocks.
- Implement Transaction Create page with dynamic line-item editor and prefill from query/context.
- Implement Transaction Detail page with editable line items and status transitions.
- Deliverable: transactions UI parity with `transactions.html`, `transaction_new.html`, `transaction_detail.html`.

## WBS-FE-PARITY-07 — Search & cross-module workflows
- Build dedicated `/admin/search` screen (filters, grouped results, quick actions).
- Ensure deep-link navigation from search results to detail pages.
- Deliverable: parity with legacy search workflow.

## WBS-FE-PARITY-08 — QA parity gates (must-pass)
- Build page-by-page parity checklist (controls, states, validations, actions).
- Add Playwright/Cypress smoke suite for all critical admin flows.
- Run side-by-side UAT against legacy app screenshots/video scripts.
- Deliverable: signed parity certificate “no missing UI”.

## 6) Suggested execution order (critical path)
1. WBS-FE-PARITY-01 (routing)
2. WBS-FE-PARITY-02 (shared UI primitives)
3. Parallel: WBS-FE-PARITY-03/04/05/06
4. WBS-FE-PARITY-07
5. WBS-FE-PARITY-08 (final gate)

## 7) Definition of Done (strict)
- 100% legacy admin pages have Svelte equivalents (list/new/detail/search/system/dashboard/login/404).
- No critical legacy action is only available via JSON editor.
- All module KPIs/segmentations from legacy are visible in Svelte pages.
- Cross-module purchase flows work from car/customer context.
- QA parity checklist has zero open “missing UI” items.
