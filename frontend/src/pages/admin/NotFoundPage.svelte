<script lang="ts">
  export let path = '';
  export let go: (path:string)=>void;
</script>
<div class="card">
  <div class="card-body">
    <h3 class="mb-2">404</h3>
    <p class="text-muted">Route not found: <code>{path}</code></p>
    <button class="btn btn-primary" onclick={() => go('/admin/dashboard')}>Back to Dashboard</button>
  </div>
</div>
