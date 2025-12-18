/// <reference types="vitest/config" />
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/setupTests.ts', 
    css: true,
    include: [
      'src/**/*.test.{ts,tsx}',
    ],
    exclude: [
      'e2e-tests/**',
      '**/*.spec.ts',       
      '**/*.spec.tsx'
    ],
  }
})
