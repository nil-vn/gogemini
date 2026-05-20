# WBS-1 Foundation & Security Hardening Design Note

## Scope
Implements WBS-1 items 1.2-1.7 (excluding 1.1 by request).

## Decisions
- Structured JSON logs with request-id, PII redaction/hash, and 20% sampling for non-error requests.
- Standard error envelope: `{"error":{"code","message","request_id"}}` with taxonomy for auth/validation/dependency/internal.
- `/livez` for process liveness; `/readyz` for dependency readiness (DB ping).
- Graceful shutdown via `http.Server.Shutdown` timeout 15s and DB close.
- Auth hardening: secure httponly cookie, logout invalidation, session TTL reduction, unauthorized envelope.
- Login protection: lockout kept, plus audit trail via structured logs and error code signaling.
