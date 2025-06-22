# 🚀 Microservices E-Commerce Platform

Nền tảng thương mại điện tử được xây dựng với kiến trúc microservices sử dụng Go (Golang), bao gồm 3 service chính: User Service, Product Service và Order Service.

## 🏗️ Kiến trúc Hệ thống

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Service  │    │ Product Service │    │  Order Service  │
│     :8080       │    │      :8082      │    │      :8083      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         v                       v                       v
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   PostgreSQL    │    │   PostgreSQL    │
│     :5432       │    │     :5433       │    │     :5431       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                v
                       ┌─────────────────┐
                       │      Redis      │
                       │      :6379      │
                       └─────────────────┘
```

## 📋 Tổng quan Services

### 🔐 User Service (Port: 8080)
- **Chức năng:** Quản lý người dùng, xác thực và phân quyền
- **Database:** PostgreSQL (Port: 5432)
- **Tính năng:**
  - Đăng ký/Đăng nhập người dùng
  - Quản lý profile người dùng
  - JWT Authentication
  - Phân quyền (buyer, partner, admin)

### 🛍️ Product Service (Port: 8082)
- **Chức năng:** Quản lý sản phẩm và danh mục
- **Database:** PostgreSQL (Port: 5433) + Redis Cache (Port: 6379)
- **Tính năng:**
  - CRUD sản phẩm
  - Tìm kiếm sản phẩm
  - Cache với Redis

### 🛒 Order Service (Port: 8083)
- **Chức năng:** Quản lý đơn hàng và giao dịch
- **Database:** PostgreSQL (Port: 5431)
- **Tính năng:**
  - Tạo và quản lý đơn hàng
  - Theo dõi trạng thái đơn hàng
  - Quản lý hoa hồng đối tác
  - Thông báo đơn hàng

## ⚡ Quick Start

### 1. Prerequisites
Đảm bảo bạn đã cài đặt:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Git](https://git-scm.com/)

### 2. Clone Repository
```bash
git clone <repository-url>
cd SF
```

### 3. Khởi chạy toàn bộ hệ thống
```bash
# Khởi chạy tất cả services
docker-compose up --build

# Hoặc chạy trong background
docker-compose up -d --build
```

### 4. Kiểm tra trạng thái services
```bash
# Xem logs của tất cả services
docker-compose logs

# Xem logs của service cụ thể
docker-compose logs user_service
docker-compose logs product_service
docker-compose logs order_service

# Kiểm tra trạng thái containers
docker-compose ps
```

### 5. Dừng hệ thống
```bash
# Dừng tất cả services
docker-compose down

# Dừng và xóa volumes (dữ liệu sẽ bị mất)
docker-compose down -v
```

### 6. Kiểm tra .dockerignore
```bash
# Xem nội dung .dockerignore của từng service
cat User_Service/.dockerignore
cat Product_Service/.dockerignore  
cat Order_Service/.dockerignore

# Nếu file .dockerignore bị thiếu hoặc không phù hợp, 
# hãy tạo/sửa chúng thay vì xóa đi
```

## 🌐 API Endpoints - Chi tiết

### 🔐 User Service - `http://localhost:8080`

#### 🌍 Public Endpoints
```bash
# Đăng ký người dùng mới
POST /api/v1/register
Content-Type: application/json
{
  "name": "Nguyễn Văn A",
  "email": "user@example.com", 
  "password": "password123",
  "role": "buyer"  // buyer, partner, admin
}

```

```bash
# Đăng nhập
POST /api/v1/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

```

#### 🔒 Protected Endpoints
```bash
# Lấy thông tin profile cá nhân
GET /api/v1/get-profile
Authorization: Bearer <jwt-token>

```

```bash
# Cập nhật thông tin profile
PATCH /api/v1/update-profile
Authorization: Bearer <jwt-token>
Content-Type: application/json
{
  "name": "Nguyễn Văn B"
}


```

```bash
# Xóa tài khoản
DELETE /api/v1/delete-profile
Authorization: Bearer <jwt-token>


```

#### 👑 Admin Endpoints
```bash
# Khóa/Mở khóa người dùng
POST /api/admin/toggle-user-lock/{user_id}
Authorization: Bearer <admin-token>

```

```bash
# Lấy danh sách tất cả người dùng
GET /api/admin/get-all-user
Authorization: Bearer <admin-token>

```

---

### 🛍️ Product Service - `http://localhost:8082`

#### 🌍 Public Endpoints  
```bash
# Lấy tất cả sản phẩm (có cache)
GET /api/v1/products

```

```bash
# Lấy chi tiết sản phẩm theo ID
GET /api/v1/product-id/{product_id}

```

```bash
# Tìm kiếm sản phẩm theo tên
GET /api/v1/product-name/{product_name}

```

