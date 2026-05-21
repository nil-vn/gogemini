<script lang="ts">
  import AdminHeader from './AdminHeader.svelte';
  import AdminSidebar from './AdminSidebar.svelte';
  
  export let activePath = '/admin/dashboard';
  export let globalSearchTerm = '';
  export let onNavigate: (path: string) => void;
  export let onSearchTermInput: (value: string) => void;
  export let onSearch: () => void;
  export let onThemeChange: (theme: 'dark' | 'light') => void;
  export let t: (key: string) => string = (key) => key;

  let isSidebarCollapsed = false;
  function toggleSidebar() {
    isSidebarCollapsed = !isSidebarCollapsed;
  }
</script>

<a href="#main-content" class="skip-link">{t('a11ySkipToContent')}</a>
<div class:pc-sidebar-hide={isSidebarCollapsed}>
  <nav class="pc-sidebar" aria-label="Admin sidebar">
    <AdminSidebar {activePath} {onNavigate} {t} />
  </nav>
</div>
<header class="pc-header" aria-label="Admin header">
  <AdminHeader onSearch={onSearch} onThemeChange={onThemeChange} searchTerm={globalSearchTerm} {onSearchTermInput} onToggleSidebar={toggleSidebar} {t} />
</header>
<div class="pc-container">
  <div class="pc-content" id="main-content" tabindex="-1">
    <slot />
  </div>
</div>
<footer class="pc-footer"><div class="footer-wrapper container-fluid"><div class="row"><div class="col my-1"><p class="m-0">Gemini CRM. Copyright © 2026</p></div></div></div></footer>
