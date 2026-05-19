<script lang="ts">
  export let title = '';
  export let items: any[] = [];
</script>

<section>
  <h3>{title}</h3>
  <pre>{JSON.stringify(items, null, 2)}</pre>
</section>
