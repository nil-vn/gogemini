import { test, expect } from '@playwright/test';

test('login page renders', async ({ page }) => {
  await page.goto('/#/auth/login');
  await expect(page.getByPlaceholder('username/email')).toBeVisible();
  await expect(page.getByPlaceholder('password')).toBeVisible();
});

test('admin route is reachable (ui shell)', async ({ page }) => {
  await page.goto('/#/admin');
  await expect(page.getByText('GoGemini Admin')).toBeVisible();
});
