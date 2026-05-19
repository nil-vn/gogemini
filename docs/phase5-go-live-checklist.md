# Phase 5 Go-Live Checklist

## Before Cutover
- [ ] Build artifacts available (`bin/server.exe`, `frontend/dist`)
- [ ] DB backup tested
- [ ] Environment variables set (DB, CORS, UPLOAD_ROOT, auth secret)
- [ ] Migration dry-run completed
- [ ] Smoke test baseline captured

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
