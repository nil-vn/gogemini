<script lang="ts">
  import { validateUploadFile } from '../lib/validation';
  export let onUpload: (payload: { module: 'cars' | 'customers'; file: File }) => Promise<void>;
  export let t: (key: string) => string = (key) => key;
  let module: 'cars' | 'customers' = 'cars';
  let fileInput: HTMLInputElement;
  let preview = '';
  let localError = '';

  function changeFile() {
    const file = fileInput?.files?.[0];
    if (!file) return;
    preview = URL.createObjectURL(file);
  }

  async function submit() {
    const file = fileInput?.files?.[0];
    const errs = validateUploadFile(file);
    if (errs.length) {
      localError = t(errs[0]);
      return;
    }
    localError = '';
    await onUpload({ module, file: file! });
  }
</script>

<section>
  <h3>{t('uploadImage')}</h3>
  <label>Module <select bind:value={module}><option value="cars">cars</option><option value="customers">customers</option></select></label>
  <label for="upload-file">Image</label>
  <input id="upload-file" bind:this={fileInput} type="file" accept="image/*" onchange={changeFile} required />
  <button onclick={submit}>{t('upload')}</button>
  {#if localError}<p style="color:red">{localError}</p>{/if}
  {#if preview}<img alt="preview" src={preview} style="max-width: 240px;display:block;" />{/if}
</section>
