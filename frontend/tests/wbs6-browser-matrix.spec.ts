import { test, expect } from '@playwright/test';

test('WBS-6.1 browser matrix smoke: login screen renders', async ({ page }) => {
  await page.goto('/#/auth/login');
  await expect(page.getByPlaceholder('username/email')).toBeVisible();
  await expect(page.getByPlaceholder('password')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Login' })).toBeVisible();
});
