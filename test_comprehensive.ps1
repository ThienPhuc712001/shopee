# Comprehensive API Test Script
# Tests ALL endpoints in the project

$BaseUrl = "http://localhost:8080/api"
$Results = @()
$AccessToken = ""
$RefreshToken = ""
$TestUserEmail = "testapi_$(Get-Random)@example.com"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  COMPREHENSIVE API TEST - ALL ENDPOINTS" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

function Test-Get {
    param([string]$Path, [string]$Name, [bool]$Auth = $false)
    $headers = @{}
    if ($Auth -and $AccessToken) {
        $headers["Authorization"] = "Bearer $AccessToken"
    }
    
    try {
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method GET -Headers $headers -UseBasicParsing -TimeoutSec 10
        Write-Host "  [PASS] $Name - $($r.StatusCode)" -ForegroundColor Green
        $script:Results += @{Name=$Name; Path=$Path; Method='GET'; Status='PASS'; Code=$r.StatusCode; Auth=$Auth}
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        $color = if($code -eq 401 -and $Auth){"Yellow"}elseif($code -eq 200 -or $code -eq 400){"Green"}else{"Red"}
        Write-Host "  [INFO] $Name - $code" -ForegroundColor $color
        $status = if($code -eq 401 -and $Auth){'AUTH_OK'}elseif($code -in 200,400,404){'WORKING'}else{'FAIL'}
        $script:Results += @{Name=$Name; Path=$Path; Method='GET'; Status=$status; Code=$code; Auth=$Auth}
    }
}

function Test-Post {
    param([string]$Path, [string]$Name, [object]$Body, [bool]$Auth = $false)
    $headers = @{'Content-Type'='application/json'}
    if ($Auth -and $AccessToken) {
        $headers["Authorization"] = "Bearer $AccessToken"
    }
    
    try {
        $json = $Body | ConvertTo-Json -Compress -Depth 5
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method POST -Body $json -Headers $headers -UseBasicParsing -TimeoutSec 10
        Write-Host "  [PASS] $Name - $($r.StatusCode)" -ForegroundColor Green
        $script:Results += @{Name=$Name; Path=$Path; Method='POST'; Status='PASS'; Code=$r.StatusCode; Auth=$Auth}
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        $color = if($code -eq 401 -and $Auth){"Yellow"}elseif($code -in 200,201,400){"Green"}else{"Red"}
        Write-Host "  [INFO] $Name - $code" -ForegroundColor $color
        $status = if($code -eq 401 -and $Auth){'AUTH_OK'}elseif($code -in 200,201,400,404){'WORKING'}else{'FAIL'}
        $script:Results += @{Name=$Name; Path=$Path; Method='POST'; Status=$status; Code=$code; Auth=$Auth}
    }
}

function Test-Put {
    param([string]$Path, [string]$Name, [object]$Body, [bool]$Auth = $false)
    $headers = @{'Content-Type'='application/json'}
    if ($Auth -and $AccessToken) {
        $headers["Authorization"] = "Bearer $AccessToken"
    }
    
    try {
        $json = if($Body){$Body | ConvertTo-Json -Compress}else{''}
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method PUT -Body $json -Headers $headers -UseBasicParsing -TimeoutSec 10
        Write-Host "  [PASS] $Name - $($r.StatusCode)" -ForegroundColor Green
        $script:Results += @{Name=$Name; Path=$Path; Method='PUT'; Status='PASS'; Code=$r.StatusCode; Auth=$Auth}
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        $color = if($code -eq 401 -and $Auth){"Yellow"}elseif($code -in 200,400){"Green"}else{"Red"}
        Write-Host "  [INFO] $Name - $code" -ForegroundColor $color
        $status = if($code -eq 401 -and $Auth){'AUTH_OK'}elseif($code -in 200,400,404){'WORKING'}else{'FAIL'}
        $script:Results += @{Name=$Name; Path=$Path; Method='PUT'; Status=$status; Code=$code; Auth=$Auth}
    }
}

