<script lang="ts">
  export let activePath = '/admin/dashboard';
  export let onNavigate: (path: string) => void;

  const isActive = (path: string) => activePath === path;
</script>

<div class="navbar-wrapper">
  <div class="m-header"><button type="button" class="b-brand text-primary border-0 bg-transparent" onclick={() => onNavigate('/admin/dashboard')}>Gemini CRM</button></div>
  <div class="navbar-content">
    <ul class="pc-navbar">
      <li class="pc-item pc-caption"><span>Navigation</span></li>
      <li class="pc-item"><button type="button" class="pc-link border-0 bg-transparent text-start w-100 {isActive('/admin/dashboard') ? 'active' : ''}" onclick={() => onNavigate('/admin/dashboard')}><span class="pc-mtext">Dashboard</span></button></li>
      {#each ['users','cars','customers','transactions'] as m}
        <li class="pc-item"><button type="button" class="pc-link border-0 bg-transparent text-start w-100 {isActive(`/admin/${m}`) ? 'active' : ''}" onclick={() => onNavigate(`/admin/${m}`)}><span class="pc-mtext text-capitalize">{m}</span></button></li>
      {/each}
      <li class="pc-item"><button type="button" class="pc-link border-0 bg-transparent text-start w-100 {isActive('/admin/system') ? 'active' : ''}" onclick={() => onNavigate('/admin/system')}><span class="pc-mtext">System</span></button></li>
    </ul>
  </div>
</div>
