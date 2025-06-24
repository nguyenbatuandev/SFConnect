# 🚀 SFConnect - Nền tảng Thương mại Điện tử Microservices

Hệ thống thương mại điện tử hiện đại được xây dựng với kiến trúc microservices sử dụng **Go (Golang)**, bao gồm 3 service độc lập: **User Service**, **Product Service** và **Order Service**. Hệ thống hỗ trợ đầy đủ chức năng từ quản lý người dùng, sản phẩm đến xử lý đơn hàng và hoa hồng đối tác.

## 🏗️ Kiến trúc Hệ thống

```
                        🌐 Client Applications
                               │
                               v
                    ┌─────────────────────┐
                    │    API Gateway      │
                    │   (Load Balancer)   │
                    └─────────────────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              v                v                v
    ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
    │   User Service  │ │ Product Service │ │  Order Service  │
    │     :8080       │ │      :8082      │ │      :8083      │
    │   Auth & Users  │ │  Products &     │ │  Orders &       │
    │   Management    │ │    Catalog      │ │  Commissions    │
    └─────────────────┘ └─────────────────┘ └─────────────────┘
              │                │                │
              │                │                │
              v                v                v
    ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
    │   PostgreSQL    │ │   PostgreSQL    │ │   PostgreSQL    │
    │ user_service_db │ │product_service  │ │order_service_db │
    │     :5432       │ │     :5433       │ │     :5431       │
    └─────────────────┘ └─────────────────┘ └─────────────────┘
                               │
                               v
                      ┌─────────────────┐
                      │   Redis Cache   │
                      │     :6379       │
                      └────────────────┘
```

## 📋 Tổng quan Services

### 🔐 User Service (Port: 8080)
- **Chức năng chính:** Quản lý người dùng, xác thực và phân quyền
- **Database:** PostgreSQL (Port: 5432)
- **Công nghệ:** Gin Framework, JWT Authentication, GORM
- **Tính năng:**
  - ✅ Đăng ký/Đăng nhập người dùng với JWT
  - ✅ Quản lý thông tin cá nhân (CRUD profile)
  - ✅ Phân quyền theo vai trò: `buyer`, `partner`, `admin`
  - ✅ Khóa/mở khóa tài khoản (Admin)
  - ✅ Xóa tài khoản người dùng

### 🛍️ Product Service (Port: 8082)
- **Chức năng chính:** Quản lý sản phẩm và tối ưu hóa hiệu suất
- **Database:** PostgreSQL (Port: 5433) + Redis Cache (Port: 6379)
- **Công nghệ:** Gin Framework, Redis Caching, GORM
- **Tính năng:**
  - ✅ CRUD sản phẩm (Create, Read, Update, Delete)
  - ✅ Tìm kiếm sản phẩm theo tên và ID
  - ✅ Cache thông minh với Redis để tăng tốc truy vấn
  - ✅ API public cho khách hàng và API admin để quản lý
  - ✅ Xử lý hình ảnh và mô tả sản phẩm

### 🛒 Order Service (Port: 8083) - **Được cập nhật mới nhất**
- **Chức năng chính:** Xử lý đơn hàng phức tạp và quản lý hoa hồng đối tác
- **Database:** PostgreSQL (Port: 5431)
- **Công nghệ:** Gin Framework, Service-to-Service Communication, GORM
- **Tính năng nâng cao:**
  - ✅ **Quản lý đơn hàng đa cấp:**
    - Tạo đơn hàng với nhiều sản phẩm từ nhiều đối tác khác nhau
    - Mỗi OrderItem có thể có trạng thái riêng biệt
    - Theo dõi chi tiết từng sản phẩm trong đơn hàng
  - ✅ **Hệ thống trạng thái đơn hàng:**
    - `pending`: Chờ xác nhận từ đối tác
    - `confirmed`: Đối tác đã xác nhận
    - `complete`: Hoàn thành giao hàng
    - `cancel`: Đã hủy đơn hàng
  - ✅ **Quản lý hoa hồng thông minh:**
    - Tự động tính hoa hồng cho đối tác theo từng OrderItem
    - Theo dõi tổng hoa hồng theo đối tác
    - Báo cáo hoa hồng chi tiết cho admin
  - ✅ **Phân quyền xử lý:**
    - **Buyer**: Tạo đơn hàng, theo dõi trạng thái, hủy đơn
    - **Partner**: Xác nhận đơn hàng, xem hoa hồng của mình
    - **Admin**: Quản lý toàn bộ đơn hàng và hoa hồng
  - ✅ **Integration với các service khác:**
    - Gọi User Service để verify thông tin người dùng
    - Gọi Product Service để lấy thông tin sản phẩm và giá
  - ✅ **Hệ thống thông báo:** Thông báo trạng thái đơn hàng realtime

