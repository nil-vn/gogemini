# A2 — Flask Endpoint & Behavior Audit (API Contract Baseline)

Mục tiêu tài liệu: khóa baseline behavior hiện tại của Flask làm mốc parity cho migration Go API.

## 1) Scope audit

- In-scope: các endpoint/flow thuộc parity v1 theo A1: auth, dashboard, users, cars, customers, transactions, search, system settings, upload, health-like test API.
- Source code gốc được audit trực tiếp từ `app/admin/controllers/*`, `app/homepage/controllers/*`, `app/__init__.py`.
- Out-of-scope: endpoint/tooling không nằm trong controllers hiện tại.

## 2) Cross-cutting behavior baseline

- Prefix admin web routes: `/admin` (Blueprint `admin_routes`).
- Prefix API test routes: `/api` (Blueprint `admin_api_routes`, `homepage_api_routes`).
- Auth framework: Flask-Login session-based.
  - `login_manager.login_view = "admin_routes.login"`.
  - Một số route đã bật `@login_required`, một số route hiện đang comment `# @login_required` (cần giữ parity khi migrate rồi hardening ở task sau).
- Login behavior:
  - Nếu đã authenticated, vào `/admin/login` sẽ redirect `/admin/dashboard`.
  - Login success redirect về `next` query param nếu có, nếu không về dashboard.
  - Login fail trả lại view login + flash message.
- Error pattern:
  - Chủ yếu flash message + redirect/render HTML.
  - API JSON rõ ràng chỉ xuất hiện ở `/api/*` test và 2 endpoint xóa ảnh.
- Upload behavior (cars/customers):
  - Nhận multi-file field `images`.
  - Filename xử lý bằng `secure_filename`, prefix UUID.
  - Lưu tại `static/uploads/cars` hoặc `static/uploads/customers`.
  - DB lưu path tương đối: `uploads/<module>/<filename>`.
  - Chưa có MIME whitelist/size limit explicit tại controller.

## 3) Endpoint inventory & contract baseline

> Ghi chú: vì hệ hiện tại chủ yếu server-rendered, “response schema” ở đây mô tả loại response thực tế (HTML redirect/render hoặc JSON).

### 3.1 Auth + Search

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/login` | GET, POST | No | Form `username`, `password`, `remember`; query `next` | HTML `login.html` hoặc Redirect | Validate form; check `User.check_password`; success -> login session + redirect `next`/dashboard; fail -> flash danger. |
| `/admin/logout` | GET | Yes | None | Redirect `/admin/login` | Logout session + flash info. |
| `/admin/search` | GET, POST (thực tế đọc `q` từ query) | Yes | Query `q` | HTML `search.html` | Nếu có `q`: gọi `Customer.search`, `Car.search`, `Transaction.search`, `User.search`; nếu rỗng trả list rỗng. |

### 3.2 Dashboard + System settings

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/` | GET | **No** (`login_required` đang comment) | None | HTML `dashboard.html` | Tính metric qua `get_metrics()`. |
| `/admin/dashboard` | GET | **No** | None | HTML `dashboard.html` | Alias route của dashboard. |
| `/admin/system` | GET, POST | **No** | Form `currency`,`theme`,`language` | HTML `system.html` hoặc Redirect | POST ghi setting qua `set_setting`, sau đó redirect về chính nó + flash success. GET đọc bằng `get_setting`. |

