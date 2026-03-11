-- ============================================
-- E-COMMERCE DATABASE INDEXES FOR SQL SERVER
-- ============================================
-- Purpose: Optimize query performance for the e-commerce platform
-- Database: Microsoft SQL Server 2019+
-- Aligned with complete_schema.sql
-- ============================================

-- ============================================
-- FIRST: Drop full-text index (depends on unique index)
-- ============================================
BEGIN TRY
    IF EXISTS (SELECT * FROM sys.fulltext_indexes WHERE object_id = OBJECT_ID('products'))
    BEGIN
        DROP FULLTEXT INDEX ON products;
        PRINT 'Dropped full-text index on products';
    END
    
    IF EXISTS (SELECT * FROM sys.fulltext_catalogs WHERE name = 'ProductSearchCatalog')
    BEGIN
        DROP FULLTEXT CATALOG ProductSearchCatalog;
        PRINT 'Dropped full-text catalog ProductSearchCatalog';
    END
END TRY
BEGIN CATCH
    PRINT 'Note: Full-Text Search not available or already dropped';
END CATCH
GO

-- ============================================
-- Drop existing indexes if they exist
-- ============================================
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_name') DROP INDEX IX_products_name ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_category_id_status') DROP INDEX IX_products_category_id_status ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_shop_id') DROP INDEX IX_products_shop_id ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_status') DROP INDEX IX_products_status ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_is_featured') DROP INDEX IX_products_is_featured ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_sold_count') DROP INDEX IX_products_sold_count ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_price') DROP INDEX IX_products_price ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_category_price') DROP INDEX IX_products_category_price ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_slug') DROP INDEX IX_products_slug ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_created_at') DROP INDEX IX_products_created_at ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_rating_avg') DROP INDEX IX_products_rating_avg ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_flash_sale') DROP INDEX IX_products_flash_sale ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_list_covering') DROP INDEX IX_products_list_covering ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_id_unique') DROP INDEX IX_products_id_unique ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_deleted_at') DROP INDEX IX_products_deleted_at ON products;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_categories_parent_id') DROP INDEX IX_categories_parent_id ON categories;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_categories_slug') DROP INDEX IX_categories_slug ON categories;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_categories_is_active') DROP INDEX IX_categories_is_active ON categories;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_categories_level_sort') DROP INDEX IX_categories_level_sort ON categories;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_categories_deleted_at') DROP INDEX IX_categories_deleted_at ON categories;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_user_id') DROP INDEX IX_orders_user_id ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_created_at') DROP INDEX IX_orders_created_at ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_user_id_status') DROP INDEX IX_orders_user_id_status ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_status') DROP INDEX IX_orders_status ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_total_amount') DROP INDEX IX_orders_total_amount ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_created_at_status') DROP INDEX IX_orders_created_at_status ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_shop_id') DROP INDEX IX_orders_shop_id ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_order_number') DROP INDEX IX_orders_order_number ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_payment_status') DROP INDEX IX_orders_payment_status ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_parent_order_id') DROP INDEX IX_orders_parent_order_id ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_orders_deleted_at') DROP INDEX IX_orders_deleted_at ON orders;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_order_items_order_id') DROP INDEX IX_order_items_order_id ON order_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_order_items_product_id') DROP INDEX IX_order_items_product_id ON order_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_order_items_variant_id') DROP INDEX IX_order_items_variant_id ON order_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_order_items_shop_id') DROP INDEX IX_order_items_shop_id ON order_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_order_items_deleted_at') DROP INDEX IX_order_items_deleted_at ON order_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_cart_items_cart_id') DROP INDEX IX_cart_items_cart_id ON cart_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_cart_items_product_id') DROP INDEX IX_cart_items_product_id ON cart_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_cart_items_variant_id') DROP INDEX IX_cart_items_variant_id ON cart_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_cart_items_shop_id') DROP INDEX IX_cart_items_shop_id ON cart_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_cart_items_deleted_at') DROP INDEX IX_cart_items_deleted_at ON cart_items;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_product_id') DROP INDEX IX_reviews_product_id ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_user_id') DROP INDEX IX_reviews_user_id ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_rating') DROP INDEX IX_reviews_rating ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_created_at') DROP INDEX IX_reviews_created_at ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_shop_id') DROP INDEX IX_reviews_shop_id ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_order_id') DROP INDEX IX_reviews_order_id ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_reviews_deleted_at') DROP INDEX IX_reviews_deleted_at ON reviews;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_shops_status') DROP INDEX IX_shops_status ON shops;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_shops_rating') DROP INDEX IX_shops_rating ON shops;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_shops_slug') DROP INDEX IX_shops_slug ON shops;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_shops_user_id') DROP INDEX IX_shops_user_id ON shops;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_shops_deleted_at') DROP INDEX IX_shops_deleted_at ON shops;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_users_email') DROP INDEX IX_users_email ON users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_users_phone') DROP INDEX IX_users_phone ON users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_users_role') DROP INDEX IX_users_role ON users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_users_status') DROP INDEX IX_users_status ON users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_users_deleted_at') DROP INDEX IX_users_deleted_at ON users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_variants_sku') DROP INDEX IX_product_variants_sku ON product_variants;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_variants_product_id') DROP INDEX IX_product_variants_product_id ON product_variants;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_variants_deleted_at') DROP INDEX IX_product_variants_deleted_at ON product_variants;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_images_product_id') DROP INDEX IX_product_images_product_id ON product_images;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_images_primary') DROP INDEX IX_product_images_primary ON product_images;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_product_images_deleted_at') DROP INDEX IX_product_images_deleted_at ON product_images;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_carts_user_id') DROP INDEX IX_carts_user_id ON carts;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_carts_deleted_at') DROP INDEX IX_carts_deleted_at ON carts;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_addresses_user_id') DROP INDEX IX_addresses_user_id ON addresses;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_addresses_deleted_at') DROP INDEX IX_addresses_deleted_at ON addresses;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_payments_order_id') DROP INDEX IX_payments_order_id ON payments;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_payments_user_id') DROP INDEX IX_payments_user_id ON payments;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_payments_transaction_id') DROP INDEX IX_payments_transaction_id ON payments;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_payments_status') DROP INDEX IX_payments_status ON payments;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_payments_deleted_at') DROP INDEX IX_payments_deleted_at ON payments;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_admin_users_email') DROP INDEX IX_admin_users_email ON admin_users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_admin_users_status') DROP INDEX IX_admin_users_status ON admin_users;
IF EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_admin_users_deleted_at') DROP INDEX IX_admin_users_deleted_at ON admin_users;
GO

