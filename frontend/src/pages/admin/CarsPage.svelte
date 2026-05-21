<script lang="ts">
  export let routeMode:'list'|'new'|'detail'='list';
  export let carForm:any; export let carSegmentation:any; export let carSegment:string; export let setCarSegment:(v:string)=>void;
  export let pendingCarImages:File[]; export let setPendingCarImages:(v:File[])=>void;
  export let saveCarRecord:any; export let clearEditor:any; export let go:any; export let selectedId=''; export let currentRouteId:any; export let tt:any;
</script>
{@const cseg = carSegmentation()}
<div class="row g-3 mb-3"><div class="col-md-4"><div class="card"><div class="card-body"><p class="text-muted mb-1">Total Cars</p><h4 class="mb-0">{cseg.cars.length}</h4></div></div></div></div>
<div class="d-flex flex-wrap gap-2 mb-3"><button class="btn btn-sm {carSegment==='all'?'btn-primary':'btn-outline-primary'}" onclick={() => setCarSegment('all')}>All</button></div>
<div class="row g-3 mb-3">
<div class="col-md-6"><label class="form-label">Name *</label><input class="form-control" bind:value={carForm.name} /></div>
<div class="col-md-6"><label class="form-label">Model</label><input class="form-control" bind:value={carForm.model} /></div>
<div class="col-md-12"><label class="form-label">Images</label><input multiple type="file" class="form-control" accept="image/*" onchange={(e)=>setPendingCarImages(Array.from((e.currentTarget as HTMLInputElement).files??[]))} /></div>
<div class="col-md-12 d-flex gap-2"><button class="btn btn-primary" onclick={saveCarRecord}>{routeMode==='detail'?tt('update'):tt('create')}</button>{#if routeMode==='detail'}<button class="btn btn-outline-success" onclick={() => go(`/admin/transaction/new?car_id=${selectedId || currentRouteId()}`)}>Add Customer Purchase</button>{/if}<button class="btn btn-outline-secondary" onclick={clearEditor}>{tt('reset')}</button></div>
</div>
