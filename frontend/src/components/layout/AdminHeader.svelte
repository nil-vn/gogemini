<script lang="ts">
  import { onMount } from 'svelte';

  export let onSearch: () => void;
  export let onThemeChange: (theme: 'dark' | 'light') => void;
  export let searchTerm = '';
  export let onSearchTermInput: (value: string) => void;
  export let onToggleSidebar: () => void;
  export let t: (key: string) => string = (key) => key;

  let searchInput: HTMLInputElement;

  function openSearch() {
    searchInput?.focus();
  }

  function setTheme(theme: 'dark' | 'light') {
    onThemeChange(theme);
    localStorage.setItem('pc-theme', theme);
  }

  onMount(() => {
    const savedTheme = localStorage.getItem('pc-theme');
    if (savedTheme === 'dark' || savedTheme === 'light') {
      onThemeChange(savedTheme);
    }
  });
</script>

<div class="header-wrapper">
  <div class="me-auto pc-mob-drp">
    <ul class="list-unstyled">
      <li class="pc-h-item pc-sidebar-collapse"><button type="button" class="pc-head-link ms-0" aria-label={t('toggleSidebar')} onclick={onToggleSidebar}><i class="ph ph-list"></i></button></li>
      <li class="dropdown pc-h-item">
        <button type="button" class="pc-head-link dropdown-toggle arrow-none m-0 trig-drp-search" aria-label={t('openSearch')} onclick={openSearch}><i class="ph ph-magnifying-glass"></i></button>
      </li>
    </ul>
  </div>
  <div class="ms-auto d-flex align-items-center gap-2 px-2">
    <input class="form-control" bind:this={searchInput} placeholder={t('searchInputPlaceholder')} value={searchTerm} oninput={(e) => onSearchTermInput((e.currentTarget as HTMLInputElement).value)} />
    <button class="btn btn-primary" onclick={onSearch}>Search</button>
    <button class="btn btn-outline-secondary" onclick={() => setTheme('dark')}>Dark</button>
    <button class="btn btn-outline-secondary" onclick={() => setTheme('light')}>Light</button>
  </div>
</div>
