<script lang="ts">
  import type { Settings } from '../lib/types';
  export let settings: Settings;
  export let onSave: (s: Settings) => Promise<void>;
</script>

<section>
  <h3>Settings</h3>
  <label>Currency <input bind:value={settings.currency} /></label>
  <label>Theme
    <select bind:value={settings.theme}>
      <option value="dark">dark</option>
      <option value="light">light</option>
    </select>
  </label>
  <label>Language
    <select bind:value={settings.language}>
      <option value="vi">vi</option>
      <option value="en">en</option>
    </select>
  </label>
  <button on:click={() => onSave(settings)}>Save</button>
</section>