-- ============================================
-- 1. PRODUCT TABLE INDEXES
-- ============================================

CREATE INDEX IX_products_name ON products(name);
CREATE INDEX IX_products_category_id_status ON products(category_id, status);
CREATE INDEX IX_products_shop_id ON products(shop_id);
CREATE INDEX IX_products_status ON products(status);
CREATE INDEX IX_products_is_featured ON products(is_featured) WHERE status = 'active';
CREATE INDEX IX_products_sold_count ON products(sold_count DESC);
CREATE INDEX IX_products_price ON products(price);
CREATE INDEX IX_products_category_price ON products(category_id, price);
CREATE UNIQUE INDEX IX_products_slug ON products(slug);
CREATE INDEX IX_products_created_at ON products(created_at DESC);
CREATE INDEX IX_products_rating_avg ON products(rating_avg DESC);
CREATE INDEX IX_products_flash_sale ON products(is_flash_sale, status) WHERE is_flash_sale = 1 AND status = 'active';
CREATE INDEX IX_products_list_covering ON products(status, created_at DESC)
    INCLUDE (id, name, price, original_price, discount_percent, rating_avg, sold_count, category_id, shop_id)
    WHERE status = 'active';
CREATE INDEX IX_products_deleted_at ON products(deleted_at);

-- ============================================
-- 2. CATEGORY TABLE INDEXES
-- ============================================

CREATE INDEX IX_categories_parent_id ON categories(parent_id);
CREATE UNIQUE INDEX IX_categories_slug ON categories(slug);
CREATE INDEX IX_categories_is_active ON categories(is_active) WHERE is_active = 1;
CREATE INDEX IX_categories_level_sort ON categories(level, sort_order);
CREATE INDEX IX_categories_deleted_at ON categories(deleted_at);

-- ============================================
-- 3. ORDER TABLE INDEXES
-- ============================================

CREATE INDEX IX_orders_user_id ON orders(user_id);
CREATE INDEX IX_orders_created_at ON orders(created_at DESC);
CREATE INDEX IX_orders_user_id_status ON orders(user_id, status);
CREATE INDEX IX_orders_status ON orders(status);
CREATE INDEX IX_orders_total_amount ON orders(total_amount DESC);
CREATE INDEX IX_orders_created_at_status ON orders(created_at, status);
CREATE INDEX IX_orders_shop_id ON orders(shop_id);
CREATE UNIQUE INDEX IX_orders_order_number ON orders(order_number);
CREATE INDEX IX_orders_payment_status ON orders(payment_status);
CREATE INDEX IX_orders_parent_order_id ON orders(parent_order_id);
CREATE INDEX IX_orders_deleted_at ON orders(deleted_at);

-- ============================================
-- 4. ORDER ITEMS TABLE INDEXES
-- ============================================

CREATE INDEX IX_order_items_order_id ON order_items(order_id);
CREATE INDEX IX_order_items_product_id ON order_items(product_id);
CREATE INDEX IX_order_items_variant_id ON order_items(variant_id);
CREATE INDEX IX_order_items_shop_id ON order_items(shop_id);
CREATE INDEX IX_order_items_deleted_at ON order_items(deleted_at);

-- ============================================
-- 5. CART & CART ITEMS TABLE INDEXES
-- ============================================

CREATE INDEX IX_carts_user_id ON carts(user_id);
CREATE INDEX IX_carts_deleted_at ON carts(deleted_at);

