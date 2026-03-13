# 🚀 QUICK START - TEST LUỒNG MUA HÀNG

## ⚡ START HERE - 5 PHÚT ĐỂ TEST FULL LUỒNG

### 1. Mở trình duyệt
```
http://localhost:3000
```

### 2. Login
```
URL: http://localhost:3000/login
Email: customer@example.com
Password: customer123
```

### 3. Thêm sản phẩm vào cart
```
1. Vào: http://localhost:3000/products
2. Click "Premium Wireless Headphones"
3. Click "Add to Cart"
4. Click "Premium Wireless Headphones"
5. Click "Add to Cart" (thêm 1 cái nữa cho tổng > $50)
```

### 4. Xem cart
```
URL: http://localhost:3000/cart
Click: "Proceed to Checkout"
```

### 5. Checkout - 3 bước

**Bước 1 - Shipping:**
```
Full Name: John Doe
Phone: 0901234567
City: Ho Chi Minh City
Address: 123 Nguyen Hue
District: District 1
→ Click "Continue to Payment"
```

**Bước 2 - Payment:**
```
Chọn: Cash on Delivery (COD)
→ Click "Review Order"
```

**Bước 3 - Review:**
```
Kiểm tra thông tin
→ Click "Place Order"
```

### 6. Success! ✅
```
Order ID: #1
Click: "View My Orders" để xem lịch sử
```

---

## 📋 CHECKLIST NHANH

- [ ] Login thành công
- [ ] Thêm 2 products vào cart
- [ ] Cart badge cập nhật
- [ ] Vào cart page
- [ ] Checkout bước 1 (Shipping)
- [ ] Checkout bước 2 (Payment)
- [ ] Checkout bước 3 (Review)
- [ ] Place order thành công
- [ ] Xem order trong Account/Orders

---

## 🎯 URLs QUAN TRỌNG

| Page | URL |
|------|-----|
| Login | http://localhost:3000/login |
| Products | http://localhost:3000/products |
| Cart | http://localhost:3000/cart |
| Checkout | http://localhost:3000/checkout |
| Account | http://localhost:3000/account |
| Orders | http://localhost:3000/account/orders |

---

## 🔧 NẾU GẶP LỖI

### Backend không chạy?
```bash
cd D:\TMDT
./server.exe
```

### Frontend không chạy?
```bash
cd D:\TMDT\frontend
npm run dev
```

### Cart không hoạt động?
- Phải login trước khi add to cart
- Refresh page và thử lại

### Order không tạo được?
- Check backend: `curl http://localhost:8080/health`
- Check console log trên browser (F12)

---

## 📞 TEST ACCOUNTS

```
Customer:
Email: customer@example.com
Password: customer123

Seller:
Email: seller@example.com
Password: seller123

Admin:
Email: admin@example.com
Password: admin123
```

---

## 🎉 EXPECTED FLOW

```
Login → Products → Add to Cart → Cart → Checkout (3 steps) → Success → View Orders
```

**Total time:** 3-5 minutes
**Products to add:** 2+ items (total > $50 for free shipping)
**Payment method:** COD (Cash on Delivery)

---

**Ready to test! 🚀**

**Chi tiết đầy đủ:** Xem `FULL_PURCHASE_FLOW_TEST.md`