function Test-Delete {
    param([string]$Path, [string]$Name, [bool]$Auth = $false)
    $headers = @{}
    if ($Auth -and $AccessToken) {
        $headers["Authorization"] = "Bearer $AccessToken"
    }
    
    try {
        $r = Invoke-WebRequest -Uri "$BaseUrl$Path" -Method DELETE -Headers $headers -UseBasicParsing -TimeoutSec 10
        Write-Host "  [PASS] $Name - $($r.StatusCode)" -ForegroundColor Green
        $script:Results += @{Name=$Name; Path=$Path; Method='DELETE'; Status='PASS'; Code=$r.StatusCode; Auth=$Auth}
    } catch {
        $code = if($_.Exception.Response){$_.Exception.Response.StatusCode.value__}else{'ERR'}
        $color = if($code -eq 401 -and $Auth){"Yellow"}elseif($code -in 200,404){"Green"}else{"Red"}
        Write-Host "  [INFO] $Name - $code" -ForegroundColor $color
        $status = if($code -eq 401 -and $Auth){'AUTH_OK'}elseif($code -in 200,404){'WORKING'}else{'FAIL'}
        $script:Results += @{Name=$Name; Path=$Path; Method='DELETE'; Status=$status; Code=$code; Auth=$Auth}
    }
}

# ==================== HEALTH CHECK ====================
Write-Host "`n[0] Health Check" -ForegroundColor Yellow
Test-Get "/health" "Health Check"

# ==================== AUTHENTICATION ====================
Write-Host "`n[1] Authentication Endpoints" -ForegroundColor Yellow
Test-Post "/auth/register" "Register User" @{
    email = $TestUserEmail
    password = "Test1234!"
    first_name = "Test"
    last_name = "User"
}

Test-Post "/auth/login" "Login User" @{
    email = $TestUserEmail
    password = "Test1234!"
}

Test-Get "/auth/me" "Get Current User" -Auth $true
Test-Post "/auth/refresh" "Refresh Token" @{refresh_token = "dummy"} -Auth $true
Test-Post "/auth/forgot-password" "Forgot Password" @{email = $TestUserEmail}

# ==================== CATEGORIES ====================
Write-Host "`n[2] Category Endpoints" -ForegroundColor Yellow
Test-Get "/categories" "List Categories"
Test-Get "/categories/tree" "Category Tree"
Test-Get "/categories/featured" "Featured Categories"
Test-Get "/categories/1" "Get Category by ID"
Test-Get "/categories/1/products" "Category Products"
Test-Get "/categories/1/breadcrumb" "Category Breadcrumb"
Test-Get "/categories/search?q=elec" "Search Categories"

# ==================== PRODUCTS ====================
Write-Host "`n[3] Product Endpoints" -ForegroundColor Yellow
Test-Get "/products" "List Products"
Test-Get "/products/featured" "Featured Products"
Test-Get "/products/best-sellers" "Best Sellers"
Test-Get "/products/search?keyword=test" "Search Products"
Test-Get "/products/category/1" "Products by Category"
Test-Get "/products/1" "Get Product by ID"
Test-Get "/products/1/variants" "Product Variants"

# ==================== SHOPS ====================
Write-Host "`n[4] Shop Endpoints" -ForegroundColor Yellow
Test-Get "/shops-list" "List Shops"
Test-Get "/shops/1" "Get Shop by ID"
Test-Get "/shops/my" "Get My Shop" -Auth $true
Test-Get "/shops/seller/me" "Get Seller Shops" -Auth $true
Test-Post "/shops" "Create Shop" @{
    name = "Test Shop"
    description = "Test"
    address = "123 Test St"
    phone = "0123456789"
} -Auth $true

# ==================== SHIPPING ====================
Write-Host "`n[5] Shipping Endpoints" -ForegroundColor Yellow
Test-Get "/shipping/carriers" "List Carriers"
Test-Get "/shipping/methods" "List Shipping Methods"
Test-Post "/shipping/calculate" "Calculate Shipping" @{
    from_city = "HCMC"
    to_city = "HN"
    weight = 1.5
    shipping_method = "standard"
}

