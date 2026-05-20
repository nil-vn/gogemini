# Phase 5 Cutover Plan

## Objective
Switch traffic from legacy Flask stack to Go + Svelte stack with controlled risk, fast verification, and defined rollback.

## Preconditions
- Phase 4 checks completed.
- `bin/server.exe` and `frontend/dist` are built.
- DB backup strategy verified.
- Smoke/E2E critical flows green in pre-prod.

## Cutover Steps
1. Create DB backup snapshot.
2. Apply final DB migration (`scripts/migrate.ps1`).
3. Start Go service (`scripts/run.ps1` or `scripts/cutover.ps1`).
4. Switch traffic:
   - Option A: IIS/Nginx route `/api` + `/` to Go/static
   - Option B: update service binding/DNS (low TTL window)
5. Execute post-switch smoke checks:
   - `GET /healthz`
   - Login
   - CRUD quick-path (1 create + 1 list)
   - Search
   - Upload image

## Monitoring Window (24–72h)
Track:
- HTTP 5xx rate
- Login failure ratio
- Upload failures
- p95 latency of `/api/admin/*`

## Rollback Criteria
Trigger rollback if any of below persists > 10 minutes:
- `/healthz` unstable
- 5xx spikes beyond agreed threshold
- Core admin flows unusable (login/CRUD/search)

## Rollback Actions
1. Stop Go service.
2. Restore latest DB backup (`scripts/rollback.ps1`).
3. Re-point traffic to legacy stack.
4. Announce incident + record root cause.


## WBS-5 Execution Addendum (2026-05-20)
- Staging dry-run executed with full checklist and timestamp evidence (`logs/cutover/staging-dryrun-20260520-0800.log`, `logs/cutover/staging-smoke-20260520-0830.log`).
- Rollback drill executed with measured RTO/RPO (RTO 11m40s, RPO <=5m).
- Backup/restore E2E drill completed with runnable restored copy validation.
- Hypercare/oncall/dashboard/alert thresholds and rollback triggers are codified in `docs/wbs5-cutover-readiness-operability-binder-2026-05-20.md`.
