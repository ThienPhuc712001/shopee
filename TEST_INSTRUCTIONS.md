# 🎯 HƯỚNG DẪN TEST HỆ THỐNG E-COMMERCE

## ✅ DATABASE ĐÃ ĐƯỢC SEED THÀNH CÔNG

### 📊 Dữ liệu đã tạo:

**Users (3 tài khoản):**
- ✅ customer@example.com / customer123 (Customer)
- ✅ seller@example.com / seller123 (Seller)
- ✅ admin@example.com / admin123 (Admin) - Đã tồn tại

**Categories (8 danh mục):**
- ✅ Electronics, Fashion, Home & Living, Books, Sports, Beauty, Toys & Games, Watches

**Shop (1 cửa hàng):**
- ✅ Tech Store (của seller)

**Products (8 sản phẩm):**
- ✅ Premium Wireless Headphones - $199.99 (was $299.99)
- ✅ Smart Watch Pro - $349.99 (was $449.99)
- ✅ Professional Camera - $899.99 (was $1199.99)
- ✅ Gaming Laptop - $1299.99 (was $1599.99)
- ✅ Bluetooth Speaker - $79.99 (was $129.99)
- ✅ Mechanical Keyboard - $129.99 (was $179.99)
- ✅ USB-C Hub - $49.99 (was $79.99)
- ✅ Wireless Charger - $39.99 (was $59.99)

**Coupons (3 mã giảm giá):**
- ✅ WELCOME10 - Giảm 10% (đơn từ $50, tối đa $50)
- ✅ SUMMER20 - Giảm 20% (đơn từ $100, tối đa $100)
- ✅ FREESHIP - Giảm $5.99 phí ship (đơn từ $30)

---

## 🧪 HƯỚNG DẪN TEST TỪNG LUỒNG

### 1. Test Homepage (Không cần login)

**URL:** http://localhost:3000

**Test:**
- ✅ Xem Hero Banner (tự động slide)
- ✅ Xem Categories (8 danh mục)
- ✅ Xem Flash Sale (với countdown timer)
- ✅ Xem Featured Products (4 sản phẩm nổi bật)
- ✅ Xem Recommended Products

**Expected:**
- Banner tự động chuyển sau mỗi 5 giây
- Hiển thị đầy đủ 8 categories
- Hiển thị 4 sản phẩm từ database
- Có indicator "🟢 Backend Connected"

---

### 2. Test Product Listing (Không cần login)

**URL:** http://localhost:3000/products

**Test:**
- ✅ Xem danh sách 8 sản phẩm
- ✅ Filter theo category
- ✅ Sort theo giá
- ✅ Pagination

**Expected:**
- Hiển thị 8 sản phẩm từ database
- Grid responsive (4 columns desktop, 2 mobile)
- Product cards có hover effect

---

### 3. Test Product Detail (Không cần login)

**URL:** http://localhost:3000/products/1

**Test:**
- ✅ Xem chi tiết sản phẩm
- ✅ Xem giá, discount
- ✅ Xem stock status
- ✅ Add to cart

**Expected:**
- Hiển thị đầy đủ thông tin sản phẩm
- Giá hiển thị đúng với database
- Nút "Add to Cart" hoạt động

---

### 4. Test Authentication Flow

#### Register User Mới:
**URL:** http://localhost:3000/register

**Test:**
- ✅ Register với email mới
- ✅ Login với tài khoản vừa tạo

**Expected:**
- Register thành công
- Tự động login sau register
- User được lưu vào database

#### Login:
**URL:** http://localhost:3000/login

**Test với các tài khoản:**
```
Email: customer@example.com
Password: customer123

Email: seller@example.com
Password: seller123

Email: admin@example.com
Password: admin123
```

**Expected:**
- Login thành công
- Header hiển thị user menu
- Cart badge hiển thị số lượng

---

### 5. Test Cart Flow (Cần login)

**URL:** http://localhost:3000/cart

**Test:**
- ✅ Add sản phẩm vào cart từ Product Detail
- ✅ Xem cart với số lượng cập nhật
- ✅ Update quantity
- ✅ Remove item
- ✅ Clear cart

**Expected:**
- Cart badge trên header cập nhật
- Subtotal tính đúng
- Quantity controls hoạt động

---

### 6. Test Checkout Flow (Cần login)

**URL:** http://localhost:3000/checkout

