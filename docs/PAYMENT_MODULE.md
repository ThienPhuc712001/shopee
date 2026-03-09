# Payment System - Business Flow & Implementation

## PART 1 — Payment Business Flow

### Complete Payment Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAYMENT PROCESSING FLOW                       │
└─────────────────────────────────────────────────────────────────┘

1. CUSTOMER CREATES ORDER
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│  Order   │────────>│ Database │
   │          │  POST   │ Service  │  CREATE │          │
   │          │ /orders │          │  Order  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Order created with status = "pending"
   ✓ Payment status = "pending"
   ✓ Inventory locked
   ✓ Payment timeout set (15-30 min)

2. ORDER STATUS = PENDING
   ┌──────────┐         ┌──────────┐
   │  Order   │────────>│ Payment  │
   │  Record  │         │ Service  │
   └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Order awaiting payment
   ✓ Countdown timer started
   ✓ Auto-cancel if not paid

3. SYSTEM GENERATES PAYMENT REQUEST
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Payment  │────────>│ Payment  │────────>│ Database │
   │ Service  │         │ Gateway  │  CREATE │          │
   └──────────┘         └──────────┘  Payment │          │
                              │
                              ▼
   ✓ Payment record created
   ✓ Payment method selected
   ✓ Amount calculated
   ✓ Transaction ID generated

4. CUSTOMER CHOOSES PAYMENT METHOD
   ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Payment  │
   │          │  Select │ Methods  │
   └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Payment methods displayed:
     - Credit/Debit Card
     - Bank Transfer
     - E-Wallet (MoMo, ZaloPay)
     - Cash on Delivery (COD)
   ✓ Customer selects preferred method

5. CUSTOMER COMPLETES PAYMENT
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Payment  │────────>│ Payment  │
   │          │  Enter  │ Gateway  │  Process│ Gateway  │
   │          │  Details│ (Stripe) │         │ Server   │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Card details entered
   ✓ 3D Secure verification (if applicable)
   ✓ Payment authorized
   ✓ Transaction ID returned

6. PAYMENT GATEWAY SENDS CONFIRMATION
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Payment  │────────>│ Webhook  │────────>│ Payment  │
   │ Gateway  │  POST   │ Endpoint │  Verify │ Service  │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Webhook received
   ✓ Signature verified
   ✓ Transaction validated
   ✓ Duplicate check performed

7. SYSTEM UPDATES PAYMENT STATUS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Payment  │────────>│ Database │────────>│ Payment  │
   │ Service  │  UPDATE │          │  UPDATE │  Record  │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Payment status = "paid"
   ✓ Transaction recorded
   ✓ Payment time stamped

8. SYSTEM UPDATES ORDER STATUS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Payment  │────────>│  Order   │────────>│ Database │
   │ Service  │  UPDATE │ Service  │  UPDATE │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Order status = "paid"
   ✓ Seller notified
   ✓ Order processing begins
   ✓ Inventory locked permanently
```

---

## PART 2 — Payment Methods

### Supported Payment Methods

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAYMENT METHODS                               │
└─────────────────────────────────────────────────────────────────┘

1. CREDIT/DEBIT CARD
   ┌─────────────────────────────────────────────────────────┐
   │ Providers: Stripe, PayPal, VNPay                        │
   │ Processing: Instant                                     │
   │ Fee: 2-3% per transaction                               │
   │ Security: 3D Secure, PCI DSS compliant                  │
   │ Flow:                                                   │
   │   Customer enters card → Gateway processes → Confirm    │
   └─────────────────────────────────────────────────────────┘

2. BANK TRANSFER
   ┌─────────────────────────────────────────────────────────┐
   │ Providers: Local banks, SWIFT                           │
   │ Processing: 1-3 business days                           │
   │ Fee: Fixed or percentage                                │
   │ Flow:                                                   │
   │   Customer transfers → Bank confirms → System updates   │
   └─────────────────────────────────────────────────────────┘

3. E-WALLET
   ┌─────────────────────────────────────────────────────────┐
   │ Providers: MoMo, ZaloPay, VNPay eWallet                 │
   │ Processing: Instant                                     │
   │ Fee: 1-2% per transaction                               │
   │ Flow:                                                   │
   │   Customer redirects to eWallet → Confirms → Callback   │
   └─────────────────────────────────────────────────────────┘

4. CASH ON DELIVERY (COD)
   ┌─────────────────────────────────────────────────────────┐
   │ Processing: On delivery                                 │
   │ Fee: None (included in shipping)                        │
   │ Risk: Order may be refused                              │
   │ Flow:                                                   │
   │   Order placed → Delivered → Customer pays cash         │
   └─────────────────────────────────────────────────────────┘
```

---

## PART 3 — Database Tables

### Payments Table

```sql
CREATE TABLE [dbo].[Payments] (
    [id]                BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [order_id]          BIGINT         NOT NULL,
    [user_id]          BIGINT         NOT NULL,
    [transaction_id]    NVARCHAR(255)  UNIQUE,
    [payment_method]    NVARCHAR(50)   NOT NULL,
    [payment_provider]  NVARCHAR(50),
    [amount]            DECIMAL(18,2)  NOT NULL,
    [currency]          NVARCHAR(3)    DEFAULT 'USD',
    [status]            NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [gateway_response]  NVARCHAR(MAX), -- JSON
    [metadata]          NVARCHAR(MAX), -- JSON
    [paid_at]           DATETIME,
    [failed_at]         DATETIME,
    [failure_reason]    NVARCHAR(500),
    [refunded_amount]   DECIMAL(18,2)  DEFAULT 0,
    [created_at]        DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]        DATETIME       NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id])
);

CREATE INDEX [IX_Payments_OrderID] ON [Payments]([order_id]);
CREATE INDEX [IX_Payments_UserID] ON [Payments]([user_id]);
CREATE INDEX [IX_Payments_TransactionID] ON [Payments]([transaction_id]);
CREATE INDEX [IX_Payments_Status] ON [Payments]([status]);
```

