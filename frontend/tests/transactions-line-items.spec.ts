import { test, expect, type Page } from '@playwright/test';

type ModuleKey = 'users' | 'cars' | 'customers' | 'transactions';
type RecordItem = { id: string; [key: string]: unknown };

async function mountMockApi(page: Page) {
  let loggedIn = false;
  const store: Record<ModuleKey, RecordItem[]> = {
    users: [{ id: 'u-1', login: 'seed-user', email: 'seed@example.com' }],
    cars: [{ id: 'c-1', name: 'Seed Car', plate_number: 'ABC-123' }],
    customers: [{ id: 'cu-1', name: 'Seed Customer', phone_number: '0900' }],
    transactions: [{ id: 't-1', customer_id: 'cu-1', car_id: 'c-1', amount: 1000, transaction_items: [] }]
  };

  const json = (body: unknown, status = 200) => ({
    status,
    contentType: 'application/json',
    body: JSON.stringify(body)
  });

  await page.route('http://localhost:8080/**', async (route) => {
    const req = route.request();
    const method = req.method();
    const url = new URL(req.url());
    const path = url.pathname;

    if (path === '/api/auth/login' && method === 'POST') {
      loggedIn = true;
      return route.fulfill(json({ ok: true }));
    }

    if (path.startsWith('/api/admin/') && !loggedIn) {
      return route.fulfill({ status: 401, body: 'Unauthorized' });
    }

    if (path === '/api/admin/dashboard' && method === 'GET') {
      return route.fulfill(json({ users: 1, cars: 1, customers: 1, transactions: store.transactions.length }));
    }

    const moduleMatch = path.match(/^\/api\/admin\/(users|cars|customers|transactions)(?:\/([^/]+))?$/);
    if (moduleMatch) {
      const module = moduleMatch[1] as ModuleKey;
      const id = moduleMatch[2];

      if (method === 'GET' && !id) {
        return route.fulfill(json({ items: store[module], total: store[module].length, page: 1, page_size: 10, sort: 'id', order: 'asc' }));
      }

      if (method === 'GET' && id) {
        const found = store[module].find((x) => x.id === id);
        return route.fulfill(found ? json(found) : { status: 404, body: 'Not found' });
      }

      if (method === 'POST' && !id) {
        const body = JSON.parse(req.postData() || '{}') as RecordItem;
        const created = { ...body, id: `t-${Date.now()}` };
        store[module].push(created);
        return route.fulfill(json(created, 201));
      }

      if (method === 'PUT' && id) {
        const body = JSON.parse(req.postData() || '{}') as RecordItem;
        const idx = store[module].findIndex((x) => x.id === id);
        if (idx < 0) return route.fulfill({ status: 404, body: 'Not found' });
        store[module][idx] = { ...store[module][idx], ...body, id };
        return route.fulfill(json(store[module][idx]));
      }
    }

    return route.fulfill({ status: 500, body: `Unhandled route: ${method} ${path}` });
  });
}

test('transactions multi-item create and update flow', async ({ page }) => {
  await mountMockApi(page);

  await page.goto('/#/auth/login');
  await page.getByPlaceholder('username/email').fill('admin');
  await page.getByPlaceholder('password').fill('secret');
  await page.getByRole('button', { name: 'Login' }).click();

  await page.getByRole('button', { name: 'transactions' }).click();
  await page.getByLabel('Customer ID *').fill('cu-1');
  await page.getByLabel('Car ID *').fill('c-1');

  await page.getByRole('button', { name: '+ Add Item' }).click();
  await page.getByRole('button', { name: '+ Add Item' }).click();

  const itemNameInputs = page.locator('table tbody tr td:nth-child(1) input');
  const itemPriceInputs = page.locator('table tbody tr td:nth-child(2) input');

  await itemNameInputs.nth(0).fill('Tint Film');
  await itemPriceInputs.nth(0).fill('1200');
  await itemNameInputs.nth(1).fill('Dash Cam');
  await itemPriceInputs.nth(1).fill('800');

  await page.getByRole('button', { name: 'Create' }).click();
  await expect(page.getByText('Record created')).toBeVisible();

  const createdRow = page.locator('tr', { hasText: 'cu-1' }).first();
  await createdRow.getByRole('button', { name: 'Detail' }).click();

  await page.locator('table tbody tr').nth(1).getByRole('button', { name: 'Remove' }).click();
  await page.getByRole('button', { name: '+ Add Item' }).click();

  await expect(itemNameInputs).toHaveCount(2);
  await itemNameInputs.nth(1).fill('Floor Mat');
  await itemPriceInputs.nth(1).fill('450');

  await page.getByRole('button', { name: 'Update' }).click();
  await expect(page.getByText('Record updated')).toBeVisible();

  await page.reload();
  await createdRow.getByRole('button', { name: 'Detail' }).click();

  await expect(itemNameInputs.nth(0)).toHaveValue('Tint Film');
  await expect(itemPriceInputs.nth(0)).toHaveValue('1200');
  await expect(itemNameInputs.nth(1)).toHaveValue('Floor Mat');
  await expect(itemPriceInputs.nth(1)).toHaveValue('450');
});
