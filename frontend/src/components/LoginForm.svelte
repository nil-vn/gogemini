<script lang="ts">
  export let onSubmit: (payload: { login: string; password: string }) => Promise<void>;
  let login = '';
  let password = '';
</script>

<form on:submit|preventDefault={() => onSubmit({ login, password })}>
  <input bind:value={login} placeholder="username/email" required />
  <input bind:value={password} type="password" placeholder="password" required />
  <button type="submit">Login</button>
</form>
