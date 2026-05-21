<script lang="ts">
  export let routeMode: 'list'|'new'|'detail' = 'list';
  export let transactionForm: any;
  export let transactionItems: any[];
  export let setError: (s:string)=>void;
  export let syncTransactionDraftFromForm: any;
  export let saveRecord: any;
  export let clearEditor: any;
  export let addTransactionItem: any;
  export let removeTransactionItem: any;
  export let transactionSummary: any;
  export let transactionStatusGroups: any;
  export let filter: string;
  export let setFilter: (s:string)=>void;
  export let tt: any;
</script>
{@const txSum = transactionSummary()}
{@const txGroup = transactionStatusGroups()}
<div class="row g-3 mb-3">
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Total Revenue</p><h4 class="mb-0">{txSum.totalRevenue.toLocaleString()}</h4></div></div></div>
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Paid Revenue</p><h4 class="mb-0">{txSum.paidRevenue.toLocaleString()}</h4></div></div></div>
  <div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Deposited Amount</p><h4 class="mb-0">{txSum.depositedAmount.toLocaleString()}</h4></div></div></div>
</div>
<div class="d-flex flex-wrap gap-2 mb-3">
  <button class="btn btn-sm {filter==='tx_all'?'btn-primary':'btn-outline-primary'}" onclick={() => setFilter('tx_all')}>All ({txGroup.all.length})</button>
  <button class="btn btn-sm {filter==='tx_deposited'?'btn-primary':'btn-outline-primary'}" onclick={() => setFilter('tx_deposited')}>Deposited ({txGroup.deposited.length})</button>
  <button class="btn btn-sm {filter==='tx_paid'?'btn-primary':'btn-outline-primary'}" onclick={() => setFilter('tx_paid')}>Paid ({txGroup.paid.length})</button>
</div>
<div class="row g-3 mb-3">
  <div class="col-md-6"><label class="form-label">Customer ID *</label><input class="form-control" bind:value={transactionForm.customer_id} /></div>
  <div class="col-md-6"><label class="form-label">Car ID *</label><input class="form-control" bind:value={transactionForm.car_id} /></div>
  <div class="col-md-12"><h6>Other Transactions (Accessories/Services)</h6>
    <div class="table-responsive"><table class="table table-bordered"><tbody>{#each transactionItems as item, idx}<tr><td><input class="form-control" bind:value={item.name} /></td><td><input class="form-control" bind:value={item.price} /></td><td><button class="btn btn-sm btn-danger" onclick={() => removeTransactionItem(idx)}>Remove</button></td></tr>{/each}</tbody></table></div>
    <button class="btn btn-sm btn-primary" onclick={addTransactionItem}>+ Add Item</button></div>
  <div class="col-md-12 d-flex flex-wrap gap-2"><button class="btn btn-primary" onclick={async () => { if (!transactionForm.customer_id || !transactionForm.car_id) { setError('Customer and Car are required'); return; } syncTransactionDraftFromForm(); await saveRecord('transactions'); }}>{routeMode==='detail' ? tt('update') : tt('create')}</button><button class="btn btn-outline-secondary" onclick={clearEditor}>{tt('reset')}</button></div>
</div>
