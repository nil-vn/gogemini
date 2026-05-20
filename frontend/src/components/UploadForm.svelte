<script lang="ts">
  export let onUpload: (payload: { module: 'cars' | 'customers'; file: File }) => Promise<void>;
  export let t: (key: string) => string = (key) => key;
  let module: 'cars' | 'customers' = 'cars';
  let fileInput: HTMLInputElement;
  let preview = '';

  function changeFile() {
    const file = fileInput?.files?.[0];
    if (!file) return;
    preview = URL.createObjectURL(file);
  }

  async function submit() {
    const file = fileInput?.files?.[0];
    if (!file) return;
    await onUpload({ module, file });
  }
</script>

<section>
  <h3>{t('uploadImage')}</h3>
  <select bind:value={module}><option value="cars">cars</option><option value="customers">customers</option></select>
  <input bind:this={fileInput} type="file" accept="image/*" on:change={changeFile} required />
  <button on:click={submit}>{t('upload')}</button>
  {#if preview}<img alt="preview" src={preview} style="max-width: 240px;display:block;" />{/if}
</section>
