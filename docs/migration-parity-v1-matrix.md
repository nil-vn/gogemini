# v1 Feature Parity Matrix (Task A1)

Mục tiêu tài liệu này là **khóa phạm vi v1 parity** cho migration Go + Svelte 5, làm baseline bắt buộc cho các task tiếp theo.

## In-scope (bắt buộc có trong v1)

| Feature | Route hiện tại (legacy Flask) | Expected behavior (v1 parity) |
|---|---|---|
| Auth (login/logout) | `GET/POST /login`, `GET /logout` | Người dùng legacy đăng nhập thành công; logout huỷ session/cookie và điều hướng về trang login; route admin yêu cầu xác thực. |
| Dashboard | `GET /admin`, `GET /admin/dashboard` | Hiển thị metrics chính theo cùng định nghĩa nghiệp vụ như hệ cũ, không thay đổi semantics số liệu trong v1. |
| Users CRUD | `GET /admin/users`, `GET /admin/users/new`, `GET /admin/users/edit/<id>` + các API tương đương | Đầy đủ list/detail/create/update/delete; hỗ trợ pagination/sort/filter theo baseline; validation và thông báo lỗi tương đương legacy. |
| Cars CRUD | `GET /admin/cars`, `GET /admin/cars/new`, `GET /admin/cars/edit/<id>` + các API tương đương | Đầy đủ list/detail/create/update/delete; giữ behavior nghiệp vụ hiện tại (trường bắt buộc, định dạng dữ liệu, thứ tự mặc định). |
| Customers CRUD | `GET /admin/customers`, `GET /admin/customers/new`, `GET /admin/customers/edit/<id>` + các API tương đương | Đầy đủ list/detail/create/update/delete; parity validation và list behavior như hệ cũ. |
| Transactions CRUD | `GET /admin/transactions`, `GET /admin/transactions/new`, `GET /admin/transactions/edit/<id>` + các API tương đương | Đầy đủ list/detail/create/update/delete; giữ đúng quy tắc nghiệp vụ giao dịch hiện có. |
| Search | `GET /admin/search` (và endpoint tìm kiếm liên quan) | Kết quả tìm kiếm tương đương legacy cho cùng input/query, bao gồm default sort/filter theo baseline audit. |
| System settings | `GET /admin/system-settings` (và endpoint cấu hình liên quan) | Đọc/cập nhật được currency/theme/language; có validation key/value; chỉ admin có quyền truy cập. |
| Upload (car/customer image) | `POST /admin/upload` và luồng upload trong form cars/customers | Upload ảnh thành công cho cars/customers, validate file cơ bản, trả path/url hợp lệ và render preview đúng. |

## Out-of-scope cho v1 (không chặn cutover v1 nếu chưa làm)

- Nâng cao accessibility beyond baseline parity (keyboard navigation nâng cao, contrast audit đầy đủ).
- Tối ưu hiệu năng chuyên sâu (beyond baseline p95 theo phase 4).
- Refactor lớn UI/UX hoặc thay đổi luồng nghiệp vụ so với legacy.
- Mở rộng tính năng mới ngoài danh sách in-scope phía trên.

## Ghi chú khóa phạm vi

- Tài liệu này phục vụ Task **A1** và là input bắt buộc cho A2/A4/E2.
- Mọi yêu cầu mới ngoài bảng in-scope phải đi qua quy trình change request và được chốt ở task khác, không gộp vào v1 parity.
