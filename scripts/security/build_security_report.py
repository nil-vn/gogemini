#!/usr/bin/env python3
from pathlib import Path

out = Path('docs/reports/security-gate-report.md')
lines = [
    '# Security Gate Report',
    '',
    '## Scope (WBS-4.3)',
    '- Dependency scan: pip-audit + govulncheck',
    '- SAST: bandit + staticcheck',
    '- App security behavior tests: csrf/cors/session/injection/xss/upload abuse',
    '',
    '## Artifacts',
]
for p in sorted(Path('.').glob('security-*')):
    lines.append(f'- `{p}`')
out.write_text('\n'.join(lines)+"\n")
print(out)
