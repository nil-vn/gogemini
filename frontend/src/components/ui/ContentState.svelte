<script lang="ts">
  export let title = '';
  export let description = '';
  export let actionLabel = '';
  export let onAction: () => void = () => {};
</script>

<div class="card">
  <div class="card-body text-center py-5">
    <h5 class="mb-2">{title}</h5>
    <p class="text-muted mb-3">{description}</p>
    {#if actionLabel}
      <button class="btn btn-outline-primary" onclick={onAction}>{actionLabel}</button>
    {/if}
  </div>
</div>
