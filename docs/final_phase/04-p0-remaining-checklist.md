# P0 Remaining Checklist (Tick-box)
Date: 2026-05-21

## P0-01 — Dedicated page model parity
- [x] Route kinds dedicated per module (`users/cars/customers/transactions/search/not-found`).
- [x] Dedicated page components created for Users/Cars/Customers/Transactions/Search/NotFound.
- [ ] Move remaining domain logic out of `App.svelte` into page-level modules/stores.
- [ ] Remove all generic-form fallback paths from operator primary flow.
- [x] Verify each route has page-specific acceptance evidence (list/new/detail).

## P0-02 — Remove generic JSON CRUD as default view
- [x] JSON payload editor is hidden behind debug flag (`VITE_UI_DEBUG_JSON_EDITOR`).
- [x] Confirm debug editor cannot be accessed in production profile.
- [x] Add test/guard for flag behavior in prod config.

## P0-03 — Transaction dynamic line-item parity
- [x] FE dynamic line-item add/remove UI exists.
- [x] Backend persists/loads/syncs transaction items (`transaction_item`).
- [x] CRUD parity test exists (`TestTransactionItemsCRUDParity`).
- [ ] Add UI-level automated flow test for multi-item create/update.

## P0-04 — Car/Customer purchase shortcut flows
- [x] Shortcut links from Car/Customer detail to Transaction new with query params.
- [x] Transaction prefill logic from `car_id` / `customer_id` query params.
- [ ] Add end-to-end test coverage for both shortcut paths.

## P0-05 — Search parity + grouped results
- [x] Dedicated `/admin/search` route and grouped module sections.
- [x] Deep-link from search results to detail routes.
- [ ] Replace generic JSON snippet preview with domain-specific row renderers.
- [ ] Align empty/no-result copy and UX parity with legacy.

## P0-06 — 404/Error parity
- [x] Dedicated NotFound page component wired for unknown routes.
- [ ] Standardize API error -> UI mapping taxonomy (401/403/validation/server).
- [ ] Add session-expiry/unauthorized redirect UX parity checks.
- [ ] Add parity evidence for 404 and error states.

## P0-07 — P0 Evidence Gate (must-pass before cutover)
- [ ] Playwright critical flow suite green:
  - [ ] login/logout
  - [ ] users CRUD
  - [ ] cars flow + purchase shortcut
  - [ ] customers flow + purchase shortcut
  - [ ] transactions line-items flow
  - [ ] search deep-link flow
- [ ] Visual baseline approved for critical pages.
- [ ] Operator UAT sign-off recorded.
- [ ] Final pass/fail matrix signed by FE/QA/Ops owners.

## Go/No-Go
- [ ] **GO** only when all unchecked boxes above are completed.
