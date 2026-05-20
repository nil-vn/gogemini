import { test, type Page } from '@playwright/test';
import fs from 'node:fs';
import path from 'node:path';

type ModuleKey = 'users' | 'cars' | 'customers' | 'transactions';
type RecordItem = { id: string; [key: string]: unknown };

async function mountMockApi(page: Page) {
  let loggedIn = false;
  const store: Record<ModuleKey, RecordItem[]> = {
    users: [{ id: 'u-1', login: 'seed-user', email: 'seed@example.com' }],
    cars: [{ id: 'c-1', name: 'Seed Car', plate_number: 'ABC-123' }],
    customers: [{ id: 'cu-1', name: 'Seed Customer', phone_number: '0900' }],
    transactions: [{ id: 't-1', customer_id: 'cu-1', car_id: 'c-1', amount: 1000 }]
  };

  const json = (body: unknown, status = 200) => ({ status, contentType: 'application/json', body: JSON.stringify(body) });

  await page.route('http://localhost:8080/**', async (route) => {
    const req = route.request();
    const method = req.method();
    const url = new URL(req.url());
    const p = url.pathname;

    if (p === '/api/auth/login' && method === 'POST') { loggedIn = true; return route.fulfill(json({ ok: true })); }
    if (p === '/api/auth/logout' && method === 'POST') { loggedIn = false; return route.fulfill(json({ ok: true })); }
    if (p.startsWith('/api/admin/') && !loggedIn) return route.fulfill({ status: 401, body: 'Unauthorized' });
    if (p === '/api/admin/dashboard' && method === 'GET') return route.fulfill(json({ users: 1, cars: 1, customers: 1, transactions: 1 }));
    if (p === '/api/admin/system' && method === 'GET') return route.fulfill(json({ currency: 'USD', theme: 'light', language: 'en' }));

    const m = p.match(/^\/api\/admin\/(users|cars|customers|transactions)(?:\/([^/]+))?$/);
    if (m && method === 'GET' && !m[2]) {
      const module = m[1] as ModuleKey;
      return route.fulfill(json({ items: store[module], total: store[module].length, page: 1, page_size: 10 }));
    }

    return route.fulfill(json({ ok: true }));
  });
}

function out(name: string) {
  const dir = path.resolve('tests/baseline');
  fs.mkdirSync(dir, { recursive: true });
  return path.join(dir, name);
}

test('capture WBS-1 visual baseline set', async ({ page }) => {
  await mountMockApi(page);

  await page.goto('/#/auth/login');
  await page.screenshot({ path: out('wbs1-login.png'), fullPage: true });

  await page.getByPlaceholder('username/email').fill('admin');
  await page.getByPlaceholder('password').fill('secret');
  await page.getByRole('button', { name: 'Login' }).click();
  await page.screenshot({ path: out('wbs1-dashboard.png'), fullPage: true });

  await page.getByRole('button', { name: 'users' }).click();
  await page.screenshot({ path: out('wbs1-list-users.png'), fullPage: true });

  await page.locator('textarea').fill('{"login":"visual-user","email":"visual@example.com"}');
  await page.screenshot({ path: out('wbs1-detail-form.png'), fullPage: true });
});
