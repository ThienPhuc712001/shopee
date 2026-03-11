//go:build cache
// +build cache

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds cache configuration
type Config struct {
	// DefaultTTL is the default time-to-live for cache entries
	DefaultTTL time.Duration

	// ProductListTTL is TTL for product list caches
	ProductListTTL time.Duration

	// ProductDetailTTL is TTL for product detail caches
	ProductDetailTTL time.Duration

	// CategoryTTL is TTL for category caches
	CategoryTTL time.Duration

	// KeyPrefix is the prefix for all cache keys
	KeyPrefix string
}

// DefaultConfig returns default cache configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultTTL:       10 * time.Minute,
		ProductListTTL:   10 * time.Minute,
		ProductDetailTTL: 30 * time.Minute,
		CategoryTTL:      1 * time.Hour,
		KeyPrefix:        "ecommerce:",
	}
}

// Cache provides caching functionality for the e-commerce platform
type Cache struct {
	redis  *redis.Client
	config *Config
}

// NewCache creates a new cache instance
func NewCache(redisClient *redis.Client, config *Config) *Cache {
	if config == nil {
		config = DefaultConfig()
	}
	return &Cache{
		redis:  redisClient,
		config: config,
	}
}

// ==================== PRODUCT LIST CACHING ====================

// ProductListCacheKey generates cache key for product list
func (c *Cache) ProductListCacheKey(page, limit int, categoryID *uint, filters string) string {
	if categoryID != nil {
		return fmt.Sprintf("%sproducts:cat%d:p%d:l%d:f%s", c.config.KeyPrefix, *categoryID, page, limit, filters)
	}
	return fmt.Sprintf("%sproducts:all:p%d:l%d:f%s", c.config.KeyPrefix, page, limit, filters)
}

// GetProductList retrieves product list from cache
func (c *Cache) GetProductList(ctx context.Context, key string) ([]interface{}, int64, bool) {
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, 0, false
	}

	var cached struct {
		Products []interface{} `json:"products"`
		Total    int64         `json:"total"`
	}

	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, 0, false
	}

	return cached.Products, cached.Total, true
}

// SetProductList stores product list in cache
func (c *Cache) SetProductList(ctx context.Context, key string, products interface{}, total int64) error {
	data, err := json.Marshal(map[string]interface{}{
		"products": products,
		"total":    total,
	})
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, c.config.ProductListTTL).Err()
}

// InvalidateProductList invalidates product list cache
func (c *Cache) InvalidateProductList(ctx context.Context, categoryID *uint) error {
	var pattern string
	if categoryID != nil {
		pattern = fmt.Sprintf("%sproducts:cat%d:*", c.config.KeyPrefix, *categoryID)
	} else {
		pattern = fmt.Sprintf("%sproducts:*", c.config.KeyPrefix)
	}

	return c.deleteByPattern(ctx, pattern)
}

// ==================== PRODUCT DETAIL CACHING ====================

// ProductDetailCacheKey generates cache key for product detail
func (c *Cache) ProductDetailCacheKey(id uint) string {
	return fmt.Sprintf("%sproduct:detail:%d", c.config.KeyPrefix, id)
}

// GetProductDetail retrieves product detail from cache
func (c *Cache) GetProductDetail(ctx context.Context, id uint) (interface{}, bool) {
	key := c.ProductDetailCacheKey(id)
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var product interface{}
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, false
	}

	return product, true
}

// SetProductDetail stores product detail in cache
func (c *Cache) SetProductDetail(ctx context.Context, id uint, product interface{}) error {
	key := c.ProductDetailCacheKey(id)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, c.config.ProductDetailTTL).Err()
}

// DeleteProductDetail removes product detail from cache
func (c *Cache) DeleteProductDetail(ctx context.Context, id uint) error {
	key := c.ProductDetailCacheKey(id)
	return c.redis.Del(ctx, key).Err()
}

// ==================== CATEGORY CACHING ====================

// CategoryCacheKey generates cache key for category
func (c *Cache) CategoryCacheKey(id uint) string {
	return fmt.Sprintf("%scategory:%d", c.config.KeyPrefix, id)
}

// CategoryListCacheKey generates cache key for category list
func (c *Cache) CategoryListCacheKey(parentID *uint) string {
	if parentID != nil {
		return fmt.Sprintf("%scategories:parent%d", c.config.KeyPrefix, *parentID)
	}
	return fmt.Sprintf("%scategories:all", c.config.KeyPrefix)
}

// GetCategory retrieves category from cache
func (c *Cache) GetCategory(ctx context.Context, id uint) (interface{}, bool) {
	key := c.CategoryCacheKey(id)
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var category interface{}
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, false
	}

	return category, true
}

// SetCategory stores category in cache
func (c *Cache) SetCategory(ctx context.Context, id uint, category interface{}) error {
	key := c.CategoryCacheKey(id)
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, c.config.CategoryTTL).Err()
}

// GetCategoryList retrieves category list from cache
func (c *Cache) GetCategoryList(ctx context.Context, key string) (interface{}, bool) {
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var categories interface{}
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, false
	}

	return categories, true
}

// SetCategoryList stores category list in cache
func (c *Cache) SetCategoryList(ctx context.Context, key string, categories interface{}) error {
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, c.config.CategoryTTL).Err()
}

// InvalidateCategory invalidates category cache
func (c *Cache) InvalidateCategory(ctx context.Context) error {
	return c.deleteByPattern(ctx, fmt.Sprintf("%scategory*", c.config.KeyPrefix))
}

// ==================== HELPER METHODS ====================

// deleteByPattern deletes all keys matching a pattern
func (c *Cache) deleteByPattern(ctx context.Context, pattern string) error {
	iter := c.redis.Scan(ctx, 0, pattern, 100).Iterator()
	keysToDelete := make([]string, 0)

	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keysToDelete) > 0 {
		return c.redis.Del(ctx, keysToDelete...).Err()
	}

	return nil
}

// Clear clears all cache entries with the configured prefix
func (c *Cache) Clear(ctx context.Context) error {
	return c.deleteByPattern(ctx, c.config.KeyPrefix+"*")
}

// GetStats returns cache statistics
func (c *Cache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.redis.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	// Parse Redis info output
	stats := make(map[string]interface{})
	stats["info"] = info

	return stats, nil
}

// Ping checks if Redis connection is alive
func (c *Cache) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}
