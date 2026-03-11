# Shipping and Order Status System Documentation

## Overview

The Shipping and Order Status System manages order shipping information and tracks order status from creation to delivery. It provides complete visibility into the order fulfillment process for both administrators and customers.

---

## Table of Contents

1. [Shipping Business Flow](#shipping-business-flow)
2. [Order Status States](#order-status-states)
3. [Database Schema](#database-schema)
4. [Go Models](#go-models)
5. [API Endpoints](#api-endpoints)
6. [Order Integration](#order-integration)
7. [Tracking History](#tracking-history)
8. [Shipping Carriers](#shipping-carriers)
9. [Frontend Integration](#frontend-integration)
10. [Usage Examples](#usage-examples)

---

## Shipping Business Flow

### Complete Order Fulfillment Process

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   USER      │     │   SYSTEM    │     │   ADMIN     │
│   PLACES    │────▶│   CREATES   │────▶│   CONFIRMS  │
│   ORDER     │     │   ORDER     │     │   ORDER     │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   ORDER     │◀────│   ORDER     │◀────│  PACKAGE    │
│  DELIVERED  │     │   SHIPPED   │     │  ORDER      │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Step-by-Step Flow

| Step | Status | Actor | Description |
|------|--------|-------|-------------|
| 1 | `pending` | User | User places order, payment pending |
| 2 | `confirmed` | Admin | Admin confirms order and payment |
| 3 | `processing` | Admin | Order being prepared |
| 4 | `packed` | Admin | Order packed and ready |
| 5 | `shipped` | Admin | Order handed to carrier |
| 6 | `in_transit` | System | Order in transit |
| 7 | `out_for_delivery` | System | Out for final delivery |
| 8 | `delivered` | System | Customer receives order |
| 9 | `completed` | System | Order completed |

---

## Order Status States

### Status Definitions

```go
type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusPaid       OrderStatus = "paid"
    OrderStatusProcessing OrderStatus = "processing"
    OrderStatusShipped    OrderStatus = "shipped"
    OrderStatusDelivered  OrderStatus = "delivered"
    OrderStatusCancelled  OrderStatus = "cancelled"
    OrderStatusRefunded   OrderStatus = "refunded"
)
```

### Shipment Status Definitions

```go
type ShipmentStatus string

const (
    ShipmentStatusPending        ShipmentStatus = "pending"
    ShipmentStatusConfirmed      ShipmentStatus = "confirmed"
    ShipmentStatusProcessing     ShipmentStatus = "processing"
    ShipmentStatusPacked         ShipmentStatus = "packed"
    ShipmentStatusShipped        ShipmentStatus = "shipped"
    ShipmentStatusInTransit      ShipmentStatus = "in_transit"
    ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
    ShipmentStatusDelivered      ShipmentStatus = "delivered"
    ShipmentStatusFailed         ShipmentStatus = "failed"
    ShipmentStatusReturned       ShipmentStatus = "returned"
    ShipmentStatusCancelled      ShipmentStatus = "cancelled"
)
```

### Status Transitions

```
Valid Status Transitions:

pending ──→ confirmed ──→ processing ──→ packed ──→ shipped
   │            │             │            │          │
   │            │             │            │          ▼
   │            │             │            │      in_transit
   │            │             │            │          │
   │            │             │            │          ▼
   │            │             │            │      out_for_delivery
   │            │             │            │       │         │
   │            │             │            │       ▼         ▼
   │            │             │            │   delivered   failed
   │            │             │            │
   ▼            ▼             ▼            ▼
cancelled   cancelled     cancelled

```

### Transition Rules

| From Status | To Status | Description |
|-------------|-----------|-------------|
| `pending` | `confirmed` | Order confirmed by admin |
| `pending` | `cancelled` | User cancels order |
| `confirmed` | `processing` | Order preparation starts |
| `confirmed` | `cancelled` | Admin cancels order |
| `processing` | `packed` | Order packed and ready |
| `processing` | `cancelled` | Admin cancels order |
| `packed` | `shipped` | Carrier picks up order |
| `shipped` | `in_transit` | Order in transit |
| `in_transit` | `out_for_delivery` | Out for final delivery |
| `in_transit` | `failed` | Delivery failed |
| `out_for_delivery` | `delivered` | Successfully delivered |
| `out_for_delivery` | `failed` | Delivery failed |

---

## Database Schema

### Shipping Addresses Table

```sql
CREATE TABLE shipping_addresses (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    recipient_name NVARCHAR(200) NOT NULL,
    phone NVARCHAR(20) NOT NULL,
    address_line NVARCHAR(500) NOT NULL,
    ward NVARCHAR(200),
    district NVARCHAR(200) NOT NULL,
    city NVARCHAR(200) NOT NULL,
    postal_code NVARCHAR(20),
    country NVARCHAR(100) DEFAULT 'Vietnam',
    is_default BIT DEFAULT 0,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    notes NVARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    deleted_at DATETIME,

    CONSTRAINT FK_shipping_addresses_user FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Shipments Table

```sql
CREATE TABLE shipments (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE,
    carrier_id BIGINT,
    carrier_name NVARCHAR(200),
    carrier_type NVARCHAR(50) DEFAULT 'third_party',
    tracking_number NVARCHAR(100),
    status NVARCHAR(50) NOT NULL DEFAULT 'pending',
    
    -- Shipping details
    shipping_from NVARCHAR(500),
    shipping_to NVARCHAR(500),
    weight DECIMAL(10,2),
    dimensions NVARCHAR(50), -- LxWxH
    package_count INT DEFAULT 1,
    
    -- Timestamps
    shipped_at DATETIME,
    estimated_delivery DATETIME,
    delivered_at DATETIME,
    failed_at DATETIME,
    failure_reason NVARCHAR(500),
    
    -- Additional info
    shipping_fee DECIMAL(18,2) DEFAULT 0,
    insurance_amount DECIMAL(18,2) DEFAULT 0,
    notes NVARCHAR(500),
    
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    deleted_at DATETIME,

    CONSTRAINT FK_shipments_order FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT FK_shipments_carrier FOREIGN KEY (carrier_id) REFERENCES shipping_carriers(id)
);
```

### Shipment Tracking Table

```sql
CREATE TABLE shipment_tracking (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    status NVARCHAR(50) NOT NULL,
    location NVARCHAR(300),
    description NVARCHAR(1000) NOT NULL,
    occurred_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_shipment_tracking_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id)
);
```

### Shipping Carriers Table

```sql
CREATE TABLE shipping_carriers (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(200) NOT NULL UNIQUE,
    code NVARCHAR(50) NOT NULL UNIQUE,
    type NVARCHAR(50) NOT NULL,
    contact_name NVARCHAR(200),
    phone NVARCHAR(20),
    email NVARCHAR(255),
    website NVARCHAR(500),
    api_endpoint NVARCHAR(500),
    api_key NVARCHAR(500),
    is_active BIT DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE()
);
```

### Relationships

```
users (1) ──────< shipping_addresses (>
                     │
orders (1) ──────< shipments (>────── shipping_carriers (1)
                     │
                     └──────< shipment_tracking (>
```

---

## Go Models

### ShippingAddress

```go
type ShippingAddress struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    UserID        uint      `gorm:"not null;index" json:"user_id"`
    RecipientName string    `gorm:"varchar(200);not null" json:"recipient_name"`
    Phone         string    `gorm:"varchar(20);not null" json:"phone"`
    AddressLine   string    `gorm:"varchar(500);not null" json:"address_line"`
    Ward          string    `gorm:"varchar(200)" json:"ward"`
    District      string    `gorm:"varchar(200);not null" json:"district"`
    City          string    `gorm:"varchar(200);not null" json:"city"`
    PostalCode    string    `gorm:"varchar(20)" json:"postal_code"`
    Country       string    `gorm:"varchar(100);default:'Vietnam'" json:"country"`
    IsDefault     bool      `gorm:"default:false" json:"is_default"`
    Latitude      float64   `gorm:"decimal(10,8)" json:"latitude"`
    Longitude     float64   `gorm:"decimal(11,8)" json:"longitude"`
    Notes         string    `gorm:"varchar(500)" json:"notes"`
    
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### Shipment

```go
type Shipment struct {
    ID              uint           `gorm:"primaryKey" json:"id"`
    OrderID         uint           `gorm:"not null;uniqueIndex" json:"order_id"`
    CarrierID       *uint          `gorm:"index" json:"carrier_id"`
    CarrierName     string         `gorm:"varchar(200)" json:"carrier_name"`
    CarrierType     CarrierType    `gorm:"varchar(50);default:'third_party'" json:"carrier_type"`
    TrackingNumber  string         `gorm:"varchar(100);index" json:"tracking_number"`
    Status          ShipmentStatus `gorm:"varchar(50);not null;default:'pending'" json:"status"`
    
    ShippingFrom    string         `gorm:"varchar(500)" json:"shipping_from"`
    ShippingTo      string         `gorm:"varchar(500)" json:"shipping_to"`
    Weight          float64        `gorm:"decimal(10,2)" json:"weight"`
    Dimensions      string         `gorm:"varchar(50)" json:"dimensions"`
    PackageCount    int            `gorm:"default:1" json:"package_count"`
    
    ShippedAt       *time.Time     `gorm:"type:datetime" json:"shipped_at"`
    EstimatedDelivery *time.Time   `gorm:"type:datetime" json:"estimated_delivery"`
    DeliveredAt     *time.Time     `gorm:"type:datetime" json:"delivered_at"`
    
    ShippingFee     float64        `gorm:"decimal(18,2);default:0" json:"shipping_fee"`
    Notes           string         `gorm:"varchar(500)" json:"notes"`
    
    Order    *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
    Tracking []ShipmentTracking `gorm:"foreignKey:ShipmentID" json:"tracking,omitempty"`
}
```

### ShipmentTracking

```go
type ShipmentTracking struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    ShipmentID  uint           `gorm:"not null;index" json:"shipment_id"`
    Status      ShipmentStatus `gorm:"varchar(50);not null" json:"status"`
    Location    string         `gorm:"varchar(300)" json:"location"`
    Description string         `gorm:"varchar(1000);not null" json:"description"`
    OccurredAt  time.Time      `gorm:"not null" json:"occurred_at"`
    
    Shipment *Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
}
```

---

## API Endpoints

### User APIs

#### Get User's Shipping Addresses
```http
GET /api/shipping/addresses
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "addresses": [
      {
        "id": 1,
        "recipient_name": "John Doe",
        "phone": "123456789",
        "address_line": "123 Main St",
        "district": "District 1",
        "city": "Ho Chi Minh City",
        "is_default": true
      }
    ],
    "count": 1
  }
}
```

#### Create Shipping Address
```http
POST /api/shipping/addresses
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
  "recipient_name": "John Doe",
  "phone": "123456789",
  "address_line": "123 Main St",
  "ward": "Ward 1",
  "district": "District 1",
  "city": "Ho Chi Minh City",
  "postal_code": "70000",
  "country": "Vietnam",
  "is_default": true
}

Response (201):
{
  "success": true,
  "message": "Address created successfully",
  "data": { ...address object... }
}
```

#### Track Order
```http
GET /api/orders/:id/tracking
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "shipment_id": 1,
    "tracking_number": "VN123456789",
    "current_status": "in_transit",
    "estimated_delivery": "2024-03-15T10:00:00Z",
    "events": [
      {
        "id": 3,
        "status": "in_transit",
        "location": "Ho Chi Minh City Hub",
        "description": "Package arrived at sorting facility",
        "occurred_at": "2024-03-14T08:30:00Z"
      },
      {
        "id": 2,
        "status": "shipped",
        "location": "Hanoi Warehouse",
        "description": "Package has been shipped",
        "occurred_at": "2024-03-13T14:00:00Z"
      }
    ]
  }
}
```

---

### Admin APIs

#### Create Shipment
```http
POST /api/shipments
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "order_id": 100,
  "carrier_id": 1,
  "carrier_name": "Giao Hàng Nhanh",
  "carrier_type": "third_party",
  "tracking_number": "GHN123456789",
  "shipping_from": "Hanoi Warehouse",
  "shipping_to": "Ho Chi Minh City",
  "weight": 1.5,
  "dimensions": "30x20x10",
  "package_count": 1,
  "estimated_delivery": "2024-03-15T10:00:00Z",
  "shipping_fee": 30000,
  "notes": "Handle with care"
}

Response (201):
{
  "success": true,
  "message": "Shipment created successfully",
  "data": { ...shipment object... }
}
```

#### Update Shipment Status
```http
PUT /api/shipments/:id/status
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "status": "shipped",
  "notes": "Handed to carrier GHN"
}

Response (200):
{
  "success": true,
  "message": "Shipment status updated successfully",
  "data": { ...shipment object... }
}
```

#### Add Tracking Event
```http
POST /api/shipments/:id/tracking
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "status": "in_transit",
  "location": "Ho Chi Minh City Hub",
  "description": "Package arrived at sorting facility",
  "occurred_at": "2024-03-14T08:30:00Z"
}

Response (201):
{
  "success": true,
  "message": "Tracking event added successfully",
  "data": { ...tracking event object... }
}
```

#### Get Shipments by Status
```http
GET /api/shipments?status=pending&page=1&limit=20
Authorization: Bearer {token}
Role: admin

Response (200):
{
  "success": true,
  "data": {
    "shipments": [...],
    "pagination": {
      "total": 50,
      "page": 1,
      "limit": 20,
      "total_pages": 3
    }
  }
}
```

#### Get Shipment Statistics
```http
GET /api/shipments/stats
Authorization: Bearer {token}
Role: admin

Response (200):
{
  "success": true,
  "data": {
    "total_shipments": 500,
    "pending_shipments": 25,
    "shipped_shipments": 150,
    "delivered_shipments": 300,
    "failed_shipments": 10,
    "average_delivery_days": 3.5
  }
}
```

#### Get Active Carriers
```http
GET /api/shipping/carriers

Response (200):
{
  "success": true,
  "data": {
    "carriers": [
      {
        "id": 1,
        "name": "Giao Hàng Nhanh",
        "code": "GHN",
        "type": "third_party",
        "website": "https://ghn.vn",
        "is_active": true
      },
      {
        "id": 2,
        "name": "Internal Shipping",
        "code": "INTERNAL",
        "type": "internal",
        "is_active": true
      }
    ],
    "count": 2
  }
}
```

---

## Order Integration

### Integration Flow

```
┌─────────────────────────────────────────────────────────────┐
│              ORDER TO SHIPMENT INTEGRATION                   │
├─────────────────────────────────────────────────────────────┤
│  1. Order placed (status: pending)                          │
│  2. Payment confirmed (status: paid)                        │
│  3. Admin confirms order (status: processing)               │
│  4. Shipment created automatically                          │
│  5. Order packed (shipment status: packed)                  │
│  6. Carrier picks up (shipment status: shipped)             │
│  7. Tracking events added                                   │
│  8. Order delivered (shipment status: delivered)            │
│  9. Order marked complete                                   │
└─────────────────────────────────────────────────────────────┘
```

### Code Integration

```go
// In order_service_enhanced.go
func (s *orderServiceEnhanced) ConfirmOrder(orderID uint, userID uint) (*model.Order, error) {
    // Update order status
    order, err := s.UpdateOrderStatus(orderID, model.OrderStatusProcessing, userID, "")
    if err != nil {
        return nil, err
    }

    // Create shipment automatically
    shipmentInput := service.CreateShipmentInput{
        OrderID:        orderID,
        CarrierName:    "Default Carrier",
        CarrierType:    "internal",
        TrackingNumber: generateTrackingNumber(),
        Status:         model.ShipmentStatusPending,
    }
    
    // This would be called via shipping service
    _ = shipmentInput
    
    return order, nil
}
```

---

## Tracking History

### Tracking Events

Tracking history is stored in the `shipment_tracking` table with the following typical events:

| Event | Status | Description |
|-------|--------|-------------|
| Order received | `pending` | Shipment created, awaiting confirmation |
| Order confirmed | `confirmed` | Order confirmed by seller |
| Processing | `processing` | Order being prepared |
| Packed | `packed` | Package ready for shipment |
| Shipped | `shipped` | Handed to carrier |
| In transit | `in_transit` | Package in transit |
| Arrived at hub | `in_transit` | At sorting facility |
| Out for delivery | `out_for_delivery` | Final delivery attempt |
| Delivered | `delivered` | Successfully delivered |
| Failed | `failed` | Delivery failed |

### Timeline View

The API returns tracking events in chronological order:

```json
{
  "shipment_id": 1,
  "tracking_number": "VN123456789",
  "current_status": "in_transit",
  "estimated_delivery": "2024-03-15T10:00:00Z",
  "events": [
    {
      "status": "in_transit",
      "location": "Ho Chi Minh City Hub",
      "description": "Package arrived at sorting facility",
      "occurred_at": "2024-03-14T08:30:00Z"
    },
    {
      "status": "shipped",
      "location": "Hanoi Warehouse",
      "description": "Package has been shipped",
      "occurred_at": "2024-03-13T14:00:00Z"
    },
    {
      "status": "pending",
      "location": "Hanoi Warehouse",
      "description": "Shipment created, awaiting confirmation",
      "occurred_at": "2024-03-13T10:00:00Z"
    }
  ]
}
```

---

## Shipping Carriers

### Carrier Types

| Type | Description | Examples |
|------|-------------|----------|
| `internal` | Company's own delivery service | Internal fleet |
| `third_party` | External courier services | GHN, GHTK, Viettel Post |
| `local` | Local delivery partners | Local taxi services |

### Default Carriers (Seeded)

```sql
-- Internal Shipping
INSERT INTO shipping_carriers (name, code, type) 
VALUES ('Internal Shipping', 'INTERNAL', 'internal');

-- Third-party carriers
INSERT INTO shipping_carriers (name, code, type, website) 
VALUES 
  ('Giao Hàng Nhanh', 'GHN', 'third_party', 'https://ghn.vn'),
  ('Giao Hàng Tiết Kiệm', 'GHTK', 'third_party', 'https://ghtk.vn'),
  ('Viettel Post', 'VIETTEL', 'third_party', 'https://viettelpost.com.vn'),
  ('Vietnam Post', 'VNPOST', 'third_party', 'https://vnpost.vn');
```

### Carrier Integration

For third-party carriers with APIs:
- Store API endpoint URL
- Store encrypted API key
- Implement carrier-specific tracking sync

---

## Frontend Integration

### Order Tracking Page

```html
<!-- Example timeline view -->
<div class="tracking-timeline">
  <div class="tracking-step completed">
    <div class="status-icon">✓</div>
    <div class="status-info">
      <h4>Order Placed</h4>
      <p>March 13, 2024 at 10:00 AM</p>
    </div>
  </div>
  
  <div class="tracking-step completed">
    <div class="status-icon">✓</div>
    <div class="status-info">
      <h4>Shipped</h4>
      <p>March 13, 2024 at 2:00 PM</p>
      <p class="location">Hanoi Warehouse</p>
    </div>
  </div>
  
  <div class="tracking-step active">
    <div class="status-icon">🚚</div>
    <div class="status-info">
      <h4>In Transit</h4>
      <p>March 14, 2024 at 8:30 AM</p>
      <p class="location">Ho Chi Minh City Hub</p>
    </div>
  </div>
  
  <div class="tracking-step">
    <div class="status-icon">📦</div>
    <div class="status-info">
      <h4>Out for Delivery</h4>
      <p>Estimated: March 15, 2024</p>
    </div>
  </div>
  
  <div class="tracking-step">
    <div class="status-icon">🏠</div>
    <div class="status-info">
      <h4>Delivered</h4>
      <p>Estimated: March 15, 2024</p>
    </div>
  </div>
</div>
```

### Status Display Components

```javascript
// Status badge colors
const statusColors = {
  pending: 'gray',
  confirmed: 'blue',
  processing: 'yellow',
  packed: 'purple',
  shipped: 'cyan',
  in_transit: 'blue',
  out_for_delivery: 'orange',
  delivered: 'green',
  failed: 'red',
  cancelled: 'red'
};

// Format status for display
function formatStatus(status) {
  const labels = {
    pending: 'Pending',
    confirmed: 'Confirmed',
    processing: 'Processing',
    packed: 'Packed',
    shipped: 'Shipped',
    in_transit: 'In Transit',
    out_for_delivery: 'Out for Delivery',
    delivered: 'Delivered',
    failed: 'Delivery Failed',
    cancelled: 'Cancelled'
  };
  return labels[status] || status;
}
```

---

## Usage Examples

### Example 1: Create Shipment (Admin)

```bash
curl -X POST http://localhost:8080/api/shipments \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 100,
    "carrier_name": "Giao Hàng Nhanh",
    "carrier_type": "third_party",
    "tracking_number": "GHN123456789",
    "shipping_from": "Hanoi Warehouse",
    "shipping_to": "Ho Chi Minh City",
    "weight": 1.5,
    "dimensions": "30x20x10",
    "estimated_delivery": "2024-03-15T10:00:00Z",
    "shipping_fee": 30000
  }'
```

### Example 2: Update Shipment Status (Admin)

```bash
curl -X PUT http://localhost:8080/api/shipments/1/status \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_transit",
    "notes": "Package arrived at Ho Chi Minh Hub"
  }'
```

### Example 3: Add Tracking Event (Admin)

```bash
curl -X POST http://localhost:8080/api/shipments/1/tracking \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_transit",
    "location": "Ho Chi Minh City Hub",
    "description": "Package processed at sorting facility",
    "occurred_at": "2024-03-14T08:30:00Z"
  }'
```

### Example 4: Track Order (User)

```bash
curl -X GET http://localhost:8080/api/orders/100/tracking \
  -H "Authorization: Bearer {user_token}"
```

### Example 5: Save Shipping Address (User)

```bash
curl -X POST http://localhost:8080/api/shipping/addresses \
  -H "Authorization: Bearer {user_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "recipient_name": "John Doe",
    "phone": "123456789",
    "address_line": "123 Main Street",
    "ward": "Ward 1",
    "district": "District 1",
    "city": "Ho Chi Minh City",
    "postal_code": "70000",
    "country": "Vietnam",
    "is_default": true
  }'
```

---

## Implementation Checklist

### Completed ✅

- [x] ShippingAddress model with all fields
- [x] Shipment model with status tracking
- [x] ShipmentTracking model for history
- [x] ShippingCarrier model for carriers
- [x] Repository layer with CRUD operations
- [x] Service layer with business logic
- [x] Handler layer with HTTP endpoints
- [x] API routes (user and admin)
- [x] Database migration script
- [x] Status transition validation
- [x] Tracking timeline view
- [x] Carrier management

### Future Enhancements 📋

- [ ] Automatic shipment creation on order confirmation
- [ ] Third-party carrier API integration
- [ ] Real-time tracking sync
- [ ] Delivery route optimization
- [ ] SMS/Email notifications
- [ ] Proof of delivery (signature/photo)
- [ ] Return shipment management
- [ ] Multi-package shipments

---

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Shipment not found" | No shipment for order | Create shipment first |
| "Invalid status transition" | Wrong status order | Follow valid transitions |
| "Shipment already exists" | Duplicate creation | Use update instead |
| No tracking events | Events not added | Add tracking events |

### Debug Queries

```sql
-- Check shipment status
SELECT s.tracking_number, s.status, s.created_at, o.order_number
FROM shipments s
JOIN orders o ON s.order_id = o.id
WHERE o.order_number = 'ORD-20240313-00001';

-- Get tracking history
SELECT st.status, st.location, st.description, st.occurred_at
FROM shipment_tracking st
WHERE st.shipment_id = 1
ORDER BY st.occurred_at DESC;

-- Check pending shipments
SELECT * FROM shipments 
WHERE status = 'pending' 
ORDER BY created_at DESC;
```

---

## License

Internal use only - E-Commerce Platform