#### 👑 Admin Endpoints
```bash
# Tạo sản phẩm mới
POST /api/admin/products
Authorization: Bearer <admin-token>
Content-Type: application/json
{
  "name": "iPhone 15 Pro",
  "description": "Smartphone cao cấp với chip A17 Pro",
  "price": 29999000,
  "image": "https://example.com/iphone15pro.jpg"
}

```

```bash
# Cập nhật sản phẩm
PATCH /api/admin/products/{product_id}
Authorization: Bearer <admin-token>
Content-Type: application/json
{
  "name": "iPhone 15 Pro (Updated)",
  "price": 27999000
}

```

```bash
# Xóa sản phẩm
DELETE /api/admin/products/{product_id}
Authorization: Bearer <admin-token>

```

---

### 🛒 Order Service - `http://localhost:8083`

#### 🛍️ Buyer Endpoints
```bash
# Tạo đơn hàng mới
POST /api/buyer/orders
Authorization: Bearer <buyer-token>
Content-Type: application/json
{
  "address": "123 Nguyễn Văn A, Quận 1, TP.HCM",
  "order_items": [
    {
      "product_id": "product-uuid-1",
      "partner_id": "partner-uuid-1", 
      "quantity": 2
    },
    {
      "product_id": "product-uuid-2",
      "partner_id": "partner-uuid-2",
      "quantity": 1
    }
  ]
}

```

```bash
# Lấy đơn hàng theo trạng thái
GET /api/buyer/orders/status?status={status}
Authorization: Bearer <buyer-token>
# Status: pending, confirmed, complete, cancel

```

```bash
# Cập nhật trạng thái đơn hàng (Buyer cancel)
PUT /api/buyer/orders/{orderItemId}
Authorization: Bearer <buyer-token>

```

#### 🤝 Partner Endpoints
```bash
# Lấy đơn hàng của partner theo trạng thái
GET /api/partner/orders/{status}
Authorization: Bearer <partner-token>
# Status: pending, confirmed, complete, cancel

```

```bash
# Xác nhận đơn hàng (Partner confirm)
PUT /api/partner/orders/{orderItemId}
Authorization: Bearer <partner-token>

```

```bash
# Xem hoa hồng của partner
GET /api/partner/commission
Authorization: Bearer <partner-token>

```

#### 👑 Admin Endpoints  
```bash
# Lấy tất cả đơn hàng
GET /api/admin/all-orders
Authorization: Bearer <admin-token>

```

```bash
# Lấy chi tiết đơn hàng
GET /api/admin/orders/{orderId}
Authorization: Bearer <admin-token>

```

```bash
# Xem hoa hồng theo OrderItem
GET /api/admin/commission/{orderItemId}
Authorization: Bearer <admin-token>

```

```bash
# Xem tổng hoa hồng theo Partner
GET /api/admin/commission-pnid/{partnerId}
Authorization: Bearer <admin-token>

```

---

## 📊 HTTP Response Codes & Error Handling

### Success Codes
- **200 OK**: Request thành công
- **201 Created**: Tạo resource thành công
- **204 No Content**: Xóa thành công

### Client Error Codes
- **400 Bad Request**: Dữ liệu không hợp lệ
- **401 Unauthorized**: Token không hợp lệ hoặc hết hạn
- **403 Forbidden**: Không có quyền truy cập
- **404 Not Found**: Resource không tồn tại
- **409 Conflict**: Conflict dữ liệu (email đã tồn tại)
- **422 Unprocessable Entity**: Validation error


---

## 🗄️ Database Schema

### User Service Database
- **users**: Thông tin người dùng, authentication
- **roles**: Phân quyền hệ thống

### Product Service Database
- **products**: Thông tin sản phẩm

### Order Service Database
- **orders**: Thông tin đơn hàng
- **order_items**: Chi tiết từng sản phẩm trong đơn hàng
- **partner_commissions**: Hoa hồng đối tác
- **notifications**: Thông báo hệ thống

## 🔧 Environment Variables

### Common Variables
```env
JWT_SECRET=qwertyuiopiuytsdfghjkllkjhgfddfghjklkjhgfddfghjkkj
GIN_MODE=release
```

### User Service
```env
DB_HOST=user_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=user_service
SERVER_PORT=8080
JWT_EXPIRES_IN=24
```

### Product Service
```env
DB_HOST=product_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=@BCxyz123
DB_NAME=product_service_db
REDIS_HOST=redis
REDIS_PORT=6379
SERVER_PORT=8082
```

### Order Service
```env
DB_HOST=order_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=@BCxyz123
DB_NAME=order_service_db
SERVER_PORT=8083
```


## 📚 Tech Stack

- **Backend Framework:** Gin (Go)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Authentication:** JWT
- **Containerization:** Docker & Docker Compose
- **Architecture:** Microservices
