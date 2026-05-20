<script lang="ts">
  import type { Settings } from '../lib/types';
  import type { Locale } from '../lib/i18n';
  export let settings: Settings;
  export let onSave: (s: Settings) => Promise<void>;
  export let t: (key: string, vars?: Record<string, string | number>) => string;
  export let onLanguageChange: (locale: Locale) => void;
</script>

<section>
  <h3>{t('settingsTitle')}</h3>
  <label>{t('currency')} <input bind:value={settings.currency} /></label>
  <label>{t('theme')}
    <select bind:value={settings.theme}>
      <option value="dark">{t('dark')}</option>
      <option value="light">{t('light')}</option>
    </select>
  </label>
  <label>{t('language')}
    <select bind:value={settings.language} on:change={() => onLanguageChange(settings.language)}>
      <option value="vi">vi</option>
      <option value="en">en</option>
    </select>
  </label>
  <button on:click={() => onSave(settings)}>{t('save')}</button>
</section>
