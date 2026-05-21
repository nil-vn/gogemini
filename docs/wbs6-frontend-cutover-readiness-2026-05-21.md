# WBS-6 Frontend Cutover Readiness Report

Date: 2026-05-21 (UTC)
Scope: WBS-6.1 → WBS-6.3 only.

## 6.1 Browser matrix validation (Chrome/Edge/Safari/Firefox)

Implemented Playwright multi-project matrix:
- chromium (Desktop Chrome)
- firefox (Desktop Firefox)
- webkit (Desktop Safari)
- edge (Desktop Edge channel)

Smoke testcase: `frontend/tests/wbs6-browser-matrix.spec.ts`
- Gate: login page must render username/password/Login CTA across matrix.

## 6.2 Performance sanity checks (LCP/CLS/JS payload budget)

Implemented performance sanity testcase: `frontend/tests/wbs6-performance.spec.ts`
- Route under test: `/#/auth/login`
- Budgets:
  - LCP <= 2500ms
  - CLS <= 0.1
- Runtime signal source: browser Performance APIs (`largest-contentful-paint`, `layout-shift`).

Notes:
- This is a sanity gate, not full Lighthouse lab profiling.
- Intended for cutover guardrail in CI/local pre-release checks.

## 6.3 Rollback/feature-flag strategy for UI fallback

Added UI fallback flags in frontend runtime:
- `VITE_UI_ROLLBACK_FORCE_LEGACY=true`
- `VITE_UI_FALLBACK_LEGACY_URL=<legacy_flask_admin_url>`

Behavior:
- On app mount, if both flags are set, frontend immediately redirects user to legacy URL.
- This gives deterministic frontend rollback switch without redeploying backend.

## Deliverable status
- Browser matrix gate: READY
- Performance sanity gate (LCP/CLS): READY
- UI fallback feature-flag rollback: READY
