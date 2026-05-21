<script lang="ts">
  export let routeMode: 'list'|'new'|'detail' = 'list';
  export let userForm: any;
  export let validateUserForm: any;
  export let syncUserDraftFromForm: any;
  export let saveRecord: any;
  export let clearEditor: any;
  export let userSegmentation: any;
  export let tt: any;
</script>
{@const seg = userSegmentation()}
<div class="row g-3 mb-3">
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Total Users</p><h4 class="mb-0">{seg.users.length}</h4></div></div></div>
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Admin Users</p><h4 class="mb-0">{seg.adminUsers.length}</h4></div></div></div>
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Staff/Members</p><h4 class="mb-0">{seg.staffUsers.length}</h4></div></div></div>
</div>
<div class="row g-3 mb-3">
  <div class="col-md-6"><label class="form-label">Username *</label><input class="form-control" bind:value={userForm.username} /></div>
  <div class="col-md-6"><label class="form-label">Email</label><input type="email" class="form-control" bind:value={userForm.email} /></div>
  <div class="col-md-6"><label class="form-label">Password</label><input type="text" class="form-control" bind:value={userForm.password} placeholder={routeMode==='detail'?'******':''} /></div>
  <div class="col-md-6"><label class="form-label">Confirm Password</label><input type="text" class="form-control" bind:value={userForm.confirm_password} /></div>
  <div class="col-md-6"><label class="form-label">Role</label><select class="form-select" bind:value={userForm.role}><option value="guest">Guest</option><option value="admin">Admin</option></select></div>
  <div class="col-md-6"><label class="form-label">Status</label><select class="form-select" bind:value={userForm.status}><option value="Active">Active</option><option value="Inactive">Inactive</option></select></div>
  <div class="col-md-12 d-flex gap-2"><button class="btn btn-primary" onclick={async () => { const e = validateUserForm(routeMode==='detail'); if (e) return; syncUserDraftFromForm(); await saveRecord('users'); }}>{routeMode==='detail' ? tt('update') : tt('create')}</button><button class="btn btn-outline-secondary" onclick={clearEditor}>{tt('reset')}</button></div>
</div>
