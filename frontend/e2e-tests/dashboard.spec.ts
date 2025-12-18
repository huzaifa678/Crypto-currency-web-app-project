import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {

  test.beforeEach(async ({ page }) => {

    await page.addInitScript(() => {
      localStorage.setItem('user', JSON.stringify({
        id: '1',
        username: 'Huzaifa',
        email: 'huzaifa210@example.com',
        role: 'user',
        is_verified: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }));
      localStorage.setItem('access_token', 'fake-token');
    });

    await page.route('**/v1/**', route =>
      route.fulfill({ status: 200, body: '{}' })
    );

    await page.route('**/v1/stream_trades**', async route => {
      const encoder = new TextEncoder();
      const stream = new ReadableStream({
        start(controller) {
          controller.enqueue(
            encoder.encode(
              JSON.stringify({ result: { symbol: 'BTCUSDT', price: '43000' } }) + '\n'
            )
          );
        }
      });

      await route.fulfill({ status: 200, body: stream as any });
    });

    await page.goto('http://localhost:5173/dashboard');
  });

  test('shows welcome message', async ({ page }) => {
    await expect(
      page.getByText(/Welcome back, Huzaifa/i)
    ).toBeVisible();
  });

  test('renders stats cards', async ({ page }) => {
    await expect(page.getByText('Portfolio Value')).toBeVisible();
    await expect(page.getByText('24h Profit')).toBeVisible();
  });

  test('renders price chart', async ({ page }) => {
    await expect(
      page.getByRole('heading', { name: 'Price Chart' })
    ).toBeVisible();
  });

});
