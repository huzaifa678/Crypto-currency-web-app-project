import { test, expect } from '@playwright/test';

test.describe('Login Page', () => {

  test.beforeEach(async ({ page }) => {
    await page.route('**/v1/auth/login**', async route => {
      const body = JSON.parse(route.request().postData() || '{}');

      if (body.email === 'huzaifa210@example.com' && body.password === 'password123') {
        return route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            user: {
              id: '1',
              username: 'Huzaifa',
              email: 'huzaifa210@example.com',
              role: 'user',
              is_verified: true,
            },
            access_token: 'fake-token',
          }),
        });
      }

      return route.fulfill({ status: 401 });
    });

    await page.route('https://accounts.google.com/**', route =>
      route.fulfill({ status: 200, body: '' })
    );

    await page.goto('http://localhost:5173/login');
  });

  test('renders login form', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /sign in to your account/i }))
      .toBeVisible();

    await expect(page.getByLabel('Email address')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByTestId('submit-button')).toBeVisible();
  });

  test('logs in successfully and redirects to dashboard', async ({ page }) => {
    await page.addInitScript(() => {
        localStorage.setItem('access_token', 'fake-token');
        localStorage.setItem(
        'user',
        JSON.stringify({
            id: '1',
            username: 'Huzaifa',
            role: 'user',
            is_verified: true,
        })
        );
    });

    await page.getByLabel('Email address').fill('huzaifa210@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByTestId('submit-button').click();

    await page.goto('http://localhost:5173/dashboard');

    await expect(page).toHaveURL(/dashboard/);
  });

  test('toggles password visibility', async ({ page }) => {
    const passwordInput = page.getByLabel('Password');

    await expect(passwordInput).toHaveAttribute('type', 'password');

    await page.getByTestId('toggle-password').click();

    await expect(passwordInput).toHaveAttribute('type', 'text');

    await page.getByTestId('toggle-password').click();

    await expect(passwordInput).toHaveAttribute('type', 'password');
  });

});