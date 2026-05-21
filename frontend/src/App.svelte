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
  import AdminLayout from './components/layout/AdminLayout.svelte';
  import Breadcrumbs from './components/ui/Breadcrumbs.svelte';
  import NoticeStack from './components/ui/NoticeStack.svelte';
  import ContentState from './components/ui/ContentState.svelte';
  import SectionCard from './components/ui/SectionCard.svelte';

  const modules: ModuleKey[] = ['users', 'cars', 'customers', 'transactions'];
  const moduleSet = new Set<ModuleKey>(modules);

  const legacyFallbackUrl = import.meta.env.VITE_UI_FALLBACK_LEGACY_URL ?? '';
  const forceLegacyFallback = String(import.meta.env.VITE_UI_ROLLBACK_FORCE_LEGACY ?? 'false').toLowerCase() === 'true';

  type AppRoute =
    | { kind: 'home' }
    | { kind: 'login' }
    | { kind: 'dashboard' }
    | { kind: 'settings' }
    | { kind: 'search' }
    | { kind: 'module-list'; module: ModuleKey }
    | { kind: 'module-new'; module: ModuleKey }
    | { kind: 'module-detail'; module: ModuleKey; id: string }
    | { kind: 'not-found'; path: string };

  let route: AppRoute = { kind: 'home' };
  let dashboard: DashboardMetrics | null = null;
  let records: Record<ModuleKey, ModuleRecord[]> = { users: [], cars: [], customers: [], transactions: [] };
  let searchResults: Record<ModuleKey, ModuleRecord[]> = { users: [], cars: [], customers: [], transactions: [] };
  let settings: Settings = { currency: 'USD', theme: 'light', language: 'en' };
  let draft: Partial<ModuleRecord> = {};
  let userForm = { username: '', email: '', password: '', confirm_password: '', role: 'guest', status: 'Active' };
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
  $: activeModule = route.kind === 'module-list' || route.kind === 'module-new' || route.kind === 'module-detail' ? route.module : 'users';

  function tt(key: string, vars: Record<string, string | number> = {}) { return get(tStore)(key, vars); }

  const singularRouteToModule: Record<string, ModuleKey> = {
    user: 'users',
    car: 'cars',
    customer: 'customers',
    transaction: 'transactions'
  };

  function parseRoute(hash: string): AppRoute {
    const normalized = (hash.replace('#', '') || '/').replace(/\/+$/, '') || '/';
    if (normalized === '/auth/login') return { kind: 'login' };
    if (normalized === '/admin' || normalized === '/admin/dashboard') return { kind: 'dashboard' };
    if (normalized === '/admin/system') return { kind: 'settings' };
    if (normalized === '/admin/search') return { kind: 'search' };
    if (normalized.startsWith('/admin/')) {
      const segments = normalized.split('/').filter(Boolean);
      const moduleList = segments[1] as ModuleKey | undefined;
      if (moduleList && moduleSet.has(moduleList) && segments.length === 2) return { kind: 'module-list', module: moduleList };

      const singular = segments[1];
      const mappedModule = singularRouteToModule[singular];
      if (!mappedModule) return { kind: 'not-found', path: normalized };
      if (segments.length === 3 && segments[2] === 'new') return { kind: 'module-new', module: mappedModule };
      if (segments.length === 3 && segments[2]) return { kind: 'module-detail', module: mappedModule, id: decodeURIComponent(segments[2]) };
      return { kind: 'not-found', path: normalized };
    }
    return normalized === '/' ? { kind: 'home' } : { kind: 'not-found', path: normalized };
  }

  function go(path: string) { window.location.hash = `#${path}`; }
  function syncDraftText() { draftText = JSON.stringify(draft, null, 2); }
  function clearEditor() { draft = {}; selectedId = ''; userForm = toUserForm(); syncDraftText(); }

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
    if (module === 'users') userForm = toUserForm(detail);
    syncDraftText();
  }


  function toUserForm(input: Partial<ModuleRecord> = {}) {
    return {
      username: String(input.username ?? ''),
      email: String(input.email ?? ''),
      password: '',
      confirm_password: '',
      role: String(input.role ?? 'guest') || 'guest',
      status: String(input.status ?? 'Active') || 'Active'
    };
  }

  function userSegmentation() {
    const users = records.users ?? [];
    const adminUsers = users.filter((u) => String(u.role ?? '').toLowerCase() === 'admin');
    const staffUsers = users.filter((u) => String(u.role ?? '').toLowerCase() !== 'admin');
    return { users, adminUsers, staffUsers };
  }

  function syncUserDraftFromForm() {
    draft = { ...draft, username: userForm.username.trim(), email: userForm.email.trim(), role: userForm.role, status: userForm.status };
    if (userForm.password.trim()) draft.password = userForm.password;
    else delete (draft as any).password;
    draft.name = userForm.username.trim();
  }

  function validateUserForm(isEdit = false) {
    if (!userForm.username.trim()) return 'Username is required';
    if (userForm.password || userForm.confirm_password) {
      if (userForm.password !== userForm.confirm_password) return 'Password confirmation does not match';
    } else if (!isEdit) {
      return 'Password is required for new user';
    }
    return '';
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
    if (route.kind !== 'module-list' && route.kind !== 'module-new' && route.kind !== 'module-detail' && route.kind !== 'dashboard' && route.kind !== 'settings' && route.kind !== 'search') return;
    const authed = isAuthenticated || await checkAuth();
    if (!authed) { go('/auth/login'); return; }
    if (route.kind === 'dashboard') await bootstrapDashboard();
    if (route.kind === 'module-list' || route.kind === 'module-new' || route.kind === 'module-detail') {
      await bootstrapAdmin(route.module);
      if (route.kind === 'module-detail') await selectRecord(route.module, route.id);
    }
    if (route.kind === 'search') await runGlobalSearch();
    if (route.kind === 'settings') await loadSettings();
  }

  syncDraftText();

  onMount(() => {
    if (forceLegacyFallback && legacyFallbackUrl) {
      window.location.href = legacyFallbackUrl;
      return;
    }

    syncRoute();
    window.addEventListener('hashchange', syncRoute);
    return () => window.removeEventListener('hashchange', syncRoute);
  });
