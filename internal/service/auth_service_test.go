package service_test

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(limit, offset int) ([]model.User, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Search(keyword string, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(keyword, limit, offset)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(id uint, hashedPassword string) error {
	args := m.Called(id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateStatus(id uint, status model.UserStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// TestAuthService_Register tests user registration
func TestAuthService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		authService := service.NewAuthService(mockRepo, "test-secret", time.Hour)

		mockRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.Anything).Return(nil)

		// Act
		user, err := authService.Register("test@example.com", "password123", "John", "Doe")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, model.RoleCustomer, user.Role)

		mockRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		authService := service.NewAuthService(mockRepo, "test-secret", time.Hour)

		existingUser := &model.User{
			ID:    1,
			Email: "test@example.com",
		}
		mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

		// Act
		user, err := authService.Register("test@example.com", "password123", "John", "Doe")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, service.ErrUserAlreadyExists, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestAuthService_Login tests user login
func TestAuthService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		authService := service.NewAuthService(mockRepo, "test-secret", time.Hour)

		user := &model.User{
			ID:       1,
			Email:    "test@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password123"
			Role:     model.RoleCustomer,
			Status:   model.StatusActive,
		}
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		mockRepo.On("UpdateLastLogin", uint(1)).Return(nil)

		// Act
		token, returnedUser, err := authService.Login("test@example.com", "password123")

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, returnedUser)
		assert.Equal(t, user.ID, returnedUser.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		authService := service.NewAuthService(mockRepo, "test-secret", time.Hour)

		mockRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)

		// Act
		token, user, err := authService.Login("test@example.com", "wrongpassword")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.Equal(t, service.ErrInvalidCredentials, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_CreateProduct tests product creation
func TestProductService_CreateProduct(t *testing.T) {
	t.Run("successful product creation", func(t *testing.T) {
		// Arrange
		mockRepo := new(repository.MockProductRepository)
		productService := service.NewProductService(mockRepo)

		mockRepo.On("Create", mock.Anything).Return(nil)

		product := &model.Product{
			ShopID:   1,
			Name:     "Test Product",
			Price:    99.99,
			Stock:    100,
			Status:   model.ProductStatusActive,
		}

		// Act
		createdProduct, err := productService.CreateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)
		assert.Equal(t, "Test Product", createdProduct.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid product - missing name", func(t *testing.T) {
		// Arrange
		mockRepo := new(repository.MockProductRepository)
		productService := service.NewProductService(mockRepo)

		product := &model.Product{
			ShopID: 1,
			Price:  99.99,
		}

		// Act
		createdProduct, err := productService.CreateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdProduct)
		assert.Equal(t, service.ErrInvalidProduct, err)
	})
}

// TestCartService_AddItem tests adding items to cart
func TestCartService_AddItem(t *testing.T) {
	t.Run("successful add to cart", func(t *testing.T) {
		// Arrange
		mockCartRepo := new(repository.MockCartRepository)
		mockProductRepo := new(repository.MockProductRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo)

		product := &model.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
			Stock: 100,
		}

		cart := &model.Cart{
			ID:         1,
			UserID:     1,
			TotalItems: 0,
			TotalPrice: 0,
		}

		mockProductRepo.On("FindByID", uint(1)).Return(product, nil)
		mockCartRepo.On("FindOrCreate", uint(1)).Return(cart, nil)
		mockCartRepo.On("FindItemByCartAndProduct", uint(1), uint(1)).Return(nil, gorm.ErrRecordNotFound)
		mockCartRepo.On("AddItem", mock.Anything).Return(nil)
		mockCartRepo.On("UpdateTotals", uint(1)).Return(nil)

		// Act
		updatedCart, err := cartService.AddItem(1, 1, 2)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)

		mockProductRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})
}
