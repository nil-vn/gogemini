<script lang="ts">
  export let onSubmit: (payload: { login: string; password: string }) => Promise<void>;
  export let t: (key: string) => string = (key) => key;
  let login = '';
  let password = '';
</script>

<form on:submit|preventDefault={() => onSubmit({ login, password })} aria-label="login-form">
  <label for="login">Login</label>
  <input id="login" bind:value={login} placeholder="username/email" required autocomplete="username" />
  <label for="password">Password</label>
  <input id="password" bind:value={password} type="password" placeholder="password" required autocomplete="current-password" />
  <button type="submit">{t('login')}</button>
</form>
