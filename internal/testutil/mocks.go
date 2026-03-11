package testutil

import (
	"context"
	"ecommerce/internal/domain/model"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ============================================================================
// MOCK USER REPOSITORY
// ============================================================================

// MockUserRepository is a mock implementation of UserRepositoryEnhanced
type MockUserRepository struct {
	Users         map[uint]*model.User
	UsersByEmail  map[string]*model.User
	UsersByPhone  map[string]*model.User
	ErrorToReturn error
	mu            sync.RWMutex
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:        make(map[uint]*model.User),
		UsersByEmail: make(map[string]*model.User),
		UsersByPhone: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) Create(user *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if user.ID == 0 {
		user.ID = uint(len(m.Users) + 1)
	}
	m.Users[user.ID] = user
	m.UsersByEmail[user.Email] = user
	if user.Phone != "" {
		m.UsersByPhone[user.Phone] = user
	}
	return nil
}

func (m *MockUserRepository) Update(user *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.Users[user.ID] = user
	m.UsersByEmail[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	delete(m.Users, id)
	return nil
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	user, exists := m.Users[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	user, exists := m.UsersByEmail[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindAll(limit, offset int) ([]model.User, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, 0, m.ErrorToReturn
	}

	users := make([]model.User, 0, len(m.Users))
	for _, user := range m.Users {
		users = append(users, *user)
	}

	total := int64(len(users))
	if offset >= len(users) {
		return []model.User{}, total, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], total, nil
}

func (m *MockUserRepository) Search(keyword string, limit, offset int) ([]model.User, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, 0, m.ErrorToReturn
	}

	var results []model.User
	for _, user := range m.Users {
		if contains(user.Email, keyword) || contains(user.FirstName, keyword) || contains(user.LastName, keyword) {
			results = append(results, *user)
		}
	}

	total := int64(len(results))
	if offset >= len(results) {
		return []model.User{}, total, nil
	}

	end := offset + limit
	if end > len(results) {
		end = len(results)
	}

	return results[offset:end], total, nil
}

func (m *MockUserRepository) UpdatePassword(id uint, hashedPassword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if user, exists := m.Users[id]; exists {
		user.Password = hashedPassword
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockUserRepository) UpdateLastLogin(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if user, exists := m.Users[id]; exists {
		now := time.Now()
		user.LastLogin = &now
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockUserRepository) UpdateStatus(id uint, status model.UserStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if user, exists := m.Users[id]; exists {
		user.Status = status
		return nil
	}
	return gorm.ErrRecordNotFound
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0)
}

// ============================================================================
// MOCK PRODUCT REPOSITORY
// ============================================================================

// MockProductRepository is a mock implementation of product repository
type MockProductRepository struct {
	Products      map[uint]*model.Product
	ProductsBySlug map[string]*model.Product
	ErrorToReturn error
	mu            sync.RWMutex
}

// NewMockProductRepository creates a new mock product repository
func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products:       make(map[uint]*model.Product),
		ProductsBySlug: make(map[string]*model.Product),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, product *model.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if product.ID == 0 {
		product.ID = uint(len(m.Products) + 1)
	}
	m.Products[product.ID] = product
	m.ProductsBySlug[product.Slug] = product
	return nil
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	product, exists := m.Products[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return product, nil
}

func (m *MockProductRepository) FindBySlug(ctx context.Context, slug string) (*model.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	product, exists := m.ProductsBySlug[slug]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return product, nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *model.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.Products[product.ID] = product
	m.ProductsBySlug[product.Slug] = product
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	delete(m.Products, id)
	return nil
}

// ============================================================================
// MOCK ORDER REPOSITORY
// ============================================================================

// MockOrderRepository is a mock implementation of order repository
type MockOrderRepository struct {
	Orders        map[uint]*model.Order
	OrdersByNumber map[string]*model.Order
	ErrorToReturn error
	mu            sync.RWMutex
}

// NewMockOrderRepository creates a new mock order repository
func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders:         make(map[uint]*model.Order),
		OrdersByNumber: make(map[string]*model.Order),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if order.ID == 0 {
		order.ID = uint(len(m.Orders) + 1)
	}
	m.Orders[order.ID] = order
	m.OrdersByNumber[order.OrderNumber] = order
	return nil
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id uint) (*model.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	order, exists := m.Orders[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return order, nil
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.Orders[order.ID] = order
	return nil
}

// ============================================================================
// MOCK CART REPOSITORY
// ============================================================================

// MockCartRepository is a mock implementation of cart repository
type MockCartRepository struct {
	Carts     map[uint]*model.Cart
	CartItems map[uint]*model.CartItem
	ErrorToReturn error
	mu        sync.RWMutex
}

// NewMockCartRepository creates a new mock cart repository
func NewMockCartRepository() *MockCartRepository {
	return &MockCartRepository{
		Carts:     make(map[uint]*model.Cart),
		CartItems: make(map[uint]*model.CartItem),
	}
}

func (m *MockCartRepository) Create(ctx context.Context, cart *model.Cart) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if cart.ID == 0 {
		cart.ID = uint(len(m.Carts) + 1)
	}
	m.Carts[cart.ID] = cart
	return nil
}

func (m *MockCartRepository) FindByUserID(ctx context.Context, userID uint) (*model.Cart, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	for _, cart := range m.Carts {
		if cart.UserID == userID {
			return cart, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// ============================================================================
// MOCK COUPON REPOSITORY
// ============================================================================

// MockCouponRepository is a mock implementation of coupon repository
type MockCouponRepository struct {
	Coupons     map[uint]*model.Coupon
	CouponsByCode map[string]*model.Coupon
	ErrorToReturn error
	mu          sync.RWMutex
}

// NewMockCouponRepository creates a new mock coupon repository
func NewMockCouponRepository() *MockCouponRepository {
	return &MockCouponRepository{
		Coupons:       make(map[uint]*model.Coupon),
		CouponsByCode: make(map[string]*model.Coupon),
	}
}

func (m *MockCouponRepository) FindByCode(ctx context.Context, code string) (*model.Coupon, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	coupon, exists := m.CouponsByCode[code]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return coupon, nil
}

// ============================================================================
// MOCK REFRESH TOKEN REPOSITORY
// ============================================================================

// MockRefreshTokenRepository is a mock implementation of refresh token repository
type MockRefreshTokenRepository struct {
	RefreshTokens map[string]*model.RefreshToken
	ErrorToReturn error
	mu            sync.RWMutex
}

// NewMockRefreshTokenRepository creates a new mock refresh token repository
func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{
		RefreshTokens: make(map[string]*model.RefreshToken),
	}
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.RefreshTokens[token.Token] = token
	return nil
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	rt, exists := m.RefreshTokens[token]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return rt, nil
}

func (m *MockRefreshTokenRepository) GetValidByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	rt, exists := m.RefreshTokens[token]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	if rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}
	return rt, nil
}

func (m *MockRefreshTokenRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	var tokens []*model.RefreshToken
	for _, rt := range m.RefreshTokens {
		if rt.UserID == userID {
			tokens = append(tokens, rt)
		}
	}
	return tokens, nil
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	if rt, exists := m.RefreshTokens[token]; exists {
		rt.Revoked = true
		now := time.Now()
		rt.RevokedAt = &now
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	for _, rt := range m.RefreshTokens {
		if rt.UserID == userID {
			rt.Revoked = true
			now := time.Now()
			rt.RevokedAt = &now
		}
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	delete(m.RefreshTokens, token)
	return nil
}

func (m *MockRefreshTokenRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorToReturn != nil {
		return 0, m.ErrorToReturn
	}

	var count int64
	for _, rt := range m.RefreshTokens {
		if rt.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *MockRefreshTokenRepository) CleanupOldTokens(ctx context.Context, retentionDays int) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return 0, m.ErrorToReturn
	}

	var count int64
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	for token, rt := range m.RefreshTokens {
		if (rt.Revoked && rt.RevokedAt != nil && rt.RevokedAt.Before(cutoffTime)) || rt.ExpiresAt.Before(time.Now()) {
			delete(m.RefreshTokens, token)
			count++
		}
	}
	return count, nil
}

func (m *MockRefreshTokenRepository) RevokeAllExceptCurrent(ctx context.Context, userID int64, currentToken string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	for _, rt := range m.RefreshTokens {
		if rt.UserID == userID && rt.Token != currentToken {
			rt.Revoked = true
			now := time.Now()
			rt.RevokedAt = &now
		}
	}
	return nil
}

// ============================================================================
// MOCK TOKEN SERVICE
// ============================================================================

// MockTokenService is a mock implementation of token service
type MockTokenService struct {
	GeneratedTokens *TokenPair
	ValidateError   error
	ExpiryTime      time.Time
}

// NewMockTokenService creates a new mock token service
func NewMockTokenService() *MockTokenService {
	return &MockTokenService{
		GeneratedTokens: &TokenPair{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
			TokenType:    "Bearer",
			ExpiresIn:    900,
			Expiry:       time.Now().Add(15 * time.Minute),
		},
		ExpiryTime: time.Now().Add(15 * time.Minute),
	}
}

func (m *MockTokenService) GenerateTokenPair(user *model.User) (*TokenPair, error) {
	if m.ValidateError != nil {
		return nil, m.ValidateError
	}
	return m.GeneratedTokens, nil
}

func (m *MockTokenService) GenerateAccessToken(user *model.User) (string, time.Time, error) {
	if m.ValidateError != nil {
		return "", time.Time{}, m.ValidateError
	}
	return m.GeneratedTokens.AccessToken, m.ExpiryTime, nil
}

func (m *MockTokenService) GenerateRefreshToken(user *model.User) (string, time.Time, error) {
	if m.ValidateError != nil {
		return "", time.Time{}, m.ValidateError
	}
	return m.GeneratedTokens.RefreshToken, m.ExpiryTime.Add(7 * 24 * time.Hour), nil
}

func (m *MockTokenService) ValidateToken(tokenString string, expectedType TokenType) (*TokenClaims, error) {
	if m.ValidateError != nil {
		return nil, m.ValidateError
	}
	return &TokenClaims{
		UserID:    1,
		Email:     "test@example.com",
		Role:      model.RoleCustomer,
		TokenType: expectedType,
	}, nil
}

func (m *MockTokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return m.ValidateToken(tokenString, TokenTypeAccess)
}

func (m *MockTokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return m.ValidateToken(tokenString, TokenTypeRefresh)
}

func (m *MockTokenService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	if m.ValidateError != nil {
		return nil, m.ValidateError
	}
	return m.GeneratedTokens, nil
}

func (m *MockTokenService) RevokeToken(tokenString string) error {
	return nil
}

func (m *MockTokenService) GetTokenExpiry(tokenString string) (time.Time, error) {
	return m.ExpiryTime, nil
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	Expiry       time.Time
}

// TokenType defines the type of token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID    uint
	Email     string
	Role      model.UserRole
	TokenType TokenType
}

// GetRefreshExpiry returns the refresh token expiry duration
func (m *MockTokenService) GetRefreshExpiry() time.Duration {
	return 7 * 24 * time.Hour
}

// ============================================================================
// ERROR HELPERS
// ============================================================================

// MockError is a generic mock error
var MockError = errors.New("mock error")

// ErrNotFound is a mock not found error
var ErrNotFound = gorm.ErrRecordNotFound
