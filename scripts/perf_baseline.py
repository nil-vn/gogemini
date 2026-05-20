#!/usr/bin/env python3
import argparse
import concurrent.futures
import statistics
import subprocess
import time
import urllib.request


def fetch(url: str, timeout: float) -> float:
    start = time.perf_counter()
    with urllib.request.urlopen(url, timeout=timeout) as resp:
        resp.read()
    return (time.perf_counter() - start) * 1000.0


def percentile(sorted_values, p):
    if not sorted_values:
        return 0.0
    k = (len(sorted_values) - 1) * p
    f = int(k)
    c = min(f + 1, len(sorted_values) - 1)
    if f == c:
        return sorted_values[f]
    return sorted_values[f] * (c - k) + sorted_values[c] * (k - f)


def sample_process(pid: int):
    out = subprocess.check_output([
        "ps", "-p", str(pid), "-o", "rss=,%cpu="
    ], text=True).strip()
    if not out:
        return 0, 0.0
    rss_kb, cpu_pct = out.split()
    return int(rss_kb), float(cpu_pct)


def main():
    ap = argparse.ArgumentParser(description="Simple HTTP perf baseline")
    ap.add_argument("--url", default="http://127.0.0.1:8080/healthz")
    ap.add_argument("--requests", type=int, default=200)
    ap.add_argument("--concurrency", type=int, default=20)
    ap.add_argument("--pid", type=int, required=True)
    ap.add_argument("--timeout", type=float, default=5.0)
    args = ap.parse_args()

    latencies = []
    rss_samples = []
    cpu_samples = []

    start = time.perf_counter()
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.concurrency) as ex:
        futures = [ex.submit(fetch, args.url, args.timeout) for _ in range(args.requests)]

        for i, fut in enumerate(concurrent.futures.as_completed(futures), start=1):
            latencies.append(fut.result())
            if i % max(1, args.requests // 10) == 0:
                rss, cpu = sample_process(args.pid)
                rss_samples.append(rss)
                cpu_samples.append(cpu)

    duration = time.perf_counter() - start
    latencies.sort()

    p50 = percentile(latencies, 0.50)
    p95 = percentile(latencies, 0.95)
    p99 = percentile(latencies, 0.99)
    throughput = args.requests / duration if duration > 0 else 0

    rss_avg_mb = (statistics.mean(rss_samples) / 1024.0) if rss_samples else 0
    rss_peak_mb = (max(rss_samples) / 1024.0) if rss_samples else 0
    cpu_avg = statistics.mean(cpu_samples) if cpu_samples else 0

    print(f"url={args.url}")
    print(f"requests={args.requests} concurrency={args.concurrency}")
    print(f"latency_ms p50={p50:.2f} p95={p95:.2f} p99={p99:.2f}")
    print(f"throughput_rps={throughput:.2f}")
    print(f"memory_rss_mb avg={rss_avg_mb:.2f} peak={rss_peak_mb:.2f}")
    print(f"cpu_pct avg={cpu_avg:.2f}")


if __name__ == "__main__":
    main()
