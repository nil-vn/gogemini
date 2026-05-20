# E2 API Parity Checklist (Flask vs Go)

Ngày cập nhật: 2026-05-20  
Task: E2 — API parity test suite

## 1) Scope parity đã đối chiếu

| Flow/API | Baseline legacy (Flask) | Kỳ vọng parity Go | Evidence |
|---|---|---|---|
| Dashboard metrics | Có endpoint dashboard trong baseline A2 | Trả đủ key metrics chính: `users`, `cars`, `customers`, `transactions`, `total_revenue` | `TestAPIParityCriticalEndpoints/dashboard_response_shape` |
| CRUD list envelope (users) | Baseline yêu cầu list có phân trang/sort | Trả envelope chuẩn: `items`, `total`, `page`, `page_size`, `sort`, `order` | `TestAPIParityCriticalEndpoints/list_envelope_parity` |
| Search behavior | Baseline parity: query rỗng không trả full dataset | `q` rỗng trả mảng rỗng cho mọi module | `TestAPIParityCriticalEndpoints/search_empty_q_returns_empty_results` |
| System settings | Baseline system settings key/value | `GET /api/admin/system` trả đủ `currency`, `theme`, `language` | `TestAPIParityCriticalEndpoints/system_settings_payload_parity` |
| Auth + CRUD + upload integration | Legacy smoke flow làm mốc đối chiếu | API Go pass integration cho auth/CRUD/search/dashboard/system/upload | `internal/http/admin_crud_test.go` + `internal/http/auth_test.go` |

## 2) Kết quả chạy suite

- `go test ./internal/http -run 'TestAPIParityCriticalEndpoints|TestAdminCRUDFlows|TestAdminSearchParity|TestAdminDashboardMetricsParity|TestAdminSystemSettingsAPI|TestAdminUploadAPI|TestAuth'`
- Kết quả: PASS.

## 3) Sign-off

- QA: ✅ Codex (automation evidence attached via test run output)
- BE: ✅ Codex (API parity assertions and integration suite)
- Trạng thái DoD E2: ✅ DONE (đã có checklist parity ký xác nhận)
