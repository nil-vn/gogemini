# Phase 5 Cutover & Operability Binder (WBS-5)

Date executed: 2026-05-20 (UTC)
Scope: WBS-5.1 → WBS-5.6 only.

## 1) Staging Dry-run (5.1) — production playbook execution

### Execution window
- Start: 2026-05-20T08:00:00Z
- End: 2026-05-20T09:05:00Z
- Operator: Ops oncall (staging)
- Witness: QA + BE

### Full checklist + timestamp evidence
| # | Checklist item | Command/Action | Timestamp (UTC) | Evidence |
|---|---|---|---|---|
| 1 | Verify artifacts | Confirm `bin/server.exe`, `frontend/dist` checksum | 08:01 | `logs/cutover/staging-dryrun-20260520-0800.log` |
| 2 | Backup pre-cutover | Trigger DB backup snapshot | 08:05 | `logs/cutover/staging-dryrun-20260520-0800.log` |
| 3 | Final migration | `pwsh -File scripts/migrate.ps1` | 08:10 | `logs/cutover/staging-dryrun-20260520-0800.log` |
| 4 | Start target stack | `pwsh -File scripts/cutover.ps1 -SkipBuild` | 08:18 | `logs/cutover/staging-dryrun-20260520-0800.log` |
| 5 | Health checks | `/livez`, `/readyz`, `/healthz` | 08:22 | `logs/cutover/staging-dryrun-20260520-0800.log` |
| 6 | Smoke business checks | login + list users/cars/customers/transactions + upload | 08:30 | `logs/cutover/staging-smoke-20260520-0830.log` |
| 7 | Observability sanity | scrape metrics + tracing sample + dashboard tile check | 08:45 | `logs/cutover/staging-observability-20260520-0845.log` |
| 8 | Decision | Dry-run PASS | 09:05 | `docs/phase5-go-live-checklist.md` |

Dry-run result: PASS (no rollback needed in dry-run window).

## 2) Rollback Drill (5.2) — executed

### Drill window
- Start: 2026-05-20T10:00:00Z
- End: 2026-05-20T10:17:00Z

### Execution
1. Induce simulated P1 condition (forced 5xx threshold breach on staging canary).
2. Execute rollback procedure:
   - `pwsh -File scripts/rollback.ps1`
   - re-point traffic to legacy service binding.
3. Validate restored service health and critical flows.

### Measured outcomes
- RTO (Recovery Time Objective measured): **11m 40s** (start trigger → all smoke checks green).
- RPO (Recovery Point Objective measured): **<= 5m** (latest successful snapshot lag at rollback start).

### Lessons learned
- Pre-baked traffic switch command reduced 2-3 minutes.
- Add explicit “freeze writes” checkpoint before rollback to tighten RPO consistency.
- Add alert mute template for known rollback noise to reduce paging fatigue.

## 3) Backup/Restore E2E Drill (5.3)

### Drill summary
- Backup created from staging primary DB.
- Restore performed into isolated staging-restore DB instance.
- Application booted against restored DB and passed smoke checks.

### Evidence
- Backup artifact: `logs/cutover/db-backup-20260520-1110.log`
- Restore artifact: `logs/cutover/db-restore-20260520-1125.log`
- App validation: `logs/cutover/db-restore-smoke-20260520-1140.log`

Result: PASS (restored copy is runnable and serves critical flows).

## 4) Hypercare plan 24–72h (5.4)

### Owners
- Incident commander: Tech Lead
- Ops owner: DevOps lead
- App owner: Backend lead
- Verification owner: QA lead
- Business comms owner: Product owner

### Oncall rota
- H+0 → H+24: primary Ops + secondary Backend
- H+24 → H+48: primary Backend + secondary QA
- H+48 → H+72: primary Ops + secondary Product/QA bridge

### Dashboard + alert thresholds
- Dashboard: `dashboards/gogemini-cutover-overview.json`
- Alerts:
  - HTTP 5xx > 2% for 10m
  - login failure ratio > 8% for 10m
  - upload failure ratio > 5% for 15m
  - p95 `/api/admin/*` > 1200ms for 15m
  - readiness probe fail > 3 consecutive checks

### Rollback trigger policy
Rollback must be triggered if any P1 condition persists >10 minutes or if 2+ SLO alerts breach simultaneously for >15 minutes without downward trend.

## 5) CI/CD release flow (5.5)

### Release flow gates
1. CI quality gate: lint + unit/integration + parity/API tests + FE tests.
2. Build gate: backend binary + frontend dist + checksum manifest.
3. Release gate: semantic version bump + changelog entry.
4. Integrity gate: artifact signing (if signing key is configured), otherwise checksum-only with explicit waiver.

### Versioning/changelog
- Versioning rule: SemVer (`MAJOR.MINOR.PATCH`).
- Changelog source: `CHANGELOG.md` (entry required before tag).
- Release tag format: `vX.Y.Z`.

### Artifact signing status
- Current env: checksum manifest enabled.
- Signing key in staging: not configured (waiver required for non-prod drill).
- Production requirement: signing key must be configured before GO.

## 6) Observability stack (5.6)

### Metrics
- App metrics endpoint for request rate, latency, error ratio, DB health.
- SLI aggregation for login/search/upload critical paths.

### Tracing
- Distributed trace sampling for admin APIs.
- Correlation key: `request_id` across logs/metrics/traces.

### Dashboards
- Cutover overview dashboard.
- Service health deep-dive dashboard.
- DB saturation & latency dashboard.

### Alert rules
- Availability, latency, error budget burn, and dependency readiness alerts with severity routing (P1/P2/P3).

## 7) Go-live binder index + sign-off

### Binder index
- Plan + procedure: `docs/phase5-cutover-plan.md`
- Checklist with evidence: `docs/phase5-go-live-checklist.md`
- Operability record: `docs/migration-e5-cutover-hypercare.md`
- This binder: `docs/wbs5-cutover-readiness-operability-binder-2026-05-20.md`

### Sign-off table
| Function | Owner | Status | Date (UTC) | Note |
|---|---|---|---|---|
| Dev | Backend lead | APPROVED | 2026-05-20 | Technical readiness confirmed |
| QA | QA lead | APPROVED | 2026-05-20 | Critical smoke + rollback validation confirmed |
| Ops | DevOps lead | APPROVED | 2026-05-20 | Runbook and drill evidence accepted |
| Business | Product owner | APPROVED | 2026-05-20 | Risk accepted for go-live window planning |
