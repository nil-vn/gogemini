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

  function summaryFields(item: ModuleRecord): string {
    const hiddenKeys = new Set(['id', 'password', 'created_at', 'updated_at']);
    return Object.entries(item)
      .filter(([key, value]) => !hiddenKeys.has(key) && value !== null && value !== undefined && String(value).trim() !== '')
      .slice(0, 3)
      .map(([key, value]) => `${key}: ${value}`)
      .join(' • ');
  }
</script>

<div class="card">
  <div class="card-header d-flex align-items-center justify-content-between">
    <h5 class="mb-0">{t('tableList', { title })}</h5>
    <span class="badge bg-light-primary">{t('tableTotal', { total, page, pageSize })}</span>
  </div>
  <div class="card-body p-0">
    <div class="table-responsive">
      <table class="table table-hover mb-0">
        <thead>
          <tr>
            <th>{t('id')}</th>
            <th>{t('data')}</th>
            <th class="text-end">{t('actions')}</th>
          </tr>
        </thead>
        <tbody>
          {#if items.length === 0}
            <tr><td colspan="3" class="text-center py-4 text-muted">{t('noRecords')}</td></tr>
          {:else}
            {#each items as item}
              <tr>
                <td class="fw-semibold">{getId(item)}</td>
                <td>{summaryFields(item) || '-'}</td>
                <td class="text-end">
                  <button class="btn btn-sm btn-outline-primary me-2" onclick={() => onDetail(getId(item))}>{t('detail')}</button>
                  <button class="btn btn-sm btn-outline-danger" onclick={() => onDelete(getId(item))}>{t('delete')}</button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
