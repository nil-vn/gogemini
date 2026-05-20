# B1 - Go Project Scaffold Baseline

Tài liệu này ghi nhận scaffold backend Go cho task **B1** theo yêu cầu phase 1.

## Scope B1
- Thiết lập module Go (`go.mod`, `go.sum`) cho backend migration.
- Xác nhận cấu trúc thư mục chuẩn: `cmd/`, `internal/`, `migrations/`, `scripts/`.
- Chốt convention package theo layout hiện tại.

## Go module
- Module name: `gogemini`.
- Entrypoint server: `cmd/server/main.go`.

## Directory scaffold
- `cmd/server/`: application entrypoint.
- `internal/config/`: đọc/chuẩn hóa config runtime.
- `internal/http/`: router + handlers + middleware HTTP.
- `internal/service/`: service layer cho business logic.
- `internal/repo/`: database access/repositories.
- `internal/domain/`: domain model structs.
- `migrations/`: migration assets baseline.
- `scripts/`: automation scripts cho dev/build/migrate/run.

## Package conventions (baseline)
- Không import trực tiếp từ `cmd` vào `internal` theo chiều ngược.
- HTTP handlers chỉ orchestration request/response; business logic đi qua service layer.
- Repo layer chịu trách nhiệm truy cập DB; domain giữ model dùng chung.
- Mọi package compile được bởi `go test ./...` kể cả khi chưa có test files.

## Evidence
- `go test ./...` pass với toàn bộ package Go hiện tại.
