<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from './lib/api';
  import type { DashboardMetrics, ModuleKey, ModuleRecord } from './lib/types';
  import LoginForm from './components/LoginForm.svelte';
  import ModuleTable from './components/ModuleTable.svelte';

  const modules: ModuleKey[] = ['users', 'cars', 'customers', 'transactions'];
  const moduleSet = new Set<ModuleKey>(modules);

  type AppRoute =
    | { kind: 'home' }
    | { kind: 'login' }
    | { kind: 'dashboard' }
    | { kind: 'admin'; module: ModuleKey };

  let route: AppRoute = { kind: 'home' };
  let dashboard: DashboardMetrics | null = null;
  let records: Record<ModuleKey, ModuleRecord[]> = { users: [], cars: [], customers: [], transactions: [] };
  let draft: Partial<ModuleRecord> = {};
  let selectedId = '';
  let filter = '';
  let draftText = '{}';
  let sort: 'asc' | 'desc' = 'asc';
  let page = 1;
  const pageSize = 10;
  let error = '';
  let message = '';
  let isAuthenticated = false;

  function parseRoute(hash: string): AppRoute {
    const normalized = (hash.replace('#', '') || '/').replace(/\/+$/, '') || '/';
    if (normalized === '/auth/login') return { kind: 'login' };
    if (normalized === '/admin' || normalized === '/admin/dashboard') return { kind: 'dashboard' };
    if (normalized.startsWith('/admin/')) {
      const maybeModule = normalized.split('/')[2] as ModuleKey | undefined;
      return { kind: 'admin', module: maybeModule && moduleSet.has(maybeModule) ? maybeModule : 'users' };
    }
    return { kind: 'home' };
  }

  function go(path: string) { window.location.hash = `#${path}`; }
  function syncDraftText() { draftText = JSON.stringify(draft, null, 2); }
  function clearEditor() { draft = {}; selectedId = ''; syncDraftText(); }

  async function guarded<T>(fn: () => Promise<T>) {
    try { error = ''; return await fn(); } catch (e) { error = (e as Error).message; return null; }
  }

  async function checkAuth() {
    const probe = await guarded(() => apiFetch('/api/admin/dashboard')) as DashboardMetrics | null;
    if (!probe) { isAuthenticated = false; dashboard = null; return false; }
    isAuthenticated = true; dashboard = probe; return true;
  }

  async function login(payload: { login: string; password: string }) {
    const ok = await guarded(() => apiFetch('/api/auth/login', { method: 'POST', body: JSON.stringify(payload) }));
    if (ok) { isAuthenticated = true; go('/admin/dashboard'); await bootstrapDashboard(); }
  }

  async function logout() {
    await guarded(() => apiFetch('/api/auth/logout', { method: 'POST' }));
    isAuthenticated = false; dashboard = null; go('/auth/login');
  }

  async function loadModule(module: ModuleKey) {
    const res = await guarded(() => apiFetch(`/api/admin/${module}`)) as { items?: ModuleRecord[] } | null;
    records[module] = res?.items ?? [];
  }

  async function bootstrapDashboard() {
    const res = await guarded(() => apiFetch('/api/admin/dashboard')) as DashboardMetrics | null;
    if (res) dashboard = res;
  }

  async function saveRecord(module: ModuleKey) {
    const id = selectedId || (draft as any)?.id;
    const method = id ? 'PUT' : 'POST';
    const path = id ? `/api/admin/${module}/${id}` : `/api/admin/${module}`;
    const payload = { ...draft };
    const result = await guarded(() => apiFetch(path, { method, body: JSON.stringify(payload) }));
    if (result) { message = id ? 'Record updated' : 'Record created'; clearEditor(); await loadModule(module); }
  }

  async function removeRecord(module: ModuleKey, id: string) {
    if (!confirm(`Delete ${module} #${id}?`)) return;
    const result = await guarded(() => apiFetch(`/api/admin/${module}/${id}`, { method: 'DELETE' }));
    if (result === null) { message = 'Record deleted'; if (selectedId === id) clearEditor(); await loadModule(module); }
  }

  async function selectRecord(module: ModuleKey, id: string) {
    const detail = await guarded(() => apiFetch(`/api/admin/${module}/${id}`)) as ModuleRecord | null;
    if (!detail) return;
    selectedId = id;
    draft = { ...detail };
    syncDraftText();
  }

  function visibleItems(module: ModuleKey) {
    const list = records[module] ?? [];
    const term = filter.trim().toLowerCase();
    const filtered = !term ? list : list.filter((item) => JSON.stringify(item).toLowerCase().includes(term));
    const sorted = [...filtered].sort((a, b) => {
      const aid = String(a.id ?? '');
      const bid = String(b.id ?? '');
      return sort === 'asc' ? aid.localeCompare(bid) : bid.localeCompare(aid);
    });
    const start = (page - 1) * pageSize;
    return { total: sorted.length, items: sorted.slice(start, start + pageSize) };
  }

  async function bootstrapAdmin(module: ModuleKey) { clearEditor(); page = 1; await loadModule(module); syncDraftText(); }

  async function syncRoute() {
    route = parseRoute(window.location.hash);
    if (route.kind !== 'admin' && route.kind !== 'dashboard') return;
    const authed = isAuthenticated || await checkAuth();
    if (!authed) { go('/auth/login'); return; }
    if (route.kind === 'dashboard') await bootstrapDashboard();
    if (route.kind === 'admin') await bootstrapAdmin(route.module);
  }

  syncDraftText();

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
  {:else if route.kind === 'dashboard' || route.kind === 'admin'}
    <nav>
      <button on:click={() => go('/admin/dashboard')} disabled={route.kind === 'dashboard'}>dashboard</button>
      {#each modules as m}<button on:click={() => go(`/admin/${m}`)} disabled={route.kind === 'admin' && route.module === m}>{m}</button>{/each}
      <button on:click={logout}>logout</button>
    </nav>

    {#if route.kind === 'dashboard'}
      <section><h3>Dashboard metrics</h3>{#if dashboard}<pre>{JSON.stringify(dashboard, null, 2)}</pre>{/if}</section>
    {:else}
      <section>
        <h3>{route.module} management</h3>
        <label>filter <input bind:value={filter} placeholder="filter" /></label>
        <label>sort by id
          <select bind:value={sort}><option value="asc">asc</option><option value="desc">desc</option></select>
        </label>
        <button on:click={() => saveRecord(route.module)}>{selectedId ? 'Update' : 'Create'}</button>
        <button on:click={clearEditor}>Reset</button>
        <textarea rows="8" bind:value={draftText} on:change={() => { try { draft = JSON.parse(draftText); } catch {} }}></textarea>
      </section>
      {@const view = visibleItems(route.module)}
      <ModuleTable title={route.module} items={view.items} total={view.total} page={page} pageSize={pageSize} onDetail={(id) => selectRecord(route.module, id)} onDelete={(id) => removeRecord(route.module, id)} />
      <section>
        <button disabled={page<=1} on:click={() => page = page - 1}>Prev</button>
        <span>page {page}</span>
        <button disabled={page*pageSize>=view.total} on:click={() => page = page + 1}>Next</button>
      </section>
    {/if}
  {:else}
    <p><a href="#/auth/login">Login</a> | <a href="#/admin/dashboard">Admin</a></p>
  {/if}
  {#if message}<p style="color:green">{message}</p>{/if}
  {#if error}<p style="color:red">{error}</p>{/if}
</main>
