# E3 Performance Baseline Report

## Scope
- Task ID: **E3 – Performance baseline**.
- Mục tiêu: tạo baseline đo **p50/p95 latency**, **throughput**, và **memory/cpu footprint** trước release candidate.

## Method
- Endpoint benchmark: `GET /healthz` (đại diện đường đi request core: router + middleware + DB ping).
- Tool: `scripts/perf_baseline.py` (Python, concurrent HTTP requests + lấy RSS/%CPU từ process Go bằng `ps`).
- Workload:
  - Total requests: `400`
  - Concurrency: `40`
- Runtime command:
  - `go run ./cmd/server`
  - `python3 scripts/perf_baseline.py --pid <server_pid> --url http://127.0.0.1:8080/healthz --requests 400 --concurrency 40`

## Environment
- Date (UTC): **2026-05-20**
- OS: Linux container (CI-like dev environment)
- Go service mode: default local config

## Results (Baseline)
- p50 latency: **21.88 ms**
- p95 latency: **71.63 ms**
- p99 latency: **104.72 ms**
- Throughput: **771.70 req/s**
- Memory RSS: **avg 27.81 MB**, **peak 27.81 MB**
- CPU: **avg 6.05%**

## Notes
- Kết quả là baseline tương đối cho môi trường hiện tại; số liệu production/staging có thể khác theo hạ tầng, DB workload và network.
- Baseline này dùng làm mốc đối chiếu cho RC và go-live checklist (monitor p95 latency).
