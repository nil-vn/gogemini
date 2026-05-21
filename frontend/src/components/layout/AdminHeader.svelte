<script lang="ts">
  export let onSearch: () => void;
  export let onThemeChange: (theme: 'dark' | 'light') => void;
  export let searchTerm = '';
  export let onSearchTermInput: (value: string) => void;
</script>

<div class="header-wrapper">
  <div class="me-auto pc-mob-drp">
    <ul class="list-unstyled">
      <li class="pc-h-item pc-sidebar-collapse"><button type="button" class="pc-head-link ms-0" aria-label="Toggle sidebar"><i class="ph ph-list"></i></button></li>
      <li class="dropdown pc-h-item">
        <button type="button" class="pc-head-link dropdown-toggle arrow-none m-0 trig-drp-search" aria-label="Open search"><i class="ph ph-magnifying-glass"></i></button>
      </li>
    </ul>
  </div>
  <div class="ms-auto d-flex align-items-center gap-2 px-2">
    <input class="form-control" placeholder="Search here..." value={searchTerm} oninput={(e) => onSearchTermInput((e.currentTarget as HTMLInputElement).value)} />
    <button class="btn btn-primary" onclick={onSearch}>Search</button>
    <button class="btn btn-outline-secondary" onclick={() => onThemeChange('dark')}>Dark</button>
    <button class="btn btn-outline-secondary" onclick={() => onThemeChange('light')}>Light</button>
  </div>
</div>
