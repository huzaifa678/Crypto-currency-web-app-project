import { defineConfig } from '@playwright/test';

export default defineConfig({
  timeout: 30_000,
  retries: process.env.CI ? 2 : 0,

  testDir: 'e2e-tests',
  testMatch: '**/*.spec.ts',
  testIgnore: [
    '**/src/**',
    '**/*.test.ts',
    '**/*.test.tsx',
  ],

  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },

  webServer: {
    command: 'npm run dev',
    port: 5173,
    reuseExistingServer: !process.env.CI,
  },
});
