<script lang="ts">
  export let title = '';
  export let subtitle = '';
  export let actions = false;
</script>

<div class="card mb-3">
  <div class="card-header d-flex justify-content-between align-items-center">
    <div>
      <h5 class="mb-0">{title}</h5>
      {#if subtitle}<small class="text-muted">{subtitle}</small>{/if}
    </div>
    {#if actions}
      <div><slot name="actions" /></div>
    {/if}
  </div>
  <div class="card-body">
    <slot />
  </div>
</div>
