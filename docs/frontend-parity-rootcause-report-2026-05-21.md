# Frontend Parity Root-Cause Investigation Report (Flask ➜ Svelte 5)

Date: 2026-05-21 (UTC)
Owner: Codex investigation

## Executive summary
UI của Svelte 5 chưa giống Flask cũ **không phải do “thiếu polish nhỏ”** mà do lệch chiến lược triển khai: code hiện tại ưu tiên một **generic CRUD shell + JSON editor** thay vì tái tạo IA/page model của Flask (list/new/detail theo module + workflow chuyên biệt).

Kết quả là dù đã có “WBS breakdown” và một phần implement, phần lớn artifact đang dừng ở mức readiness/planning; implementation production cho parity chưa đi đến các page/component domain-specific cần thiết.

## Evidence snapshot

### 1) Svelte app đang dùng kiến trúc generic module screen
- Route model gom về `module-list`, `module-new`, `module-detail` dùng chung 4 module users/cars/customers/transactions thay vì màn hình chuyên biệt theo từng module logic. (`frontend/src/App.svelte`)
- Table component render dữ liệu kiểu generic (`id`, `data`, `detail/delete`) và cắt summary fields động, không có cột domain-specific như Flask. (`frontend/src/components/ModuleTable.svelte`)

### 2) Legacy Flask có page semantics đậm domain + workflow
- Cars page có tab segmentation nhiều trạng thái + KPI cards + chart placeholders và bảng cột nghiệp vụ chi tiết. (`templates/admin/cars.html`)
- Transactions page có tab All/Deposited/Paid + financial KPI cards + table lồng quan hệ customer/car + status-specific rendering. (`templates/admin/transactions.html`)
- Transaction create có dynamic line items `item_name[]/item_price[]`, sticky action bar, datepicker/choices integrations. (`templates/admin/transaction_new.html`)
- Search page legacy có grouping result theo entity và deep-link tới detail pages. (`templates/admin/search.html`)

### 3) Project docs đã tự ghi nhận migration gap
- Gap assessment xác nhận “chưa migrate hoàn chỉnh parity” và liệt kê thiếu create/detail/search/purchase-flow/404/UI blocks parity. (`docs/wbs6-frontend-migration-gap-assessment-2026-05-21.md`)

## Root causes (5 Why condensed)

### RC1 — Sai “unit of migration”: migrate theo API module thay vì operator workflow/page contract
- Why UI khác: vì FE mới map 4 module vào một renderer generic.
- Why làm vậy: ưu tiên API closure + tốc độ scaffold.
- Why ảnh hưởng parity: legacy UX dựa trên tác vụ vận hành (segmentation, KPI, shortcuts), không chỉ CRUD raw.

### RC2 — WBS execution skew: hoàn thành planning artifacts nhưng implementation chưa đạt feature depth
- WBS có breakdown, nhưng code hiện hữu chưa phản ánh đầy đủ page set legacy (new/detail/search parity chuyên biệt).
- “Implemented 6 WBS” có khả năng là completion theo checklist/task admin, không phải parity acceptance theo màn hình thực thi.

### RC3 — Acceptance criteria thiếu “pixel/interaction parity gate” bắt buộc trước sign-off
- Chưa thấy bằng chứng gate kiểu side-by-side visual diff/page contract cho từng template legacy trước khi nhận done.
- Do đó “done” bị hiểu là API chạy + route tồn tại thay vì UX tương đương.

### RC4 — Reuse static theme assets không đồng nghĩa với behavior parity
- Static CSS/JS đã được copy sang frontend public, nhưng interactive behavior của template cũ (datatable config, dynamic form rows, contextual actions) chưa được tái hiện đầy đủ trong Svelte components.

### RC5 — Scope pressure + sequencing chưa khóa critical-path parity
- Triển khai có vẻ đi theo horizontal platform tasks (layout, generic table, base forms), trong khi parity cần vertical slice theo từng module workflow end-to-end.

## Impact
- Ops users mất mental model cũ → tăng thời gian thao tác và lỗi nhập liệu.
- Flow trọng yếu (transaction itemization, purchase shortcut từ car/customer) không tương đương.
- Rủi ro cutover cao vì test pass kỹ thuật không phản ánh readiness vận hành.

## Corrective action plan (pragmatic)
1. Freeze claim “parity done”; re-baseline scope theo legacy page contract (list/new/detail/search/404 + workflow hooks).
2. Đổi execution strategy sang **vertical module parity sprints** (Users → Cars → Customers → Transactions → Search).
3. Bắt buộc parity gate cho mỗi page:
   - route parity,
   - control parity,
   - interaction parity,
   - navigation/deep-link parity,
   - operator UAT sign-off.
4. Loại bỏ JSON editor khỏi primary path (giữ debug-only behind feature flag).
5. Thêm automated evidence:
   - visual baseline per critical page,
   - Playwright flows cho add/edit/delete + domain actions,
   - checklist pass/fail có owner và deadline.

## Immediate next 7-day priorities
- Day 1-2: chốt route/page matrix & acceptance checklist.
- Day 3-4: hoàn thiện Transactions parity (new/detail + item rows + summary cards) vì business critical nhất.
- Day 5: Cars parity (segmentation tabs + KPI + detail actions).
- Day 6: Search + cross-link flows.
- Day 7: UAT replay với team vận hành + defect burn-down.

## Conclusion
Root cause chính là **định nghĩa “migrate xong” sai tầng**: hoàn thành technical scaffold/WBS artifacts nhưng chưa hoàn thành UI workflow parity theo contract của Flask legacy. Cần reset tiêu chí nghiệm thu sang parity-by-page + parity-by-workflow, không chỉ parity-by-endpoint.
