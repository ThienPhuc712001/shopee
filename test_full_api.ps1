# Comprehensive Test - ALL Endpoints in Project
$BaseUrl = "http://localhost:8080/api"
$Results = @()

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  FULL API ENDPOINT TEST - ALL MODULES  " -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

function Test-Get {
    param([string]$Path, [string]$Name)
    try {
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method GET -UseBasicParsing -TimeoutSec 5
        Write-Host "  [PASS] $Name" -ForegroundColor Green
        $status = 'PASS'
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        Write-Host "  [FAIL] $Name - $code" -ForegroundColor $(if($code -eq 401){"Yellow"}else{"Red"})
        $status = if($code -eq 401){'AUTH'}else{'FAIL'}
    }
    $script:Results += @{Name=$Name; Path=$Path; Method='GET'; Status=$status; Code=$code}
}

function Test-Post {
    param([string]$Path, [string]$Name, [object]$Body)
    try {
        $json = $Body | ConvertTo-Json -Compress
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method POST -Body $json -ContentType "application/json" -UseBasicParsing -TimeoutSec 5
        Write-Host "  [PASS] $Name" -ForegroundColor Green
        $status = 'PASS'
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        Write-Host "  [FAIL] $Name - $code" -ForegroundColor $(if($code -eq 401){"Yellow"}else{"Red"})
        $status = if($code -eq 401){'AUTH'}else{'FAIL'}
    }
    $script:Results += @{Name=$Name; Path=$Path; Method='POST'; Status=$status; Code=$code}
}

function Test-Put {
    param([string]$Path, [string]$Name, [object]$Body)
    try {
        $json = if($Body){$Body | ConvertTo-Json -Compress}else{''}
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method PUT -Body $json -ContentType "application/json" -UseBasicParsing -TimeoutSec 5
        Write-Host "  [PASS] $Name" -ForegroundColor Green
        $status = 'PASS'
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        Write-Host "  [FAIL] $Name - $code" -ForegroundColor $(if($code -eq 401){"Yellow"}else{"Red"})
        $status = if($code -eq 401){'AUTH'}else{'FAIL'}
    }
    $script:Results += @{Name=$Name; Path=$Path; Method='PUT'; Status=$status; Code=$code}
}

function Test-Delete {
    param([string]$Path, [string]$Name)
    try {
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method DELETE -UseBasicParsing -TimeoutSec 5
        Write-Host "  [PASS] $Name" -ForegroundColor Green
        $status = 'PASS'
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        Write-Host "  [FAIL] $Name - $code" -ForegroundColor $(if($code -eq 401){"Yellow"}else{"Red"})
        $status = if($code -eq 401){'AUTH'}else{'FAIL'}
    }
    $script:Results += @{Name=$Name; Path=$Path; Method='DELETE'; Status=$status; Code=$code}
}

# ==================== PUBLIC ENDPOINTS ====================
Write-Host "`n[1] PUBLIC ENDPOINTS (No Auth)" -ForegroundColor Yellow
Write-Host "  --- Health ---" -ForegroundColor Gray
Test-Get "/health" "Health Check"
Test-Get "/health/upload" "Health Upload"

Write-Host "  --- Categories ---" -ForegroundColor Gray
Test-Get "/categories" "List Categories"
Test-Get "/categories/tree" "Category Tree"
Test-Get "/categories/featured" "Featured Categories"
Test-Get "/categories/1" "Get Category by ID"
Test-Get "/categories/1/products" "Category Products"
Test-Get "/categories/1/breadcrumb" "Category Breadcrumb"
Test-Get "/categories/search?q=elec" "Search Categories"

Write-Host "  --- Products ---" -ForegroundColor Gray
Test-Get "/products" "List Products"
Test-Get "/products/featured" "Featured Products"
Test-Get "/products/best-sellers" "Best Sellers"
Test-Get "/products/search?keyword=test" "Search Products"
Test-Get "/products/category/1" "Products by Category"
Test-Get "/products/1" "Get Product by ID"
Test-Get "/products/1/variants" "Product Variants"

Write-Host "  --- Shops ---" -ForegroundColor Gray
Test-Get "/shops" "List Shops"
Test-Get "/shops/1" "Get Shop by ID"

Write-Host "  --- Shipping ---" -ForegroundColor Gray
Test-Get "/shipping/carriers" "List Carriers"
Test-Get "/shipping/methods" "List Shipping Methods"
Test-Post "/shipping/calculate" "Calculate Shipping" @{
    from_city = "HCMC"; to_city = "HN"; weight = 1; shipping_method = "standard"
}