## ⚡ Quick Start

### 1. Prerequisites
Đảm bảo bạn đã cài đặt đầy đủ:
- [Docker](https://www.docker.com/) (version 20.0+)
- [Docker Compose](https://docs.docker.com/compose/) (version 2.0+)
- [Git](https://git-scm.com/)
- [Postman](https://www.postman.com/) (để test API - tùy chọn)

### 2. Clone Repository
```bash
git clone <repository-url>
cd SFConnect
ls -la  # Kiểm tra cấu trúc thư mục
```

### 3. Cấu hình Environment (Tùy chọn)
```bash
# Các biến môi trường đã được cấu hình sẵn trong docker-compose.yml
# Bạn có thể tùy chỉnh theo nhu cầu
```

### 4. Khởi chạy toàn bộ hệ thống
```bash
# Build và khởi chạy tất cả services
docker-compose up --build

# Hoặc chạy trong background (khuyến nghị)
docker-compose up -d --build

# Chỉ khởi chạy một service cụ thể
docker-compose up user_service
docker-compose up product_service  
docker-compose up order_service
```

### 5. Kiểm tra trạng thái hệ thống
```bash
# Xem trạng thái tất cả containers
docker-compose ps

# Xem logs của tất cả services
docker-compose logs -f

# Xem logs của service cụ thể
docker-compose logs -f user_service
docker-compose logs -f product_service
docker-compose logs -f order_service

# Xem logs database
docker-compose logs -f user_db
docker-compose logs -f product_db
docker-compose logs -f order_db
docker-compose logs -f redis
```

### 6. Health Check các Services
```bash
# Kiểm tra User Service
curl http://localhost:8080/api/v1/health

# Kiểm tra Product Service  
curl http://localhost:8082/api/v1/health

# Kiểm tra Order Service
curl http://localhost:8083/api/v1/health

# Kiểm tra kết nối Database
docker-compose exec user_db psql -U postgres -d user_service -c "SELECT 1;"
docker-compose exec product_db psql -U postgres -d product_service_db -c "SELECT 1;"
docker-compose exec order_db psql -U postgres -d order_service_db -c "SELECT 1;"

# Kiểm tra Redis
docker-compose exec redis redis-cli ping
```

### 7. Dừng hệ thống
```bash
# Dừng tất cả services (giữ nguyên data)
docker-compose down

# Dừng và xóa volumes (⚠️ CẢNH BÁO: sẽ mất hết dữ liệu)
docker-compose down -v

# Dừng và xóa images (để rebuild hoàn toàn)
docker-compose down --rmi all
```

### 8. Troubleshooting
```bash
# Rebuild một service cụ thể
docker-compose build --no-cache user_service

# Restart một service
docker-compose restart order_service

# Xem resource usage
docker stats

# Clean up Docker system
docker system prune -a
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

### 🛒 Order Service - `http://localhost:8083` - **🆕 CẬP NHẬT MỚI**

#### 🛍️ Buyer Endpoints
```bash
# 🆕 Tạo đơn hàng mới (Hỗ trợ đa đối tác)
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
# Response: Tạo 1 Order với 2 OrderItems riêng biệt
```

```bash
# 🆕 Lấy đơn hàng theo trạng thái (Cải tiến)
GET /api/buyer/orders/{status}
Authorization: Bearer <buyer-token>
# Status: pending | confirmed | complete | cancel
# Response: Danh sách OrderItems với thông tin đầy đủ
```

```bash
# 🆕 Cập nhật trạng thái đơn hàng (Buyer cancel OrderItem)
PUT /api/buyer/orders/{orderItemId}
Authorization: Bearer <buyer-token>
# Chỉ buyer có thể cancel các OrderItem của mình
```

```bash
# 🆕 Lấy danh sách đối tác khả dụng
GET /api/buyer/partners
Authorization: Bearer <buyer-token>
# Response: Danh sách partners đang hoạt động
```

#### 🤝 Partner Endpoints - **🔥 NÂNG CẤP MẠNH**
```bash
# 🆕 Lấy OrderItems của partner theo trạng thái
GET /api/partner/orders/{status}
Authorization: Bearer <partner-token>
# Status: pending | confirmed | complete | cancel
# Response: Chỉ OrderItems liên quan đến partner này
```

```bash
# 🆕 Xác nhận OrderItem (Partner confirm)
PUT /api/partner/orders/{orderItemId}
Authorization: Bearer <partner-token>
# Partner chỉ có thể confirm OrderItems của mình
# Tự động tính hoa hồng khi confirm
```

```bash
# 🆕 Xem hoa hồng của partner (Chi tiết)
GET /api/partner/commission
Authorization: Bearer <partner-token>
# Response: Tổng hoa hồng + danh sách OrderItems đã tạo hoa hồng
```

#### 👑 Admin Endpoints - **🚀 SIÊU MẠNH**
```bash
# 🆕 Lấy tất cả đơn hàng (Phân trang)
GET /api/admin/all-orders?page=1&limit=10
Authorization: Bearer <admin-token>
# Response: Danh sách Orders với OrderItems nested
```

```bash
# 🆕 Lấy chi tiết đơn hàng (Deep Info)
GET /api/admin/orders/{orderId}
Authorization: Bearer <admin-token>
# Response: Order + OrderItems + User Info + Product Info
```

```bash
# 🆕 Xem hoa hồng theo OrderItem
GET /api/admin/commission/{orderItemId}
Authorization: Bearer <admin-token>
# Response: Chi tiết hoa hồng của OrderItem cụ thể
```

```bash
# 🆕 Xem tổng hoa hồng theo Partner
GET /api/admin/commission-pnid/{partnerId}
Authorization: Bearer <admin-token>
# Response: Tổng hoa hồng + breakdown theo OrderItem
```

#### 🔥 Order Service - Luồng xử lý mới:
```
1. Buyer tạo Order → Tạo nhiều OrderItems (status: pending)
2. Từng Partner nhận thông báo → Confirm OrderItems của minha 
3. OrderItem confirmed → Tự động tính hoa hồng
4. Buyer nhận hàng → Partner đánh dấu complete
5. Admin có thể theo dõi toàn bộ luồng + hoa hồng
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

## 🗄️ Database Schema - Chi tiết

### 🔐 User Service Database (`user_service`)
```sql
-- Users table
Table: users
- id (UUID, Primary Key)
- name (VARCHAR(255), NOT NULL)  
- email (VARCHAR(255), UNIQUE, NOT NULL)
- password_hash (VARCHAR(255), NOT NULL)
- role (VARCHAR(20), DEFAULT 'buyer') -- buyer, partner, admin
- is_locked (BOOLEAN, DEFAULT false)
- created_at (TIMESTAMP, DEFAULT NOW())
- updated_at (TIMESTAMP, DEFAULT NOW())
```

### 🛍️ Product Service Database (`product_service_db`)
```sql
-- Products table
Table: products
- id (UUID, Primary Key)
- name (VARCHAR(255), NOT NULL)
- description (TEXT)
- price (DECIMAL(10,2), NOT NULL, DEFAULT 0.00)
- image (VARCHAR(500)) -- URL to image
- created_at (TIMESTAMP, DEFAULT NOW())
- updated_at (TIMESTAMP, DEFAULT NOW())
- is_active (BOOLEAN, DEFAULT true)
```

### 🛒 Order Service Database (`order_service_db`) - **🆕 CẬP NHẬT**
```sql
-- Orders table (Đơn hàng chính)
Table: orders
- id (UUID, Primary Key)
- buyer_id (UUID, NOT NULL) -- Reference to User Service
- create_time (TIMESTAMP, DEFAULT NOW())
- total_price (DECIMAL(10,2), DEFAULT 0.00) -- Tự động tính từ OrderItems
- address (VARCHAR(255), NOT NULL)

-- Order Items table (Chi tiết từng sản phẩm trong đơn hàng)
Table: order_items  
- id (UUID, Primary Key)
- order_id (UUID, NOT NULL, Foreign Key → orders.id)
- product_id (UUID, NOT NULL) -- Reference to Product Service
- partner_id (UUID, NOT NULL) -- Reference to User Service (role=partner)
- quantity (INTEGER, NOT NULL, DEFAULT 1)
- price (DECIMAL(10,2), NOT NULL, DEFAULT 0.00) -- Giá tại thời điểm đặt hàng
- status (VARCHAR(20), DEFAULT 'pending') -- pending, confirmed, complete, cancel

-- Partner Commissions table (Hoa hồng đối tác)
Table: partner_commissions
- id (UUID, Primary Key)
- partner_id (UUID, NOT NULL) -- Reference to User Service
- order_item_id (UUID, NOT NULL, Foreign Key → order_items.id)
- commission_amount (DECIMAL(10,2), NOT NULL, DEFAULT 0.00)
- commission_rate (DECIMAL(5,2), DEFAULT 10.00) -- % hoa hồng
- created_at (TIMESTAMP, DEFAULT NOW())
- status (VARCHAR(20), DEFAULT 'pending') -- pending, paid

-- Notifications table (Thông báo hệ thống)
Table: notifications
- id (UUID, Primary Key)
- user_id (UUID, NOT NULL) -- Reference to User Service
- order_id (UUID) -- Reference to orders.id (optional)
- message (TEXT, NOT NULL)
- type (VARCHAR(50)) -- order_created, order_confirmed, etc.
- is_read (BOOLEAN, DEFAULT false)
- created_at (TIMESTAMP, DEFAULT NOW())
```

### 🔄 Relationships & Constraints
```sql
-- Order Service relationships
ALTER TABLE order_items 
ADD CONSTRAINT fk_order_items_order 
FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE partner_commissions 
ADD CONSTRAINT fk_commissions_order_item 
FOREIGN KEY (order_item_id) REFERENCES order_items(id) ON DELETE CASCADE;

-- Indexes for performance
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX idx_order_items_status ON order_items(status);
CREATE INDEX idx_order_items_partner_id ON order_items(partner_id);
CREATE INDEX idx_commissions_partner_id ON partner_commissions(partner_id);
```

## 🔧 Environment Variables - Cấu hình chi tiết

### 🌐 Common Variables (Tất cả services)
```env
# JWT Configuration
JWT_SECRET=qwertyuiopiuytsdfghjkllkjhgfddfghjklkjhgfddfghjkkj
JWT_EXPIRES_IN=24  # hours

# Application Mode
GIN_MODE=release  # debug, release, test
```

### 🔐 User Service Environment
```env
# Database Configuration
DB_HOST=user_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=user_service

# Server Configuration  
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# Security
BCRYPT_COST=12
```

### 🛍️ Product Service Environment
```env
# Database Configuration
DB_HOST=product_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=@BCxyz123
DB_NAME=product_service_db

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_EXPIRATION=3600  # seconds

# Server Configuration
SERVER_PORT=8082
```

### 🛒 Order Service Environment - **🆕 CẬP NHẬT**
```env
# Database Configuration
DB_HOST=order_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=@BCxyz123
DB_NAME=order_service_db

# Server Configuration
SERVER_PORT=8083

# External Services URLs
USER_SERVICE_URL=http://user_service:8080
PRODUCT_SERVICE_URL=http://product_service:8082

# Business Logic
DEFAULT_COMMISSION_RATE=10.00  # 10% commission
ORDER_TIMEOUT=3600  # seconds
```

### 🐳 Docker Compose Override
```yaml
# Tạo file docker-compose.override.yml để tùy chỉnh
version: '3.8'
services:
  order_service:
    environment:
      - DEFAULT_COMMISSION_RATE=15.00
      - GIN_MODE=debug
    ports:
      - "8084:8083"  # Đổi port nếu cần
```


## 📚 Tech Stack

- **Backend Framework:** Gin (Go)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Authentication:** JWT
- **Containerization:** Docker & Docker Compose
- **Architecture:** Microservices