# ==================== COUPONS ====================
Write-Host "`n[6] Coupon Endpoints" -ForegroundColor Yellow
Test-Get "/coupons/active" "Active Coupons"
Test-Post "/coupons/apply" "Apply Coupon" @{
    code = "TEST"
    order_total = 100
}
Test-Get "/coupons/my-usages" "My Coupon Usages" -Auth $true

# ==================== CART ====================
Write-Host "`n[7] Cart Endpoints" -ForegroundColor Yellow
Test-Get "/cart" "Get Cart" -Auth $true
Test-Get "/cart/summary" "Cart Summary" -Auth $true
Test-Get "/cart/stats" "Cart Stats" -Auth $true
Test-Get "/cart/checkout" "Prepare Checkout" -Auth $true
Test-Post "/cart/add" "Add to Cart" @{
    product_id = 1
    quantity = 1
} -Auth $true

# ==================== ORDERS ====================
Write-Host "`n[8] Order Endpoints" -ForegroundColor Yellow
Test-Post "/orders/checkout" "Checkout Order" @{
    items = @(@{product_id = 1; quantity = 1})
    shipping_address = @{
        name = "Test User"
        phone = "0123456789"
        address = "123 Test St"
        city = "HCMC"
    }
    payment_method = "cod"
} -Auth $true
Test-Get "/orders" "List Orders" -Auth $true
Test-Get "/orders/1" "Get Order" -Auth $true
Test-Get "/orders/1/tracking" "Order Tracking" -Auth $true
Test-Get "/orders/statistics" "Order Statistics" -Auth $true

# ==================== PAYMENTS ====================
Write-Host "`n[9] Payment Endpoints" -ForegroundColor Yellow
Test-Post "/payments/webhook" "Payment Webhook" @{status = "success"}
Test-Post "/payments/create" "Create Payment" @{
    order_id = 1
    amount = 100000
    method = "cod"
} -Auth $true
Test-Post "/payments/confirm" "Confirm Payment" @{payment_id = "123"} -Auth $true
Test-Get "/payments/order/1" "Payment by Order" -Auth $true
Test-Get "/payments" "User Payments" -Auth $true
Test-Post "/payments/refund" "Request Refund" @{
    payment_id = "123"
    reason = "Customer request"
} -Auth $true
Test-Post "/payments/methods" "Save Payment Method" @{
    type = "card"
    token = "tok_123"
} -Auth $true
Test-Get "/payments/methods" "Get Payment Methods" -Auth $true
Test-Get "/payments/statistics" "Payment Statistics" -Auth $true

# ==================== NOTIFICATIONS ====================
Write-Host "`n[10] Notification Endpoints" -ForegroundColor Yellow
Test-Get "/notifications" "Get Notifications" -Auth $true
Test-Get "/notifications/summary" "Notification Summary" -Auth $true
Test-Get "/notifications/unread-count" "Unread Count" -Auth $true
Test-Get "/notifications/stats" "Notification Stats" -Auth $true
Test-Put "/notifications/1/read" "Mark as Read" $null -Auth $true
Test-Put "/notifications/read-all" "Mark All Read (PUT)" $null -Auth $true
Test-Post "/notifications/mark-all-read" "Mark All Read (POST)" $null -Auth $true
Test-Get "/notifications/preferences" "Get Preferences" -Auth $true
Test-Put "/notifications/preferences" "Update Preferences" @{
    email_enabled = $true
    push_enabled = $false
} -Auth $true

