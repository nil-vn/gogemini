import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  server: { port: 5173 },
  test: {
    include: ['src/**/*.test.ts'],
    environment: 'jsdom',
    coverage: { reporter: ['text', 'json-summary', 'html'] }
  }
});