Write-Host "  --- Coupons (Public) ---" -ForegroundColor Gray
Test-Get "/coupons/active" "Active Coupons"
Test-Post "/coupons/apply" "Apply Coupon" @{code = "TEST"; order_total = 100}

Write-Host "  --- Inventory (Public) ---" -ForegroundColor Gray
Test-Post "/inventory/check" "Check Stock" @{product_id = 1; quantity = 1}

# ==================== AUTH ENDPOINTS ====================
Write-Host "`n[2] AUTH ENDPOINTS (Require Token)" -ForegroundColor Yellow

Write-Host "  --- Auth ---" -ForegroundColor Gray
Test-Post "/auth/register" "Register User" @{
    email = "test$(Get-Random)@example.com"; password = "Test1234!"; first_name = "Test"
}
Test-Post "/auth/login" "Login" @{email = "admin@example.com"; password = "admin123"}
Test-Get "/auth/me" "Current User (mock)"
Test-Post "/auth/refresh" "Refresh Token" @{refresh_token = "dummy"}
Test-Post "/auth/forgot-password" "Forgot Password" @{email = "test@example.com"}

# ==================== CART ENDPOINTS ====================
Write-Host "`n[3] CART ENDPOINTS" -ForegroundColor Yellow
Test-Get "/cart" "Get Cart"
Test-Get "/cart/summary" "Cart Summary"
Test-Get "/cart/stats" "Cart Stats"
Test-Get "/cart/checkout" "Prepare Checkout"
Test-Post "/cart/add" "Add to Cart" @{product_id = 1; quantity = 1}
Test-Post "/cart/items/1" "Update Cart Item" @{quantity = 2}
Test-Delete "/cart/items/1" "Remove from Cart"
Test-Delete "/cart/clear" "Clear Cart"

# ==================== ORDER ENDPOINTS ====================
Write-Host "`n[4] ORDER ENDPOINTS" -ForegroundColor Yellow
Test-Post "/orders/checkout" "Checkout Order" @{
    items = @(@{product_id = 1; quantity = 1})
    shipping_address = @{name = "Test"; phone = "0123"; address = "Test"; city = "HCMC"}
    payment_method = "cod"
}
Test-Get "/orders" "List Orders"
Test-Get "/orders/1" "Get Order"
Test-Get "/orders/1/tracking" "Order Tracking"
Test-Get "/orders/statistics" "Order Statistics"
Test-Post "/orders/1/cancel" "Cancel Order"

# ==================== PAYMENT ENDPOINTS ====================
Write-Host "`n[5] PAYMENT ENDPOINTS" -ForegroundColor Yellow
Test-Post "/webhook" "Payment Webhook" @{status = "success"}
Test-Post "/payments/create" "Create Payment" @{order_id = 1; method = "cod"}
Test-Post "/payments/confirm" "Confirm Payment" @{payment_id = "123"}
Test-Get "/payments/order/1" "Payment by Order"
Test-Get "/payments" "User Payments"
Test-Post "/payments/refund" "Request Refund" @{payment_id = "123"; reason = "test"}
Test-Post "/payments/methods" "Save Payment Method" @{type = "card"; token = "123"}
Test-Get "/payments/methods" "Get Payment Methods"
Test-Delete "/payments/methods/1" "Delete Payment Method"
Test-Post "/payments/methods/1/default" "Set Default Method"
Test-Get "/payments/statistics" "Payment Stats"

# ==================== NOTIFICATION ENDPOINTS ====================
Write-Host "`n[6] NOTIFICATION ENDPOINTS" -ForegroundColor Yellow
Test-Get "/notifications" "Get Notifications"
Test-Get "/notifications/summary" "Notification Summary"
Test-Get "/notifications/unread-count" "Unread Count"
Test-Get "/notifications/stats" "Notification Stats"
Test-Put "/notifications/1/read" "Mark as Read"
Test-Put "/notifications/read-all" "Mark All Read"
Test-Post "/notifications/mark-all-read" "Mark All Read (POST)"
Test-Delete "/notifications/1" "Delete Notification"
Test-Get "/notifications/preferences" "Get Preferences"
Test-Put "/notifications/preferences" "Update Preferences" @{email_enabled = $true}

# Admin notifications
Test-Post "/admin/notifications" "Create Notification" @{title = "Test"; message = "Test"}
Test-Post "/admin/notifications/batch" "Batch Notifications"
Test-Post "/admin/notifications/promotion" "Promotion Notification"
Test-Get "/admin/notifications/delivery-stats" "Delivery Stats"
Test-Post "/admin/notifications/cleanup" "Cleanup Notifications"