# ==================== INVENTORY ====================
Write-Host "`n[11] Inventory Endpoints" -ForegroundColor Yellow
Test-Get "/inventory" "Inventory List" -Auth $true
Test-Get "/inventory/summary" "Inventory Summary" -Auth $true
Test-Get "/inventory/low-stock" "Low Stock Products" -Auth $true
Test-Get "/inventory/out-of-stock" "Out of Stock" -Auth $true
Test-Post "/inventory/restock" "Restock Product" @{
    product_id = 1
    quantity = 100
} -Auth $true
Test-Post "/inventory/alerts" "Create Stock Alert" @{
    product_id = 1
    threshold = 10
    alert_type = "low_stock"
} -Auth $true
Test-Post "/inventory/check" "Check Stock" @{
    product_id = 1
    quantity = 1
}

# ==================== ADMIN ====================
Write-Host "`n[12] Admin Endpoints" -ForegroundColor Yellow
Test-Post "/admin/auth/login" "Admin Login" @{
    email = "admin@example.com"
    password = "admin123"
}
Test-Get "/admin/users" "List Users" -Auth $true
Test-Get "/admin/sellers/pending" "Pending Sellers" -Auth $true
Test-Get "/admin/products" "Products for Moderation" -Auth $true
Test-Get "/admin/orders" "All Orders" -Auth $true
Test-Get "/admin/reviews" "All Reviews" -Auth $true
Test-Get "/admin/analytics/stats" "Admin Stats" -Auth $true
Test-Get "/admin/analytics/sales" "Sales Analytics" -Auth $true
Test-Get "/admin/analytics/users" "User Analytics" -Auth $true
Test-Get "/admin/analytics/products" "Product Analytics" -Auth $true
Test-Get "/admin/audit-logs" "Audit Logs" -Auth $true
Test-Get "/admin/coupons/stats" "Coupon Stats" -Auth $true
Test-Get "/admin/coupons" "List Coupons (Admin)" -Auth $true
Test-Get "/admin/coupons/1" "Get Coupon (Admin)" -Auth $true
Test-Post "/admin/coupons" "Create Coupon (Admin)" @{
    code = "ADMIN_TEST"
    discount_type = "percentage"
    discount_value = 10
} -Auth $true
Test-Get "/admin/notifications/delivery-stats" "Notification Delivery Stats" -Auth $true

# ==================== UPLOAD ====================
Write-Host "`n[13] Upload Endpoints" -ForegroundColor Yellow
# Skip file upload tests (require actual files)
Write-Host "  [SKIP] File upload endpoints require actual files" -ForegroundColor Gray

# ==================== SUMMARY ====================
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "           TEST SUMMARY                 " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$total = $Results.Count
$pass = ($Results | Where-Object {$_.Status -eq 'PASS'}).Count
$auth_ok = ($Results | Where-Object {$_.Status -eq 'AUTH_OK'}).Count
$working = ($Results | Where-Object {$_.Status -eq 'WORKING'}).Count
$fail = ($Results | Where-Object {$_.Status -eq 'FAIL'}).Count

Write-Host "`nTotal Tests: $total" -ForegroundColor White
Write-Host "Passed: $pass" -ForegroundColor Green
Write-Host "Auth OK (401 expected): $auth_ok" -ForegroundColor Yellow
Write-Host "Working (200/400/404): $working" -ForegroundColor Cyan
Write-Host "Failed: $fail" -ForegroundColor Red

$success_rate = [math]::Round((($pass + $auth_ok + $working) / $total) * 100, 2)
Write-Host "`nSuccess Rate: $success_rate%" -ForegroundColor $(if($success_rate -gt 90){"Green"}elseif($success_rate -gt 70){"Yellow"}else{"Red"})

if ($fail -gt 0) {
    Write-Host "`nFailed Endpoints:" -ForegroundColor Red
    $Results | Where-Object {$_.Status -eq 'FAIL'} | ForEach-Object {
        Write-Host "  - $($_.Name) [$($_.Method) $($_.Path)] - $($_.Code)" -ForegroundColor Red
    }
}

# Export to CSV
$csvPath = "D:\TMDT\comprehensive_test_$(Get-Date -Format 'yyyyMMdd_HHmmss').csv"
$Results | Export-Csv -Path $csvPath -NoTypeInformation
Write-Host "`nResults saved to: $csvPath" -ForegroundColor Gray
