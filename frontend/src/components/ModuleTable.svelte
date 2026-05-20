<script lang="ts">
  import type { ModuleRecord } from '../lib/types';
  export let title = '';
  export let items: ModuleRecord[] = [];
  export let total = 0;
  export let page = 1;
  export let pageSize = 10;
  export let onDetail: (id: string) => void;
  export let onDelete: (id: string) => void;

  function getId(item: ModuleRecord) {
    return String(item.id ?? '');
  }
</script>

<section>
  <h4>{title} list</h4>
  <p>total: {total} | page: {page} | page size: {pageSize}</p>
  <table>
    <thead><tr><th>id</th><th>data</th><th>actions</th></tr></thead>
    <tbody>
      {#if items.length === 0}
        <tr><td colspan="3">No records</td></tr>
      {:else}
        {#each items as item}
          <tr>
            <td>{getId(item)}</td>
            <td><pre>{JSON.stringify(item, null, 2)}</pre></td>
            <td>
              <button on:click={() => onDetail(getId(item))}>Detail</button>
              <button on:click={() => onDelete(getId(item))}>Delete</button>
            </td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</section>
