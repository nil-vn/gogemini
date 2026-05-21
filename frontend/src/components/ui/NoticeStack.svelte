<script lang="ts">
  export let isLoading = false;
  export let message = '';
  export let error = '';
  export let onRetry: () => void = () => {};
  export let t: (key: string) => string = (key) => key;
</script>

<section aria-live="polite" aria-label={t('a11yStatusLabel')} class="mb-3">
  {#if isLoading}
    <div class="alert alert-info mb-2" role="status">{t('loading')}</div>
  {/if}
  {#if message}
    <div class="alert alert-success mb-2" role="status">{message}</div>
  {/if}
  {#if error}
    <div class="alert alert-danger d-flex justify-content-between align-items-center gap-2 mb-2" role="alert">
      <span>{t('errorTitle')}: {error}</span>
      <button class="btn btn-sm btn-outline-danger" onclick={onRetry}>{t('retry')}</button>
    </div>
  {/if}
</section>
