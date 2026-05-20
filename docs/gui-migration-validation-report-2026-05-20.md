# GUI Migration Validation Report (Flask -> Svelte 5)

Date: 2026-05-20  
Scope: Frontend/UI migration validation only (backend logic intentionally out-of-scope)

## Executive Summary

Current Svelte 5 GUI is functionally wired to backend APIs, but visual migration is incomplete: there is no legacy stylesheet import pipeline, and the current UI renders as mostly unstyled semantic HTML + JSON text dumps.

## Validation Evidence

### 1) Svelte app bootstrap has no stylesheet import
- `frontend/src/main.ts` mounts `App.svelte` but imports no global CSS file.
- This means Vite bundle does not include any visual base style by default.

### 2) App-level UI is currently scaffold-style, not parity UI
- `frontend/src/App.svelte` uses minimal controls (`button`, `input`, `textarea`, `pre`) and JSON rendering for dashboard/search/module data.
- This is useful for API parity/debug, but not equivalent to previous Flask admin UX.

### 3) Legacy Flask UI had full style system available
- Legacy static CSS assets exist under `static/admin/css/*`.
- Flask templates still define structured admin shells and partials under `templates/admin/*`.
- These are the primary sources to migrate visual structure and design tokens/components.

### 4) Migration planning docs exist, but GUI parity closure is not done
- Repository already has migration planning/report docs.
- Frontend currently looks like phase API integration completed faster than phase visual parity.

## Root Causes (Likely)

1. **No CSS bridge in Svelte entrypoint**: legacy stylesheet not imported in frontend bundle.
2. **No layout component migration**: Flask base/header/sidebar/footer templates not re-implemented as Svelte layout.
3. **No component-level mapping**: module pages render raw objects rather than UX components (table cards/forms with consistent spacing, states, typography).
4. **No asset-path strategy**: no explicit approach documented in code for consuming `static/admin/*` assets from Vite/Svelte runtime.
5. **Parity acceptance criteria not encoded in FE tests**: current FE tests focus utility logic, not visual/DOM parity.

## WBS / Breakdown Tasks to Close GUI Gap

### WBS-1: Styling Foundation Recovery
- 1.1 Inventory legacy CSS dependency order from Flask base template (`style.css`, presets, plugins, custom overrides).
- 1.2 Define Svelte import strategy (global import in `main.ts` or app-level layout).
- 1.3 Validate icon/font/plugin CSS path resolution inside Vite.
- 1.4 Produce before/after visual baseline screenshot set (login, dashboard, list, detail/form).
- Deliverable: Svelte app renders with non-broken base styles and consistent typography/spacing.

### WBS-2: Layout Parity (Shell)
- 2.1 Map Flask layout partials: base, header, sidebar, footer, auth layout.
- 2.2 Rebuild as Svelte layout components with slot/content regions.
- 2.3 Port navigation states (active route, collapse behavior, accessibility labels).
- 2.4 Preserve responsive breakpoints and page container behavior.
- Deliverable: Visual shell parity across all admin routes.

### WBS-3: Screen-by-Screen UI Component Parity
- 3.1 Dashboard widgets/cards/charts parity.
- 3.2 CRUD list screens (users/cars/customers/transactions): table, filters, pagination UI, actions.
- 3.3 Detail/create/edit forms: field grouping, validation messages, button hierarchy.
- 3.4 Search and system/settings pages parity.
- 3.5 File upload UX parity (dropzone/state/preview/error).
- Deliverable: Feature pages match legacy UX behavior and look-and-feel within agreed tolerance.

### WBS-4: Asset + Interaction Integration
- 4.1 Move/reuse static assets (images/icons/fonts/js plugins) via deterministic Vite paths.
- 4.2 Replace jQuery/legacy JS behaviors with Svelte-native interactions (no hidden regressions).
- 4.3 Ensure i18n labels/translations preserve legacy wording.
- Deliverable: No missing icon/image/style and no broken interactive states.

### WBS-5: QA Gate for GUI Parity
- 5.1 Define parity checklist by route + state (default, loading, error, empty, populated).
- 5.2 Add Playwright visual smoke tests with baseline snapshots.
- 5.3 Add DOM-level assertions for key classes/landmarks/components.
- 5.4 Run accessibility pass (keyboard nav/focus contrast/ARIA landmarks).
- Deliverable: Repeatable sign-off gate proving migration quality.

### WBS-6: Cutover Readiness for Frontend
- 6.1 Browser matrix validation (Chrome/Edge/Safari/Firefox).
- 6.2 Performance sanity checks (LCP/CLS/JS payload budget).
- 6.3 Rollback/feature-flag strategy for UI fallback.
- Deliverable: Controlled production readiness for migrated GUI.

## Priority Sequence (Recommended)

1. WBS-1 (restore style pipeline)  
2. WBS-2 (layout shell)  
3. WBS-3 (page parity)  
4. WBS-4 (assets/interactions hardening)  
5. WBS-5 (QA gates)  
6. WBS-6 (final readiness)

## Immediate Next Actions (48h)

1. Add/verify global stylesheet import path in Svelte entrypoint.
2. Build `AdminLayout.svelte` from legacy base/header/sidebar/footer.
3. Convert one reference route (recommended: `dashboard`) to full visual parity first.
4. Freeze a screenshot baseline and open parity checklist per route.

## Definition of Done (GUI Migration)

- No blank/unstyled views on any admin route.
- 100% critical admin routes have migrated layout/component parity.
- Visual regression suite green against approved baselines.
- Accessibility smoke pass completed with no critical issues.
- Product owner sign-off for UI parity.