### 3.3 Users

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/user/new` | GET, POST | **No** | RegisterForm: username,email,password,role,status | HTML `user_new.html` hoặc Redirect | Check duplicate username/email; password hash bằng Werkzeug; create user. |
| `/admin/users` | GET | **No** | None | HTML `users.html` | Lấy `User.get_all()`, tách `admin_users` và `staff_users`. |
| `/admin/user/<int:user_id>/delete` | GET | **No** | Path `user_id` | Redirect `/admin/users` | Xóa user theo ID, flash success/fail. |
| `/admin/user/<user_id>` | GET, POST | **No** | UserUpdateForm + optional new password | HTML `user_detail.html` | Check duplicate với user khác; update fields; nếu có password thì hash lại. |

### 3.4 Cars

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/cars` | GET | Yes | None | HTML `cars.html` | Trả all cars + nhóm status/situation (available/awaiting/sold/refurbish...). |
| `/admin/car/new` | GET, POST | Yes | CarForm + files `images[]` | HTML `car_new.html` hoặc Redirect | Create car từ form; upload ảnh nhiều file; save DB `CarImage`. |
| `/admin/car/<car_id>` | GET, POST | Yes | CarForm + files `images[]` | HTML `car_detail.html` hoặc Redirect | Update car + append ảnh mới; trả 404 page nếu không có xe. |
| `/admin/car/<int:car_id>/delete` | GET | Yes | Path `car_id` | Redirect `/admin/cars` | Delete car. |
| `/admin/car/<int:car_id>/purchase/new` | GET, POST | Yes | TransactionForm (prefill `car_id`) | Redirect car detail | Tạo purchase transaction cho car. |
| `/admin/car/image/<int:image_id>/delete` | POST | Yes | Path `image_id` | JSON `{success: bool, error?: string}` | Xóa file vật lý + record DB; 404 nếu không có ảnh. |

### 3.5 Customers

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/customers` | GET | Yes | None | HTML `customers.html` | Trả all customers + `active_customers` có transaction. |
| `/admin/customer/new` | GET, POST | Yes | CustomerForm + files `images[]` | HTML `customer_new.html` hoặc Redirect | Create customer + upload multi-image tương tự cars. |
| `/admin/customer/<int:customer_id>` | GET, POST | Yes | CustomerForm + files `images[]` | HTML `customer_detail.html` | Update customer + append ảnh; 404 page nếu không tồn tại. |
| `/admin/customer/<int:customer_id>/purchase/new` | GET, POST | Yes | TransactionForm (prefill `customer_id`) | Redirect customer detail | Tạo purchase transaction cho customer. |
| `/admin/customer/<int:customer_id>/delete` | GET | Yes | Path `customer_id` | Redirect `/admin/customers` | Delete customer. |
| `/admin/customer/image/<int:image_id>/delete` | POST | Yes | Path `image_id` | JSON `{success: bool, error?: string}` | Xóa file vật lý + record DB; 404 nếu không có ảnh. |

### 3.6 Transactions

| Endpoint | Methods | Auth | Request | Response baseline | Business behavior |
|---|---|---|---|---|---|
| `/admin/transactions` | GET | Yes | None | HTML `transactions.html` | List transaction + nhóm deposited/paid + revenue aggregate. |
| `/admin/transaction/new` | GET, POST | Yes | TransactionForm + query prefill `customer_id`,`car_id` + repeated `item_name[]`,`item_price[]` | HTML `transaction_new.html` hoặc Redirect | Create transaction + add `TransactionItem` từ mảng item form. |
| `/admin/transaction/<transaction_id>` | GET, POST | Yes | TransactionForm + repeated item arrays | HTML `transaction_detail.html` | Update transaction; clear item cũ rồi add item mới. |
| `/admin/transaction/<int:transaction_id>/delete` | GET | Yes | Path `transaction_id` | Redirect `/admin/transactions` | Delete transaction. |

### 3.7 Utility API endpoints (JSON)

| Endpoint | Methods | Auth | Response baseline |
|---|---|---|---|
| `/api/test_admin` | GET | No | `200 {"status":"OK"}` |
| `/api/test_home` | GET | No | `200 {"status":"OK"}` |

## 4) Behavior gaps/notes cần carry sang phase sau

- Auth enforcement chưa đồng nhất: dashboard/system/users routes đang mở (login decorator bị comment).
- Nhiều hành động delete đang dùng `GET` thay vì `DELETE/POST`.
- Form/HTML flow là baseline hiện tại; khi migrate qua REST cần map thành request/response JSON tương đương behavior nghiệp vụ.
- Validation chủ yếu nằm trong WTForms; API Go cần explicit error envelope và field-level validation tương đương.

## 5) Traceability

- Tài liệu này hoàn thành DoD Task A2: “Có API contract baseline (OpenAPI draft hoặc Markdown spec)”.
- Dùng làm input cho B1/C2/C5/C7/E2.
