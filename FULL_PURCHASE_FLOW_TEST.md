# 🛒 HƯỚNG DẪN TEST LUỒNG MUA HÀNG ĐẦY ĐỦ

## 📋 TỔNG QUAN LUỒNG MUA HÀNG

```
Homepage → Products → Product Detail → Cart → Checkout → Success
```

**Các bước:**
1. ✅ Xem sản phẩm
2. ✅ Thêm vào giỏ
3. ✅ Xem giỏ hàng
4. ✅ Thanh toán (3 bước)
   - Bước 1: Shipping Information
   - Bước 2: Payment Method
   - Bước 3: Review Order
5. ✅ Order Success

---

## 🎯 HƯỚNG DẪN TEST CHI TIẾT

### BƯỚC 1: LOGIN

**URL:** http://localhost:3000/login

**Credentials:**
```
Email: customer@example.com
Password: customer123
```

**Expected:**
- ✅ Login thành công
- ✅ Header hiển thị user name
- ✅ Cart badge hiển thị (nếu có)

---

### BƯỚC 2: XEM SẢN PHẨM

**URL:** http://localhost:3000/products

**Test:**
1. ✅ Scroll xem danh sách 8 sản phẩm
2. ✅ Click vào sản phẩm bất kỳ

**Expected:**
- ✅ Hiển thị 8 sản phẩm từ database
- ✅ Product cards có hover effect
- ✅ Click vào product chuyển sang detail page

---

### BƯỚC 3: XEM CHI TIẾT SẢN PHẨM

**URL:** http://localhost:3000/products/1

**Test:**
1. ✅ Xem thông tin sản phẩm
2. ✅ Xem giá, discount
3. ✅ Click "Add to Cart"

**Expected:**
- ✅ Hiển thị đầy đủ thông tin sản phẩm
- ✅ Giá hiển thị đúng (ví dụ: $199.99, was $299.99)
- ✅ Stock status: "In Stock" hoặc "Only X left"
- ✅ Click "Add to Cart" → Cart badge cập nhật số lượng

**Sản phẩm test:**
```
Premium Wireless Headphones
- Price: $199.99 (was $299.99)
- Stock: 50
- Rating: 4.8/5 (256 reviews)
```

---

### BƯỚC 4: THÊM NHIỀU SẢN PHẨM VÀO CART

**Lặp lại Bước 2-3 với các sản phẩm:**

1. **Smart Watch Pro** - $349.99
2. **Bluetooth Speaker** - $79.99
3. **USB-C Hub** - $49.99

**Expected:**
- ✅ Cart badge trên header cập nhật: 4 items
- ✅ Mỗi lần click "Add to Cart" có thông báo/animation

---

### BƯỚC 5: XEM GIỎ HÀNG

**URL:** http://localhost:3000/cart

**Test:**
1. ✅ Xem danh sách sản phẩm trong cart
2. ✅ Update quantity (tăng/giảm)
3. ✅ Xem subtotal cập nhật
4. ✅ Click "Proceed to Checkout"

**Expected:**
- ✅ Hiển thị đầy đủ 4 sản phẩm
- ✅ Quantity controls hoạt động (+/-)
- ✅ Subtotal tính đúng: $679.96
- ✅ Shipping: FREE (vì > $50)
- ✅ Total: $679.96

**Cart Summary:**
```
Subtotal: $679.96
Shipping: FREE
Total: $679.96
```

---

### BƯỚC 6: CHECKOUT - BƯỚC 1/3 (SHIPPING)

**URL:** http://localhost:3000/checkout

**Step Indicator:**
```
📦 Shipping → 💳 Payment → 🔒 Review → ✅ Complete
   (Current)
```

**Điền thông tin shipping:**

```
Full Name: John Doe
Phone: 0901234567
City: Ho Chi Minh City
Street Address: 123 Nguyen Hue, Ward 1
District: District 1
```

**Expected:**
- ✅ Form hiển thị đầy đủ fields
- ✅ Auto-fill nếu user đã có thông tin
- ✅ Validation: Required fields có mark *
- ✅ Click "Continue to Payment" → Sang bước 2

---

### BƯỚC 7: CHECKOUT - BƯỚC 2/3 (PAYMENT)

**Step Indicator:**
```
📦 Shipping → 💳 Payment → 🔒 Review → ✅ Complete
              (Current)
```

**Chọn payment method:**

**Options:**
1. ✅ **Cash on Delivery (COD)** - Pay when received
2. ✅ **Bank Transfer** - Transfer to account
3. ✅ **Credit/Debit Card** - Card payment
4. ✅ **PayPal** - PayPal account

**Test:**
1. ✅ Click chọn COD
2. ✅ Click "Review Order"

**Expected:**
- ✅ Payment methods hiển thị rõ ràng
- ✅ Selection có visual feedback
- ✅ Click "Review Order" → Sang bước 3

---

### BƯỚC 8: CHECKOUT - BƯỚC 3/3 (REVIEW)

**Step Indicator:**
```
📦 Shipping → 💳 Payment → 🔒 Review → ✅ Complete
                          (Current)
```

**Review thông tin:**