</script>

<main>
  <h1 class="d-none">{tt('appTitle')}</h1>
  {#if route.kind === 'login'}
    <LoginForm onSubmit={login} t={tt} />
  {:else if route.kind === 'dashboard' || route.kind === 'module-list' || route.kind === 'module-new' || route.kind === 'module-detail' || route.kind === 'settings' || route.kind === 'search'}
    <AdminLayout
      activePath={route.kind === 'dashboard' ? '/admin/dashboard' : route.kind === 'settings' ? '/admin/system' : route.kind === 'search' ? '/admin/search' : `/admin/${route.module}`}
      {globalSearchTerm}
      onNavigate={go}
      onSearch={runGlobalSearch}
      onSearchTermInput={(value) => globalSearchTerm = value}
      onThemeChange={(theme) => document.body.setAttribute('data-pc-theme', theme)}
      t={tt}
    >
      <Breadcrumbs onNavigate={go} items={[{ label: "Dashboard", path: "/admin/dashboard" }, ...(route.kind === "dashboard" ? [] : [{ label: route.kind === "settings" ? "System" : route.kind === "search" ? "Search" : activeModule }])]} />
      <NoticeStack isLoading={isLoading} message={message} error={error} onRetry={retryLastAction} t={tt} />
    {#if globalSearchTerm.trim().length > 0}
      <SectionCard title={tt('searchResults')}>
        <div class="row g-3">
          {#each modules as module}
            <div class="col-md-6">
              <div class="border rounded p-3 h-100">
                <h6 class="mb-2 text-capitalize">{module}</h6>
                <p class="text-muted mb-2">{searchResults[module]?.length ?? 0} kết quả</p>
                {#if (searchResults[module]?.length ?? 0) > 0}
                  <ul class="mb-0 ps-3">
                    {#each searchResults[module].slice(0, 3) as result}
                      <li>{String(result.id ?? '-')} - {JSON.stringify(result).slice(0, 80)}...</li>
                    {/each}
                  </ul>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </SectionCard>
    {/if}

    {#if route.kind === 'dashboard'}
      <section>
        <h3 class="mb-3">{tt('dashboardTitle')}</h3>
        <div class="row">
          {#if dashboard && Object.keys(dashboard).length > 0}
            {#each Object.entries(dashboard) as [key, value]}
              <div class="col-md-3 col-sm-6 mb-3">
                <div class="card"><div class="card-body"><p class="text-muted mb-1 text-capitalize">{key.replaceAll('_', ' ')}</p><h4 class="mb-0">{String(value)}</h4></div></div>
              </div>
            {/each}
          {:else}
            <ContentState title="No dashboard data" description="No metrics returned yet." />
          {/if}
        </div>
      </section>
    {:else if route.kind === 'settings'}
      <SettingsForm settings={settings} onSave={saveSettings} t={tt} onLanguageChange={(next: Locale) => setLocale(next)} />
    {:else if route.kind === 'search'}
      <SectionCard title="Search" subtitle="Dedicated search workflow parity at /admin/search">
        <p class="text-muted mb-2">Dedicated search page parity route: <code>/admin/search</code></p>
        <div class="d-flex gap-2">
          <input class="form-control" bind:value={globalSearchTerm} placeholder="Search across modules..." />
          <button class="btn btn-primary" onclick={runGlobalSearch}>Run Search</button>
        </div>
      </SectionCard>
    {:else if route.kind === 'module-list' || route.kind === 'module-new' || route.kind === 'module-detail'}
      <SectionCard title={tt('managementTitle', { module: activeModule })} subtitle="Shared form + table density parity" actions={true}>
          <div slot="actions" class="d-flex gap-2">
            <button class="btn btn-sm btn-outline-primary" onclick={() => go(`/admin/${activeModule}`)}>List</button>
            <button class="btn btn-sm btn-outline-primary" onclick={() => go(`/admin/${activeModule.slice(0, -1)}/new`)}>New</button>
            {#if selectedId}<button class="btn btn-sm btn-outline-primary" onclick={() => go(`/admin/${activeModule.slice(0, -1)}/${selectedId}`)}>Detail</button>{/if}
          </div>
                    {#if route.kind === 'module-new'}
            <p class="text-muted">Create route active: <code>/admin/{activeModule.slice(0, -1)}/new</code></p>
          {/if}
          {#if route.kind === 'module-detail'}
            <p class="text-muted">Detail route active for ID <strong>{route.id}</strong>.</p>
          {/if}
          {#if activeModule === 'users'}
            {@const seg = userSegmentation()}
            <div class="row g-3 mb-3">
              <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Total Users</p><h4 class="mb-0">{seg.users.length}</h4></div></div></div>
              <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Admin Users</p><h4 class="mb-0">{seg.adminUsers.length}</h4></div></div></div>
              <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Staff/Members</p><h4 class="mb-0">{seg.staffUsers.length}</h4></div></div></div>
            </div>
            <div class="row g-3 mb-3">
              <div class="col-md-6"><label class="form-label">Username *</label><input class="form-control" bind:value={userForm.username} /></div>
              <div class="col-md-6"><label class="form-label">Email</label><input type="email" class="form-control" bind:value={userForm.email} /></div>
              <div class="col-md-6"><label class="form-label">Password</label><input type="text" class="form-control" bind:value={userForm.password} placeholder={route.kind==='module-detail'?'******':''} /></div>
              <div class="col-md-6"><label class="form-label">Confirm Password</label><input type="text" class="form-control" bind:value={userForm.confirm_password} /></div>
              <div class="col-md-6"><label class="form-label">Role</label><select class="form-select" bind:value={userForm.role}><option value="guest">Guest</option><option value="admin">Admin</option></select></div>
              <div class="col-md-6"><label class="form-label">Status</label><select class="form-select" bind:value={userForm.status}><option value="Active">Active</option><option value="Inactive">Inactive</option></select></div>
              <div class="col-md-12 d-flex gap-2"><button class="btn btn-primary" onclick={async () => { const e = validateUserForm(route.kind==='module-detail'); if (e) { error = e; return; } syncUserDraftFromForm(); await saveRecord('users'); }}>{route.kind==='module-detail' ? tt('update') : tt('create')}</button><button class="btn btn-outline-secondary" onclick={clearEditor}>{tt('reset')}</button></div>
            </div>
          {:else}
          <div class="row g-3 mb-3">
            <div class="col-md-4"><label class="form-label" for="filter-input">{tt('filterLabel')}</label><input id="filter-input" class="form-control" bind:value={filter} placeholder={tt('filterPlaceholder')} /></div>
            <div class="col-md-3"><label class="form-label" for="sort-input">{tt('sortById')}</label><select id="sort-input" class="form-select" bind:value={sort}><option value="asc">asc</option><option value="desc">desc</option></select></div>
            <div class="col-md-5 d-flex align-items-end gap-2"><button class="btn btn-primary" onclick={() => saveRecord(activeModule)}>{selectedId ? tt('update') : tt('create')}</button><button class="btn btn-outline-secondary" onclick={clearEditor}>{tt('reset')}</button></div>
            <div class="col-12"><label class="form-label" for="draft-json">JSON payload</label><textarea id="draft-json" class="form-control" rows="8" bind:value={draftText} onchange={() => {
              try {
                draft = JSON.parse(draftText);
              } catch (parseError) {
                error = (parseError as Error).message;
              }
            }}></textarea></div>
          </div>
          {/if}
      </SectionCard>
      {#if activeModule === 'cars' || activeModule === 'customers'}
        <UploadForm onUpload={uploadImage} t={tt} />
      {/if}
      {@const view = visibleItems(activeModule)}
      <ModuleTable title={activeModule} items={view.items} total={view.total} page={page} pageSize={pageSize} onDetail={(id) => selectRecord(activeModule, id)} onDelete={(id) => removeRecord(activeModule, id)} t={tt} />
      <div class="d-flex justify-content-end align-items-center gap-2 mt-3">
        <button class="btn btn-outline-secondary btn-sm" disabled={page<=1} onclick={() => page = page - 1}>{tt('prev')}</button>
        <span>{tt('page')} {page}</span>
        <button class="btn btn-outline-secondary btn-sm" disabled={page*pageSize>=view.total} onclick={() => page = page + 1}>{tt('next')}</button>
      </div>
    {/if}
    </AdminLayout>
{:else if route.kind === 'not-found'}
    <div class="card">
      <div class="card-body">
        <h3 class="mb-2">404</h3>
        <p class="mb-3">Route not found: <code>{route.path}</code></p>
        <button class="btn btn-primary" onclick={() => go('/admin/dashboard')}>Go to dashboard</button>
      </div>
    </div>
  {:else}
    <p><a href="#/auth/login">{tt('login')}</a> | <a href="#/admin/dashboard">{tt('admin')}</a></p>
  {/if}

</main>
