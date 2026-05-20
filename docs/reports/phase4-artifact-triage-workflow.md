# Phase 4.2 - Artifact publish & failure triage workflow

## Artifact policy
- Playwright must publish `frontend/playwright-report` and `frontend/test-results` on every run (`if: always()`).
- Failure artifacts include trace/video/screenshot (Playwright defaults from retry/failure capture).
- Security gate publishes `security-*.json|txt` and `docs/reports/security-gate-report.md`.
- Perf gate publishes `perf-baseline.json`, `perf-server.log`, and `docs/reports/perf-gate-report.md`.

## Failure triage workflow
1. Open failed workflow run.
2. Download artifacts by job name: `playwright-report`, `security-report`, `perf-gate`.
3. Classify failure:
   - Test defect (flake/regression)
   - Security gate breach (new vuln/high severity)
   - Perf SLO breach
   - Infra/transient issue
4. Create issue with label: `phase4-gate` + subtype (`e2e`, `security`, `perf`, `infra`) and attach artifacts.
5. Apply fix, rerun pipeline until green streak reaches N consecutive runs.

## Green streak evidence
- Minimum threshold: **N = 3** consecutive successful `phase4-gates` runs.
- Evidence is recorded in `docs/reports/phase4-ci-green-evidence.md`.
