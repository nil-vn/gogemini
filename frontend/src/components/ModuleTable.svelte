<script lang="ts">
  import type { ModuleRecord } from '../lib/types';
  export let title = '';
  export let items: ModuleRecord[] = [];
  export let total = 0;
  export let page = 1;
  export let pageSize = 10;
  export let onDetail: (id: string) => void;
  export let onDelete: (id: string) => void;
  export let t: (key: string, vars?: Record<string, string | number>) => string;

  function getId(item: ModuleRecord) {
    return String(item.id ?? '');
  }
</script>

<section>
  <h4>{t('tableList', { title })}</h4>
  <p>{t('tableTotal', { total, page, pageSize })}</p>
  <table>
    <thead><tr><th>{t('id')}</th><th>{t('data')}</th><th>{t('actions')}</th></tr></thead>
    <tbody>
      {#if items.length === 0}
        <tr><td colspan="3">{t('noRecords')}</td></tr>
      {:else}
        {#each items as item}
          <tr>
            <td>{getId(item)}</td>
            <td><pre>{JSON.stringify(item, null, 2)}</pre></td>
            <td>
              <button on:click={() => onDetail(getId(item))}>{t('detail')}</button>
              <button on:click={() => onDelete(getId(item))}>{t('delete')}</button>
            </td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</section>
