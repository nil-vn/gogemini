<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from './lib/api';
  import type { ModuleKey, Settings } from './lib/types';
  import LoginForm from './components/LoginForm.svelte';
  import ModuleTable from './components/ModuleTable.svelte';
  import SettingsForm from './components/SettingsForm.svelte';
  import UploadForm from './components/UploadForm.svelte';

  const modules: ModuleKey[] = ['users', 'cars', 'customers', 'transactions'];
  const moduleSet = new Set<ModuleKey>(modules);

  type AppRoute =
    | { kind: 'home' }
    | { kind: 'login' }
    | { kind: 'admin'; module: ModuleKey };

  let route: AppRoute = { kind: 'home' };
  let dashboard: Record<string, number> | null = null;
  let items: Record<ModuleKey, any[]> = { users: [], cars: [], customers: [], transactions: [] };
  let settings: Settings = { currency: 'JPY', theme: 'dark', language: 'vi' };
  let searchQ = '';
  let searchResult: Record<string, any[]> | null = null;
  let error = '';
  let message = '';
  let isAuthenticated = false;

  function parseRoute(hash: string): AppRoute {
    const normalized = (hash.replace('#', '') || '/').replace(/\/+$/, '') || '/';
    if (normalized === '/auth/login') return { kind: 'login' };
    if (normalized.startsWith('/admin')) {
      const maybeModule = normalized.split('/')[2] as ModuleKey | undefined;
      return { kind: 'admin', module: maybeModule && moduleSet.has(maybeModule) ? maybeModule : 'users' };
    }
    return { kind: 'home' };
  }

  function go(path: string) {
    window.location.hash = `#${path}`;
  }

  async function guarded<T>(fn: () => Promise<T>) {
    try {
      error = '';
      return await fn();
    } catch (e) {
      error = (e as Error).message;
      return null;
    }
  }

  async function checkAuth() {
    const probe = await guarded(() => apiFetch('/api/admin/dashboard')) as Record<string, number> | null;
    if (!probe) {
      isAuthenticated = false;
      dashboard = null;
      return false;
    }
    isAuthenticated = true;
    dashboard = probe;
    return true;
  }

  async function login(payload: { login: string; password: string }) {
    const ok = await guarded(() => apiFetch('/api/auth/login', { method: 'POST', body: JSON.stringify(payload) }));
    if (ok) {
      isAuthenticated = true;
      go('/admin/users');
      await bootstrapAdmin('users');
    }
  }

  async function logout() {
    await guarded(() => apiFetch('/api/auth/logout', { method: 'POST' }));
    isAuthenticated = false;
    dashboard = null;
    go('/auth/login');
  }

  async function loadModule(m: ModuleKey) {
    const res = await guarded(() => apiFetch(`/api/admin/${m}`)) as any;
    if (res?.items) items[m] = res.items;
  }

  async function loadSettings() {
    const res = await guarded(() => apiFetch('/api/admin/system')) as any;
    if (res) settings = { ...settings, ...res };
  }

  async function saveSettings(next: Settings) {
    const res = await guarded(() => apiFetch('/api/admin/system', { method: 'PUT', body: JSON.stringify(next) }));
    if (res === null) message = 'Settings saved';
  }

  async function runSearch() {
    searchResult = await guarded(() => apiFetch(`/api/admin/search?q=${encodeURIComponent(searchQ)}`)) as any;
  }

  async function upload(payload: { module: 'cars' | 'customers'; file: File }) {
    const form = new FormData();
    form.append('file', payload.file);
    const base = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080';
    const res = await fetch(`${base}/api/admin/upload/${payload.module}`, { method: 'POST', credentials: 'include', body: form });
    if (res.status === 401) {
      isAuthenticated = false;
      go('/auth/login');
      return;
    }
    if (!res.ok) {
      error = await res.text();
      return;
    }
    const data = await res.json();
    message = `Uploaded: ${data.path}`;
  }

  async function bootstrapAdmin(module: ModuleKey) {
    await Promise.all([loadModule(module), loadSettings()]);
  }

  async function syncRoute() {
    route = parseRoute(window.location.hash);
    if (route.kind !== 'admin') return;
    const authed = isAuthenticated || await checkAuth();
    if (!authed) {
      go('/auth/login');
      return;
    }
    await bootstrapAdmin(route.module);
  }

  onMount(() => {
    syncRoute();
    window.addEventListener('hashchange', syncRoute);
    return () => window.removeEventListener('hashchange', syncRoute);
  });
</script>

<main>
  <h1>GoGemini Admin</h1>
  {#if route.kind === 'login'}
    <LoginForm onSubmit={login} />
  {:else if route.kind === 'admin'}
    <nav>
      {#each modules as m}
        <button on:click={() => go(`/admin/${m}`)} disabled={route.module === m}>{m}</button>
      {/each}
      <button on:click={logout}>Logout</button>
    </nav>
    <section>
      <h3>Dashboard</h3>
      {#if dashboard}
        <ul>
          <li>Users: {dashboard.users}</li>
          <li>Cars: {dashboard.cars}</li>
          <li>Customers: {dashboard.customers}</li>
          <li>Transactions: {dashboard.transactions}</li>
        </ul>
      {/if}
    </section>
    <ModuleTable title={route.module} items={items[route.module]} />
    <section>
      <h3>Search</h3>
      <input bind:value={searchQ} placeholder="keyword" />
      <button on:click={runSearch}>Search</button>
      {#if searchResult}
        <pre>{JSON.stringify(searchResult, null, 2)}</pre>
      {/if}
    </section>
    <SettingsForm {settings} onSave={saveSettings} />
    <UploadForm onUpload={upload} />
  {:else}
    <p><a href="#/auth/login">Login</a> | <a href="#/admin/users">Admin</a></p>
  {/if}
  {#if message}<p style="color:green">{message}</p>{/if}
  {#if error}<p style="color:red">{error}</p>{/if}
</main>