### PaymentMethods Table

```sql
CREATE TABLE [dbo].[PaymentMethods] (
    [id]           BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]      BIGINT,
    [type]         NVARCHAR(50)  NOT NULL,
    [provider]     NVARCHAR(50)  NOT NULL,
    [name]         NVARCHAR(100),
    [last_four]    NVARCHAR(4),
    [expiry_month] INT,
    [expiry_year]  INT,
    [is_default]   BIT           NOT NULL DEFAULT 0,
    [token]        NVARCHAR(500), -- Payment provider token
    [metadata]     NVARCHAR(MAX),
    [created_at]   DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE
);

CREATE INDEX [IX_PaymentMethods_UserID] ON [PaymentMethods]([user_id]);
```

### PaymentTransactions Table

```sql
CREATE TABLE [dbo].[PaymentTransactions] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [payment_id]      BIGINT        NOT NULL,
    [type]            NVARCHAR(50)  NOT NULL, -- charge, refund, chargeback
    [amount]          DECIMAL(18,2) NOT NULL,
    [status]          NVARCHAR(20)  NOT NULL,
    [gateway_id]      NVARCHAR(255),
    [gateway_response] NVARCHAR(MAX),
    [processed_at]    DATETIME,
    [created_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([payment_id]) REFERENCES [Payments]([id])
);

CREATE INDEX [IX_PaymentTransactions_PaymentID] ON [PaymentTransactions]([payment_id]);
```

### Refunds Table

```sql
CREATE TABLE [dbo].[Refunds] (
    [id]                BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [payment_id]        BIGINT         NOT NULL,
    [order_id]          BIGINT         NOT NULL,
    [refund_number]     NVARCHAR(50)   UNIQUE,
    [amount]            DECIMAL(18,2)  NOT NULL,
    [reason]            NVARCHAR(500)  NOT NULL,
    [status]            NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [type]              NVARCHAR(20)   NOT NULL, -- full, partial
    [requested_by]      BIGINT,
    [approved_by]       BIGINT,
    [approved_at]       DATETIME,
    [processed_at]      DATETIME,
    [gateway_refund_id] NVARCHAR(255),
    [notes]             NVARCHAR(500),
    [created_at]        DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]        DATETIME       NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([payment_id]) REFERENCES [Payments]([id]),
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    FOREIGN KEY ([requested_by]) REFERENCES [Users]([id]),
    FOREIGN KEY ([approved_by]) REFERENCES [Users]([id])
);

CREATE INDEX [IX_Refunds_PaymentID] ON [Refunds]([payment_id]);
CREATE INDEX [IX_Refunds_OrderID] ON [Refunds]([order_id]);
CREATE INDEX [IX_Refunds_Status] ON [Refunds]([status]);
```

---

## PART 4-12 — Implementation

The complete implementation follows in the code files below.

### Files Created:

1. **Models** - Payment, PaymentMethod, PaymentTransaction, Refund
2. **Repository** - Payment repository with all CRUD operations
3. **Service** - Payment service with gateway integration
4. **Handler** - REST API handlers and webhook endpoint
5. **Routes** - Payment route definitions

### Key Features:

✅ **Multiple Payment Methods** - Card, Bank Transfer, E-Wallet, COD
✅ **Payment Gateway Integration** - Stripe, PayPal, VNPay ready
✅ **Webhook Handling** - Secure webhook endpoint with signature verification
✅ **Transaction Tracking** - Complete transaction history
✅ **Refund Support** - Full and partial refunds
✅ **Security** - Signature verification, duplicate prevention
✅ **Status Management** - Complete payment lifecycle
✅ **Audit Trail** - All payment changes logged

---

## Payment Status State Machine

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAYMENT STATUS FLOW                           │
└─────────────────────────────────────────────────────────────────┘

                    ┌─────────────┐
                    │   PENDING   │
                    │  (created)  │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        Payment      Payment      Timeout/
        Success      Failure      Cancel
              │            │            │
              ▼            ▼            ▼
       ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
       │    PAID     │ │   FAILED    │ │ CANCELLED   │
       │ (confirmed) │ │  (rejected) │ │  (expired)  │
       └──────┬──────┘ └─────────────┘ └─────────────┘
              │
              │ Refund Request
              │
              ▼
       ┌─────────────┐
       │  REFUNDED   │
       │  (returned) │
       └─────────────┘
```

---

## Security Measures

### Webhook Security
```
1. Verify webhook signature
2. Validate gateway IP
3. Check transaction ID uniqueness
4. Verify amount matches order
5. Idempotency check
```

### Payment Security
```
1. PCI DSS compliance for cards
2. Token storage (not raw card data)
3. Encrypted transmission
4. Duplicate payment prevention
5. Fraud detection integration
```

---

This payment module is production-ready with full gateway integration and security measures.
