# Phase 5 Go-Live Checklist

## Before Cutover
- [x] Build artifacts available (`bin/server.exe`, `frontend/dist`) _(evidence: `logs/cutover/staging-dryrun-20260520-0800.log`)_
- [x] DB backup tested _(evidence: `logs/cutover/db-backup-20260520-1110.log`)_
- [ ] Environment variables set (DB, CORS, UPLOAD_ROOT, auth secret)
- [x] Migration dry-run completed _(evidence: `logs/cutover/staging-dryrun-20260520-0800.log`)_
- [x] Smoke test baseline captured _(evidence: `logs/cutover/staging-smoke-20260520-0830.log`)_

## During Cutover
- [ ] Execute `scripts/cutover.ps1`
- [ ] Verify `/healthz` returns `status=ok`
- [ ] Verify login success
- [ ] Verify list endpoints for users/cars/customers/transactions
- [ ] Verify upload endpoint for cars/customers

## After Cutover (24–72h)
- [ ] Monitor logs and error rate
- [ ] Monitor p95 latency
- [ ] Validate no data regression
- [ ] Confirm no P1/P2 incidents

## Rollback Readiness
- [ ] `scripts/rollback.ps1` reviewed
- [ ] Latest `app.db.backup.*` available
- [ ] Traffic-switch owner assigned


## Evidence Notes
- Attach command outputs and artifacts under `logs/cutover/*.log` for every checklist line executed.
- Record incident status (P1/P2) for first 24h in `docs/migration-traceability-log.md` E5 entry updates.


## Execution Snapshot (2026-05-20 UTC)
- Dry-run: PASS (08:00-09:05).
- Rollback drill: PASS, RTO 11m40s, RPO <=5m.
- Backup/restore E2E: PASS (restored copy runnable).
- Evidence binder: `docs/wbs5-cutover-readiness-operability-binder-2026-05-20.md`.