**1. Shipping Address:**
```
John Doe
123 Nguyen Hue, Ward 1
District 1, Ho Chi Minh City
Phone: 0901234567
```

**2. Payment Method:**
```
Cash on Delivery (COD)
```

**3. Order Items:**
```
1x Premium Wireless Headphones - $199.99
1x Smart Watch Pro - $349.99
1x Bluetooth Speaker - $79.99
1x USB-C Hub - $49.99
```

**4. Order Summary:**
```
Subtotal: $679.96
Shipping: FREE ✓
Discount: $0.00
Total: $679.96
```

**Test:**
1. ✅ Review tất cả thông tin
2. ✅ Click "Place Order - $679.96"

**Expected:**
- ✅ Thông tin chính xác
- ✅ Total tính đúng
- ✅ Click "Place Order" → API call → Success page

---

### BƯỚC 9: ORDER SUCCESS

**Expected:**
```
✅ Order Placed Successfully!

Thank you for your purchase.
Order ID: #1

What's Next?
1. ✓ Order confirmation email will be sent
2. ✓ We'll process your order within 1-2 business days
3. ✓ You'll receive tracking information via email
4. ✓ Estimated delivery: 5-7 business days

[View My Orders] [Continue Shopping]
```

**Test:**
1. ✅ Xem order confirmation
2. ✅ Click "View My Orders" → Account page
3. ✅ Xem order trong history

---

## 📊 CHECKLIST TEST

### Pre-Checkout:
- [ ] Login successful
- [ ] Browse products
- [ ] View product detail
- [ ] Add 4 products to cart
- [ ] Cart badge updates

### Cart:
- [ ] View cart page
- [ ] Update quantity (increase)
- [ ] Update quantity (decrease)
- [ ] Subtotal updates correctly
- [ ] Click "Proceed to Checkout"

### Checkout Step 1 (Shipping):
- [ ] Form displays correctly
- [ ] Fill in all required fields
- [ ] City dropdown works
- [ ] Click "Continue to Payment"

### Checkout Step 2 (Payment):
- [ ] All payment methods display
- [ ] Select COD
- [ ] Visual feedback on selection
- [ ] Click "Review Order"

### Checkout Step 3 (Review):
- [ ] Shipping address correct
- [ ] Payment method correct
- [ ] Order items correct
- [ ] Quantities correct
- [ ] Prices correct
- [ ] Subtotal correct
- [ ] Shipping fee correct (FREE > $50)
- [ ] Total correct
- [ ] Click "Place Order"

### Order Success:
- [ ] Success page displays
- [ ] Order ID shown
- [ ] Next steps shown
- [ ] "View My Orders" link works
- [ ] "Continue Shopping" link works

### Post-Order:
- [ ] Cart is emptied
- [ ] Cart badge shows 0
- [ ] Order appears in account/orders

---

## 🎯 TEST KỊCH BẢN MẪU

### Kịch bản 1: Standard Purchase

```
1. Login: customer@example.com / customer123
2. Browse products
3. Add: Headphones ($199.99) + Watch ($349.99)
4. Cart total: $549.98
5. Checkout with:
   - Address: 123 Test St, HCMC
   - Payment: COD
6. Place order
7. Success! Order #1 created
```

### Kịch bản 2: With Free Shipping

```
1. Login
2. Add products totaling > $50
3. Cart shows: "✓ You qualified for FREE shipping!"
4. Checkout shows: Shipping: FREE
5. Complete order
```

### Kịch bản 3: Update Cart During Checkout

```
1. Add 4 products to cart
2. Go to checkout step 2
3. Click "Back to Shipping"
4. Go back to cart
5. Remove 1 item
6. Continue checkout
7. Review shows updated items
```

---

## 🐛 CÁC LỖI THƯỜNG GẶP

### Lỗi 1: "Cart is empty"
**Nguyên nhân:** Chưa login hoặc cart đã bị clear
**Fix:** Login lại và add products vào cart

### Lỗi 2: "Login Required"
**Nguyên nhân:** Chưa login khi vào checkout
**Fix:** Login trước khi checkout

### Lỗi 3: Order không tạo được
**Nguyên nhân:** Backend API lỗi hoặc database connection
**Fix:** 
1. Check backend đang chạy: `curl http://localhost:8080/health`
2. Check console log trên browser
3. Verify database có products

### Lỗi 4: Subtotal không đúng
**Nguyên nhân:** Cart state không update
**Fix:** Refresh page và add lại products

---

## 📞 CREDENTIALS

| Role | Email | Password |
|------|-------|----------|
| Customer | customer@example.com | customer123 |

---

## 🎉 EXPECTED RESULTS

### Sau khi hoàn thành luồng:

✅ **Frontend:**
- Order success page hiển thị
- Cart được clear
- Order hiển thị trong Account/Orders

✅ **Backend:**
- Order được lưu vào database
- Cart items được clear
- User có thể xem order history

✅ **User Experience:**
- Luồng mượt mà từ A-Z
- Step indicator rõ ràng
- Form validation hoạt động
- Loading states hiển thị

---

**Happy Testing! 🛒**

**Last Updated:** March 13, 2026
**Status:** ✅ READY FOR FULL FLOW TEST