CREATE INDEX IX_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IX_cart_items_product_id ON cart_items(product_id);
CREATE INDEX IX_cart_items_variant_id ON cart_items(variant_id);
CREATE INDEX IX_cart_items_shop_id ON cart_items(shop_id);
CREATE INDEX IX_cart_items_deleted_at ON cart_items(deleted_at);

-- ============================================
-- 6. REVIEW TABLE INDEXES
-- ============================================

CREATE INDEX IX_reviews_product_id ON reviews(product_id);
CREATE INDEX IX_reviews_user_id ON reviews(user_id);
CREATE INDEX IX_reviews_shop_id ON reviews(shop_id);
CREATE INDEX IX_reviews_order_id ON reviews(order_id);
CREATE INDEX IX_reviews_rating ON reviews(rating DESC);
CREATE INDEX IX_reviews_created_at ON reviews(created_at DESC);
CREATE INDEX IX_reviews_deleted_at ON reviews(deleted_at);

-- ============================================
-- 7. SHOP TABLE INDEXES
-- ============================================

CREATE INDEX IX_shops_status ON shops(status);
CREATE INDEX IX_shops_rating ON shops(rating DESC);
CREATE UNIQUE INDEX IX_shops_slug ON shops(slug);
CREATE INDEX IX_shops_user_id ON shops(user_id);
CREATE INDEX IX_shops_deleted_at ON shops(deleted_at);

-- ============================================
-- 8. USER TABLE INDEXES
-- ============================================

CREATE UNIQUE INDEX IX_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IX_users_phone ON users(phone) WHERE phone IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX IX_users_role ON users(role);
CREATE INDEX IX_users_status ON users(status);
CREATE INDEX IX_users_deleted_at ON users(deleted_at);

-- ============================================
-- 9. PRODUCT VARIANT INDEXES
-- ============================================

CREATE UNIQUE INDEX IX_product_variants_sku ON product_variants(sku) WHERE deleted_at IS NULL;
CREATE INDEX IX_product_variants_product_id ON product_variants(product_id) WHERE deleted_at IS NULL;
CREATE INDEX IX_product_variants_deleted_at ON product_variants(deleted_at);

-- ============================================
-- 10. PRODUCT IMAGE INDEXES
-- ============================================

CREATE INDEX IX_product_images_product_id ON product_images(product_id);
CREATE INDEX IX_product_images_primary ON product_images(product_id, is_primary) WHERE is_primary = 1;
CREATE INDEX IX_product_images_deleted_at ON product_images(deleted_at);

-- ============================================
-- 11. ADDRESSES TABLE INDEXES
-- ============================================

CREATE INDEX IX_addresses_user_id ON addresses(user_id);
CREATE INDEX IX_addresses_deleted_at ON addresses(deleted_at);

-- ============================================
-- 12. PAYMENTS TABLE INDEXES
-- ============================================

CREATE INDEX IX_payments_order_id ON payments(order_id);
CREATE INDEX IX_payments_user_id ON payments(user_id);
CREATE UNIQUE INDEX IX_payments_transaction_id ON payments(transaction_id);
CREATE INDEX IX_payments_status ON payments(status);
CREATE INDEX IX_payments_deleted_at ON payments(deleted_at);

-- ============================================
-- 13. ADMIN TABLES INDEXES
-- ============================================

CREATE INDEX IX_admin_users_email ON admin_users(email);
CREATE INDEX IX_admin_users_status ON admin_users(status);
CREATE INDEX IX_admin_users_deleted_at ON admin_users(deleted_at);

-- ============================================
-- 14. FULL-TEXT SEARCH INDEXES (OPTIONAL)
-- ============================================

GO

BEGIN TRY
    -- Create full-text catalog
    IF NOT EXISTS (SELECT * FROM sys.fulltext_catalogs WHERE name = 'ProductSearchCatalog')
    BEGIN
        CREATE FULLTEXT CATALOG ProductSearchCatalog AS DEFAULT;
    END

    -- Create unique index for full-text key
    IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'IX_products_id_unique' AND object_id = OBJECT_ID('products'))
    BEGIN
        CREATE UNIQUE INDEX IX_products_id_unique ON products(id);
    END

    -- Create full-text index
    IF NOT EXISTS (SELECT * FROM sys.fulltext_indexes WHERE object_id = OBJECT_ID('products'))
    BEGIN
        CREATE FULLTEXT INDEX ON products(
            name LANGUAGE 1033,
            description LANGUAGE 1033,
            short_description LANGUAGE 1033,
            brand LANGUAGE 1033
        )
        KEY INDEX IX_products_id_unique
        ON ProductSearchCatalog
        WITH CHANGE_TRACKING AUTO;
    END
END TRY
BEGIN CATCH
    PRINT 'WARNING: Full-Text Search is not available. Skipping full-text index creation.';
    PRINT 'Error details: ' + ERROR_MESSAGE();
END CATCH

-- ============================================
-- END OF INDEXES SCRIPT
-- ============================================