# ==================== COUPON ENDPOINTS ====================
Write-Host "`n[7] COUPON ENDPOINTS" -ForegroundColor Yellow
Test-Get "/coupons/my-usages" "My Coupon Usages"
Test-Get "/admin/coupons/stats" "Coupon Stats"
Test-Get "/admin/coupons" "List Coupons"
Test-Get "/admin/coupons/1" "Get Coupon"
Test-Post "/admin/coupons" "Create Coupon" @{code = "TEST123"; discount_type = "percent"; discount_value = 10}
Test-Put "/admin/coupons/1" "Update Coupon"
Test-Delete "/admin/coupons/1" "Delete Coupon"

# ==================== ADMIN ENDPOINTS ====================
Write-Host "`n[8] ADMIN ENDPOINTS" -ForegroundColor Yellow
Test-Post "/admin/auth/login" "Admin Login" @{email = "admin@example.com"; password = "admin123"}
Test-Get "/admin/users" "List Users"
Test-Post "/admin/users/ban" "Ban User" @{user_id = 1; reason = "test"}
Test-Get "/admin/sellers/pending" "Pending Sellers"
Test-Post "/admin/sellers/approve" "Approve Seller" @{seller_id = 1}
Test-Get "/admin/products" "Products for Moderation"
Test-Delete "/admin/products/1" "Delete Product"
Test-Get "/admin/orders" "All Orders"
Test-Post "/admin/orders/refund" "Refund Order" @{order_id = 1}
Test-Get "/admin/reviews" "All Reviews"
Test-Get "/admin/analytics/stats" "Admin Stats"
Test-Get "/admin/analytics/sales" "Sales Analytics"
Test-Get "/admin/analytics/users" "User Analytics"
Test-Get "/admin/analytics/products" "Product Analytics"
Test-Get "/admin/audit-logs" "Audit Logs"
Test-Get "/admin/settings/test" "Get Setting"
Test-Put "/admin/settings/test" "Update Setting" @{value = "test"}

# ==================== UPLOAD ENDPOINTS ====================
Write-Host "`n[9] UPLOAD ENDPOINTS" -ForegroundColor Yellow
Test-Post "/upload/product" "Upload Product Image"
Test-Post "/upload/product/multiple" "Upload Multiple Images"
Test-Delete "/upload/product/1" "Delete Product Image"
Test-Get "/upload/product/images" "Get Product Images"
Test-Post "/upload/review" "Upload Review Image"
Test-Post "/upload/avatar" "Upload Avatar"

# ==================== INVENTORY ENDPOINTS ====================
Write-Host "`n[10] INVENTORY ENDPOINTS" -ForegroundColor Yellow
Test-Get "/inventory" "Inventory List"
Test-Get "/inventory/summary" "Inventory Summary"
Test-Get "/inventory/low-stock" "Low Stock Products"
Test-Get "/inventory/out-of-stock" "Out of Stock"
Test-Post "/inventory/restock" "Restock Product" @{product_id = 1; quantity = 10}
Test-Post "/inventory/alerts" "Create Stock Alert" @{product_id = 1; threshold = 10; alert_type = "low_stock"}

# ==================== SUMMARY ====================
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "           TEST SUMMARY                 " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$total = $Results.Count
$pass = ($Results | Where-Object {$_.Status -eq 'PASS'}).Count
$auth = ($Results | Where-Object {$_.Status -eq 'AUTH'}).Count
$fail = ($Results | Where-Object {$_.Status -eq 'FAIL'}).Count

Write-Host "Total: $total" -ForegroundColor White
Write-Host "Pass: $pass" -ForegroundColor Green
Write-Host "Auth Required (401): $auth" -ForegroundColor Yellow
Write-Host "Fail: $fail" -ForegroundColor Red

if ($fail -gt 0) {
    Write-Host "`nFailed Endpoints:" -ForegroundColor Red
    $Results | Where-Object {$_.Status -eq 'FAIL'} | ForEach-Object {
        Write-Host "  - $($_.Name) [$($_.Method) $($_.Path)]" -ForegroundColor Red
    }
}

# Export to CSV
$csvPath = "D:\TMDT\full_test_results_$(Get-Date -Format 'yyyyMMdd_HHmmss').csv"
$Results | Export-Csv -Path $csvPath -NoTypeInformation
Write-Host "`nResults saved to: $csvPath" -ForegroundColor Gray