**Test:**
- ✅ Nhập shipping address
- ✅ Chọn payment method (COD)
- ✅ Áp dụng coupon code (WELCOME10, SUMMER20, FREESHIP)
- ✅ Place order

**Expected:**
- Order được tạo trong database
- Cart được clear
- Hiển thị order confirmation

---

### 7. Test User Dashboard (Cần login)

**URL:** http://localhost:3000/account

**Test:**
- ✅ Xem account info
- ✅ Xem orders history
- ✅ Xem wishlist
- ✅ Xem addresses

**Expected:**
- Hiển thị thông tin user
- Danh sách orders (nếu có)

---

### 8. Test Admin Panel (Cần admin login)

**URL:** http://localhost:3000/admin

**Login với:** admin@example.com / admin123

**Test:**
- ✅ Xem dashboard stats
- ✅ Product management
- ✅ Order management
- ✅ User management
- ✅ Analytics

**Expected:**
- Hiển thị stats từ database
- CRUD operations hoạt động

---

### 9. Test Help & Info Pages (Không cần login)

**URLs:**
- http://localhost:3000/help
- http://localhost:3000/track-order
- http://localhost:3000/about
- http://localhost:3000/contact
- http://localhost:3000/faq
- http://localhost:3000/shipping
- http://localhost:3000/terms
- http://localhost:3000/privacy
- http://localhost:3000/returns

**Test:**
- ✅ Tất cả pages load thành công
- ✅ Không còn warnings trong console
- ✅ Design responsive

**Expected:**
- Tất cả pages hiển thị nội dung
- Form contact hoạt động
- FAQ search hoạt động

---

## 🔧 KIỂM TRA BACKEND API

### Test với curl:

```bash
# Health check
curl http://localhost:8080/health

# Get products
curl http://localhost:8080/api/products

# Get product by ID
curl http://localhost:8080/api/products/1

# Get categories
curl http://localhost:8080/api/categories

# Get shops
curl http://localhost:8080/api/shops-list

# Get featured products
curl http://localhost:8080/api/products/featured

# Search products
curl "http://localhost:8080/api/products/search?keyword=laptop"

# Get shipping calculate
curl -X POST http://localhost:8080/api/shipping/calculate \
  -H "Content-Type: application/json" \
  -d "{\"from_city\":\"HCMC\",\"to_city\":\"HN\",\"weight\":1.5,\"shipping_method\":\"standard\"}"
```

---

## 🎯 TEST CHECKLIST

### Public Pages:
- [ ] Homepage loads with all sections
- [ ] Product listing shows 8 products
- [ ] Product detail shows correct info
- [ ] Categories filter works
- [ ] Search works
- [ ] All info pages load (help, about, contact, etc.)

### Authentication:
- [ ] Register new user
- [ ] Login with customer account
- [ ] Login with seller account
- [ ] Login with admin account
- [ ] Logout works
- [ ] User menu displays in header

### Cart & Checkout:
- [ ] Add to cart works
- [ ] Cart displays items
- [ ] Update quantity works
- [ ] Remove item works
- [ ] Checkout form works
- [ ] Apply coupon code works
- [ ] Place order creates order in DB

### User Dashboard:
- [ ] Account page shows user info
- [ ] Orders page shows order history
- [ ] Wishlist page works
- [ ] Addresses page works

### Admin Panel:
- [ ] Admin dashboard shows stats
- [ ] Product management works
- [ ] Order management works
- [ ] User management works

### API Endpoints:
- [ ] All public endpoints return 200
- [ ] Auth endpoints return 401 without token
- [ ] Auth endpoints work with valid token
- [ ] Database queries return correct data

---

## 🐛 NẾU GẶP LỖI

### Frontend không load được products:
1. Kiểm tra backend đang chạy: `curl http://localhost:8080/health`
2. Kiểm tra console log trên browser
3. Verify API proxy trong vite.config.ts

### Login không hoạt động:
1. Verify database có users: Check seed log
2. Try register user mới
3. Check browser console for errors

### Cart không hoạt động:
1. Phải login trước khi add to cart
2. Check network tab cho API calls
3. Verify token trong localStorage

---

## 📞 CREDENTIALS TÓM TẮT

| Role | Email | Password |
|------|-------|----------|
| Customer | customer@example.com | customer123 |
| Seller | seller@example.com | seller123 |
| Admin | admin@example.com | admin123 |

---

**Happy Testing! 🎉**

**Last Updated:** March 13, 2026
**Status:** ✅ READY FOR TESTING
