<script lang="ts">
  import type { Settings } from '../lib/types';
  import type { Locale } from '../lib/i18n';
  export let settings: Settings;
  export let onSave: (s: Settings) => Promise<void>;
  export let t: (key: string, vars?: Record<string, string | number>) => string;
  export let onLanguageChange: (locale: Locale) => void;
</script>

<div class="card">
  <div class="card-header"><h5 class="mb-0">{t('settingsTitle')}</h5></div>
  <div class="card-body">
    <div class="row g-3">
      <div class="col-md-4"><label class="form-label" for="settings-currency">{t('currency')}</label><input id="settings-currency" class="form-control" bind:value={settings.currency} /></div>
      <div class="col-md-4"><label class="form-label" for="settings-theme">{t('theme')}</label>
        <select id="settings-theme" class="form-select" bind:value={settings.theme}><option value="dark">{t('dark')}</option><option value="light">{t('light')}</option></select>
      </div>
      <div class="col-md-4"><label class="form-label" for="settings-language">{t('language')}</label>
        <select id="settings-language" class="form-select" bind:value={settings.language} onchange={() => onLanguageChange(settings.language)}><option value="vi">vi</option><option value="en">en</option></select>
      </div>
      <div class="col-12"><button class="btn btn-primary" onclick={() => onSave(settings)}>{t('save')}</button></div>
    </div>
  </div>
</div>
