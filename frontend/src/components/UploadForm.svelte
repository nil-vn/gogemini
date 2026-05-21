<script lang="ts">
  import { validateUploadFile } from '../lib/validation';
  export let onUpload: (payload: { module: 'cars' | 'customers'; file: File }) => Promise<void>;
  export let t: (key: string) => string = (key) => key;
  let module: 'cars' | 'customers' = 'cars';
  let fileInput: HTMLInputElement;
  let preview = '';
  let localError = '';
  let fileName = '';

  function changeFile() {
    const file = fileInput?.files?.[0];
    if (!file) return;
    fileName = file.name;
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

<div class="card mt-3">
  <div class="card-header"><h5 class="mb-0">{t('uploadImage')}</h5></div>
  <div class="card-body">
    <div class="row g-3 align-items-end">
      <div class="col-md-3">
        <label class="form-label" for="upload-module">Module</label>
        <select id="upload-module" class="form-select" bind:value={module}><option value="cars">cars</option><option value="customers">customers</option></select>
      </div>
      <div class="col-md-6">
        <label class="form-label" for="upload-file">Image</label>
        <input class="form-control" id="upload-file" bind:this={fileInput} type="file" accept="image/*" onchange={changeFile} required />
      </div>
      <div class="col-md-3">
        <button class="btn btn-primary w-100" onclick={submit}>{t('upload')}</button>
      </div>
    </div>
    {#if localError}<p class="text-danger mt-2 mb-0">{localError}</p>{/if}
    {#if fileName}<p class="text-muted mt-2 mb-1">{fileName}</p>{/if}
    {#if preview}<img alt="preview" src={preview} class="img-fluid rounded border mt-2" style="max-height:240px" />{/if}
  </div>
</div>
