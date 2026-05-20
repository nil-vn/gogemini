# A4 Legacy Smoke Baseline Report

Date: 2026-05-20  
Task: A4 - Baseline smoke tests (hệ cũ)

## Scope smoke covered
- Login flow (`/admin/login`)
- Users CRUD smoke (create/list/update/delete)
- Cars CRUD smoke (create/list) + upload ảnh trong create form
- Customers CRUD smoke (create/list)
- Transactions CRUD smoke (create/list)
- Search flow (`/admin/search?q=...`)

## Test command
- `PYTHONPATH=. pytest -q tests/smoke/test_legacy_smoke.py`

## Result
- PASS: `1 passed`.
- Warnings observed from legacy stack dependencies (SQLAlchemy legacy API, datetime utcnow deprecation) nhưng không làm fail smoke baseline.

## Notes
- Test setup tự tạo DB baseline local tại `instance/db.sqlite3` trong môi trường test.
- CSRF được tắt trong chế độ test để ưu tiên xác thực behavior smoke của flow nghiệp vụ baseline.
