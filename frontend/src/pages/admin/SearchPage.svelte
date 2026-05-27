<script lang="ts">
  export let globalSearchTerm = '';
  export let setGlobalSearchTerm: (v: string) => void;
  export let runGlobalSearch: any;
</script>

<div class="d-flex flex-column gap-2">
  <div class="d-flex gap-2">
    <input
      class="form-control"
      value={globalSearchTerm}
      oninput={(e) => setGlobalSearchTerm((e.currentTarget as HTMLInputElement).value)}
      placeholder="Tìm theo mã giao dịch, biển số xe, tên khách hàng..."
    />
    <button class="btn btn-primary" onclick={runGlobalSearch}>Tìm kiếm</button>
  </div>
  <p class="text-muted mb-0">Nhập từ khóa để tra cứu nhanh theo từng nhóm: users, cars, customers, transactions.</p>
</div>
