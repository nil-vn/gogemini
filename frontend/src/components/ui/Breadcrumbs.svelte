<script lang="ts">
  export type Crumb = { label: string; path?: string };
  export let items: Crumb[] = [];
  export let onNavigate: (path: string) => void;
</script>

<nav aria-label="breadcrumb" class="mb-3">
  <ol class="breadcrumb mb-0">
    {#each items as item, idx}
      <li class="breadcrumb-item {idx === items.length - 1 ? 'active' : ''}" aria-current={idx === items.length - 1 ? 'page' : undefined}>
        {#if item.path && idx !== items.length - 1}
          <button type="button" class="btn btn-link btn-sm p-0 text-decoration-none" onclick={() => onNavigate(item.path!)}>{item.label}</button>
        {:else}
          {item.label}
        {/if}
      </li>
    {/each}
  </ol>
</nav>
