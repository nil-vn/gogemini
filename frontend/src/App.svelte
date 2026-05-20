<script lang="ts">
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { apiFetch } from './lib/api';
  import type { DashboardMetrics, ModuleKey, ModuleRecord, Settings } from './lib/types';
  import { setLocale, t as tStore, type Locale } from './lib/i18n';
  import { validateRecord, validateSettings } from './lib/validation';
  import LoginForm from './components/LoginForm.svelte';
  import ModuleTable from './components/ModuleTable.svelte';
  import SettingsForm from './components/SettingsForm.svelte';
  import UploadForm from './components/UploadForm.svelte';

  const modules: ModuleKey[] = ['users', 'cars', 'customers', 'transactions'];
  const moduleSet = new Set<ModuleKey>(modules);

  type AppRoute =
    | { kind: 'home' }
    | { kind: 'login' }
    | { kind: 'dashboard' }
    | { kind: 'settings' }
    | { kind: 'admin'; module: ModuleKey };

  let route: AppRoute = { kind: 'home' };
  let dashboard: DashboardMetrics | null = null;
  let records: Record<ModuleKey, ModuleRecord[]> = { users: [], cars: [], customers: [], transactions: [] };
  let searchResults: Record<ModuleKey, ModuleRecord[]> = { users: [], cars: [], customers: [], transactions: [] };
  let settings: Settings = { currency: 'USD', theme: 'light', language: 'en' };
  let draft: Partial<ModuleRecord> = {};
  let selectedId = '';
  let filter = '';
  let globalSearchTerm = '';
  let draftText = '{}';
  let sort: 'asc' | 'desc' = 'asc';
  let page = 1;
  const pageSize = 10;
  let error = '';
  let message = '';
  let isLoading = false;
  let lastAction: (() => Promise<unknown>) | null = null;
  let isAuthenticated = false;
  $: activeModule = route.kind === 'admin' ? route.module : 'users';

  function tt(key: string, vars: Record<string, string | number> = {}) { return get(tStore)(key, vars); }

  function parseRoute(hash: string): AppRoute {
    const normalized = (hash.replace('#', '') || '/').replace(/\/+$/, '') || '/';
    if (normalized === '/auth/login') return { kind: 'login' };
    if (normalized === '/admin' || normalized === '/admin/dashboard') return { kind: 'dashboard' };
    if (normalized === '/admin/system') return { kind: 'settings' };
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
    try {
      isLoading = true;
      error = '';
      lastAction = fn;
      return await fn();
    } catch (e) {
      error = (e as Error).message;
      return null;
    } finally {
      isLoading = false;
    }
  }

  async function retryLastAction() {
    if (!lastAction) return;
    await guarded(lastAction);
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

  async function runGlobalSearch() {
    const term = globalSearchTerm.trim();
    if (!term) {
      searchResults = { users: [], cars: [], customers: [], transactions: [] };
      return;
    }
    const params = new URLSearchParams({ q: term });
    const res = await guarded(() => apiFetch(`/api/admin/search?${params.toString()}`)) as Record<ModuleKey, ModuleRecord[]> | null;
    if (!res) return;
    searchResults = { users: res.users ?? [], cars: res.cars ?? [], customers: res.customers ?? [], transactions: res.transactions ?? [] };
  }

  async function loadSettings() {
    const res = await guarded(() => apiFetch('/api/admin/system')) as Settings | null;
    if (!res) return;
    settings = res;
    setLocale(res.language);
  }

  async function saveSettings(payload: Settings) {
    const settingsErrors = validateSettings(payload);
    if (settingsErrors.length) { error = tt(settingsErrors[0]); return; }
    const res = await guarded(() => apiFetch('/api/admin/system', { method: 'PUT', body: JSON.stringify(payload) })) as Settings | null;
    if (!res) return;
    settings = res;
    setLocale(res.language);
    message = tt('recordUpdated');
  }

  async function bootstrapDashboard() {
    const res = await guarded(() => apiFetch('/api/admin/dashboard')) as DashboardMetrics | null;
    if (res) dashboard = res;
  }

  async function saveRecord(module: ModuleKey) {
    const validationErrors = validateRecord(module, draft);
    if (validationErrors.length) { error = tt(validationErrors[0]); return; }
    const id = selectedId || (draft as any)?.id;
    const method = id ? 'PUT' : 'POST';
    const path = id ? `/api/admin/${module}/${id}` : `/api/admin/${module}`;
    const payload = { ...draft };
    const result = await guarded(() => apiFetch(path, { method, body: JSON.stringify(payload) }));
    if (result) { message = id ? tt('recordUpdated') : tt('recordCreated'); clearEditor(); await loadModule(module); }
  }

  async function uploadImage(payload: { module: 'cars' | 'customers'; file: File }) {
    const form = new FormData();
    form.append('image', payload.file);
    const uploaded = await guarded(() => apiFetch(`/api/admin/upload/${payload.module}`, { method: 'POST', body: form }));
    if (uploaded) message = tt('uploadSuccess');
  }

  async function removeRecord(module: ModuleKey, id: string) {
    if (!confirm(tt('deleteConfirm', { module, id }))) return;
    const result = await guarded(() => apiFetch(`/api/admin/${module}/${id}`, { method: 'DELETE' }));
    if (result === null) { message = tt('recordDeleted'); if (selectedId === id) clearEditor(); await loadModule(module); }
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
    if (route.kind !== 'admin' && route.kind !== 'dashboard' && route.kind !== 'settings') return;
    const authed = isAuthenticated || await checkAuth();
    if (!authed) { go('/auth/login'); return; }
    if (route.kind === 'dashboard') await bootstrapDashboard();
    if (route.kind === 'admin') await bootstrapAdmin(route.module);
    if (route.kind === 'settings') await loadSettings();
  }

  syncDraftText();

  onMount(() => {
    syncRoute();
    window.addEventListener('hashchange', syncRoute);
    return () => window.removeEventListener('hashchange', syncRoute);
  });
</script>

<a href="#main-content" class="skip-link">{tt('a11ySkipToContent')}</a>
<main id="main-content" tabindex="-1">
  <h1>{tt('appTitle')}</h1>
  {#if route.kind === 'login'}
    <LoginForm onSubmit={login} t={tt} />
  {:else if route.kind === 'dashboard' || route.kind === 'admin' || route.kind === 'settings'}
    <nav aria-label={tt('a11yNavLabel')}>
      <button onclick={() => go('/admin/dashboard')} disabled={route.kind === 'dashboard'}>{tt('navDashboard')}</button>
      {#each modules as m}<button onclick={() => go(`/admin/${m}`)} disabled={route.kind === 'admin' && route.module === m}>{m}</button>{/each}
      <button onclick={() => go('/admin/system')} disabled={route.kind === 'settings'}>{tt('navSettings')}</button>
      <button onclick={logout}>{tt('navLogout')}</button>
    </nav>

    <section>
      <label for="global-search">{tt('searchLabel')}</label> <input id="global-search" bind:value={globalSearchTerm} placeholder={tt('searchPlaceholder')} />
      <button onclick={runGlobalSearch}>{tt('searchLabel')}</button>
      {#if globalSearchTerm.trim().length > 0}
        <h4>{tt('searchResults')}</h4>
        <pre>{JSON.stringify(searchResults, null, 2)}</pre>
      {/if}
    </section>

    {#if route.kind === 'dashboard'}
      <section><h3>{tt('dashboardTitle')}</h3>{#if dashboard}<pre>{JSON.stringify(dashboard, null, 2)}</pre>{/if}</section>
    {:else if route.kind === 'settings'}
      <SettingsForm settings={settings} onSave={saveSettings} t={tt} onLanguageChange={(next: Locale) => setLocale(next)} />
    {:else if route.kind === 'admin'}
      <section>
        <h3>{tt('managementTitle', { module: activeModule })}</h3>
        <label>{tt('filterLabel')} <input bind:value={filter} placeholder={tt('filterPlaceholder')} /></label>
        <label>{tt('sortById')}
          <select bind:value={sort}><option value="asc">asc</option><option value="desc">desc</option></select>
        </label>
        <button onclick={() => saveRecord(activeModule)}>{selectedId ? tt('update') : tt('create')}</button>
        <button onclick={clearEditor}>{tt('reset')}</button>
        <textarea rows="8" bind:value={draftText} onchange={() => {
          try {
            draft = JSON.parse(draftText);
          } catch (parseError) {
            error = (parseError as Error).message;
          }
        }}></textarea>
      </section>
      {#if activeModule === 'cars' || activeModule === 'customers'}
        <UploadForm onUpload={uploadImage} t={tt} />
      {/if}
      {@const view = visibleItems(activeModule)}
      <ModuleTable title={activeModule} items={view.items} total={view.total} page={page} pageSize={pageSize} onDetail={(id) => selectRecord(activeModule, id)} onDelete={(id) => removeRecord(activeModule, id)} t={tt} />
      <section>
        <button disabled={page<=1} onclick={() => page = page - 1}>{tt('prev')}</button>
        <span>{tt('page')} {page}</span>
        <button disabled={page*pageSize>=view.total} onclick={() => page = page + 1}>{tt('next')}</button>
      </section>
    {/if}
  {:else}
    <p><a href="#/auth/login">{tt('login')}</a> | <a href="#/admin/dashboard">{tt('admin')}</a></p>
  {/if}
  <section aria-live="polite" aria-label={tt('a11yStatusLabel')}>
    {#if isLoading}<p>{tt('loading')}</p>{/if}
    {#if message}<p style="color:green">{message}</p>{/if}
    {#if error}<p style="color:red">{tt('errorTitle')}: {error} <button onclick={retryLastAction}>{tt('retry')}</button></p>{/if}
  </section>
</main>

<style>
  .skip-link { position:absolute; left:-9999px; }
  .skip-link:focus { left: 8px; top:8px; background:#111; color:#fff; padding:8px; z-index:1000; }
</style>
