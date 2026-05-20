#!/usr/bin/env python3
import argparse, json, sys
from pathlib import Path

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--input', required=True)
    ap.add_argument('--slo', required=True)
    ap.add_argument('--report', required=True)
    args = ap.parse_args()

    data = json.loads(Path(args.input).read_text())
    slo = json.loads(Path(args.slo).read_text())

    checks = {
        'p50_ms': data['latency_ms']['p50'] <= slo['p50_ms_max'],
        'p95_ms': data['latency_ms']['p95'] <= slo['p95_ms_max'],
        'throughput_rps': data['throughput_rps'] >= slo['throughput_rps_min'],
        'cpu_avg_pct': data['cpu_pct']['avg'] <= slo['cpu_avg_pct_max'],
        'rss_peak_mb': data['memory_rss_mb']['peak'] <= slo['rss_peak_mb_max'],
    }
    passed = all(checks.values())

    lines = [
        '# Performance Gate Report', '',
        f"- Status: {'PASS' if passed else 'FAIL'}",
        f"- Input: `{args.input}`",
        '', '## Metrics vs SLO'
    ]
    lines.append('| Metric | Value | SLO | Result |')
    lines.append('|---|---:|---:|---|')
    lines.append(f"| p50 (ms) | {data['latency_ms']['p50']:.2f} | <= {slo['p50_ms_max']} | {'PASS' if checks['p50_ms'] else 'FAIL'} |")
    lines.append(f"| p95 (ms) | {data['latency_ms']['p95']:.2f} | <= {slo['p95_ms_max']} | {'PASS' if checks['p95_ms'] else 'FAIL'} |")
    lines.append(f"| throughput (rps) | {data['throughput_rps']:.2f} | >= {slo['throughput_rps_min']} | {'PASS' if checks['throughput_rps'] else 'FAIL'} |")
    lines.append(f"| CPU avg (%) | {data['cpu_pct']['avg']:.2f} | <= {slo['cpu_avg_pct_max']} | {'PASS' if checks['cpu_avg_pct'] else 'FAIL'} |")
    lines.append(f"| RSS peak (MB) | {data['memory_rss_mb']['peak']:.2f} | <= {slo['rss_peak_mb_max']} | {'PASS' if checks['rss_peak_mb'] else 'FAIL'} |")
    Path(args.report).write_text('\n'.join(lines)+"\n")
    print('\n'.join(lines))
    if not passed:
        sys.exit(1)

if __name__ == '__main__':
    main()
