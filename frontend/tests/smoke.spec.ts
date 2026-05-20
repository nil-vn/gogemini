import { test, expect, type Page } from '@playwright/test';

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

    if (path === '/api/auth/logout' && method === 'POST') {
      loggedIn = false;
      return route.fulfill(json({ ok: true }));
    }

    if (path.startsWith('/api/admin/') && !loggedIn) {
      return route.fulfill({ status: 401, body: 'Unauthorized' });
    }

    if (path === '/api/admin/dashboard' && method === 'GET') {
      return route.fulfill(json({ users: store.users.length, cars: store.cars.length, customers: store.customers.length, transactions: store.transactions.length }));
    }

    if (path === '/api/admin/search' && method === 'GET') {
      const q = (url.searchParams.get('q') ?? '').toLowerCase();
      const result = (Object.keys(store) as ModuleKey[]).reduce((acc, module) => {
        acc[module] = store[module].filter((item) => JSON.stringify(item).toLowerCase().includes(q));
        return acc;
      }, {} as Record<ModuleKey, RecordItem[]>);
      return route.fulfill(json(result));
    }

    if (path === '/api/admin/upload/cars' && method === 'POST') {
      return route.fulfill(json({ path: '/uploads/cars/test.png' }, 201));
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
        const created = { ...body, id: `${module.slice(0, 1)}-${Date.now()}` };
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

      if (method === 'DELETE' && id) {
        store[module] = store[module].filter((x) => x.id !== id);
        return route.fulfill({ status: 204, body: '' });
      }
    }

    if (path === '/api/admin/system' && method === 'GET') {
      return route.fulfill(json({ currency: 'USD', theme: 'light', language: 'en' }));
    }

    if (path === '/api/admin/system' && method === 'PUT') {
      return route.fulfill(json(JSON.parse(req.postData() || '{}')));
    }

    return route.fulfill({ status: 500, body: `Unhandled route: ${method} ${path}` });
  });
}

test.beforeEach(async ({ page }) => {
  await mountMockApi(page);
});

test('critical flow: login + users CRUD + search + upload + logout', async ({ page }) => {
  page.on('dialog', (dialog) => dialog.accept());

  await page.goto('/#/auth/login');
  await page.getByPlaceholder('username/email').fill('admin');
  await page.getByPlaceholder('password').fill('secret');
  await page.getByRole('button', { name: 'Login' }).click();
  await expect(page.getByText('Dashboard')).toBeVisible();

  await page.getByRole('button', { name: 'users' }).click();
  await expect(page.getByText('seed-user')).toBeVisible();

  const editor = page.locator('textarea');
  await editor.fill('{"login":"playwright-user","email":"pw@example.com"}');
  await page.getByRole('button', { name: 'Create' }).click();
  await expect(page.getByText('Record created')).toBeVisible();
  await expect(page.getByText('playwright-user')).toBeVisible();

  const newRow = page.locator('tr', { hasText: 'playwright-user' }).first();
  await newRow.getByRole('button', { name: 'Detail' }).click();
  await editor.fill('{"login":"playwright-user-updated","email":"pw@example.com"}');
  await page.getByRole('button', { name: 'Update' }).click();
  await expect(page.getByText('Record updated')).toBeVisible();
  await expect(page.getByText('playwright-user-updated')).toBeVisible();

  await page.getByPlaceholder('Search all modules').fill('updated');
  await page.getByRole('button', { name: 'Search' }).first().click();
  await expect(page.getByText('playwright-user-updated')).toBeVisible();

  await page.getByRole('button', { name: 'cars' }).click();
  await page.locator('input[type="file"]').setInputFiles({
    name: 'car.png',
    mimeType: 'image/png',
    buffer: Buffer.from('fake-image')
  });
  await page.getByRole('button', { name: 'Upload' }).click();
  await expect(page.getByText('Image uploaded successfully')).toBeVisible();

  await page.getByRole('button', { name: 'users' }).click();
  const deleteRow = page.locator('tr', { hasText: 'playwright-user-updated' }).first();
  await deleteRow.getByRole('button', { name: 'Delete' }).click();
  await expect(page.getByText('Record deleted')).toBeVisible();

  await page.getByRole('button', { name: 'Logout' }).click();
  await expect(page.getByPlaceholder('username/email')).toBeVisible();
});
