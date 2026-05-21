import { test, expect } from '@playwright/test';

// Sanity budgets for CI/local readiness gate (non-lab, mocked backend-free path)
const LCP_BUDGET_MS = 2500;
const CLS_BUDGET = 0.1;

test('WBS-6.2 performance sanity: LCP/CLS budget on login route', async ({ page }) => {
  await page.goto('/#/auth/login');

  const metrics = await page.evaluate(() => {
    const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming | undefined;
    const lcpEntries = performance.getEntriesByType('largest-contentful-paint') as PerformanceEntry[];
    const lcp = lcpEntries.length > 0 ? lcpEntries[lcpEntries.length - 1].startTime : nav?.domContentLoadedEventEnd ?? 0;

    let cls = 0;
    for (const entry of performance.getEntriesByType('layout-shift') as any[]) {
      if (!entry.hadRecentInput) cls += entry.value;
    }

    return { lcp, cls };
  });

  expect(metrics.lcp, `LCP must be <= ${LCP_BUDGET_MS}ms, got ${metrics.lcp.toFixed(2)}ms`).toBeLessThanOrEqual(LCP_BUDGET_MS);
  expect(metrics.cls, `CLS must be <= ${CLS_BUDGET}, got ${metrics.cls.toFixed(4)}`).toBeLessThanOrEqual(CLS_BUDGET);
});
