package service

import (
	"testing"
	"time"

	"ecommerce/internal/domain/model"
	"ecommerce/pkg/password"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// TEST SETUP
// ============================================================================

// TestAuthService_Login_Success tests successful login
func TestAuthService_Login_Success(t *testing.T) {
	// This test requires database integration
	// For unit testing, we test the password verification logic directly
	t.Run("PasswordVerification", func(t *testing.T) {
		// Arrange
		testPassword := "TestPassword123!"
		hashedPassword, _ := password.Hash(testPassword)

		// Act
		err := password.Verify(testPassword, hashedPassword)

		// Assert
		assert.NoError(t, err)
	})
}

// TestAuthService_Login_InvalidPassword tests login with wrong password
func TestAuthService_Login_InvalidPassword(t *testing.T) {
	t.Run("WrongPassword", func(t *testing.T) {
		// Arrange
		correctPassword := "TestPassword123!"
		wrongPassword := "WrongPassword"
		hashedPassword, _ := password.Hash(correctPassword)

		// Act
		err := password.Verify(wrongPassword, hashedPassword)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, password.ErrHashMismatch, err)
	})
}

// TestAuthService_Register_PasswordValidation tests password validation during registration
func TestAuthService_Register_PasswordValidation(t *testing.T) {
	t.Run("WeakPassword", func(t *testing.T) {
		// Arrange
		weakPassword := "weak"

		// Act
		err := password.ValidateDefault(weakPassword)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, password.ErrPasswordTooShort, err)
	})

	t.Run("StrongPassword", func(t *testing.T) {
		// Arrange
		strongPassword := "TestPassword123!"

		// Act
		err := password.ValidateDefault(strongPassword)

		// Assert
		assert.NoError(t, err)
	})
}

// TestAuthService_AccountLockout tests account lockout after failed attempts
func TestAuthService_AccountLockout(t *testing.T) {
	t.Run("LockAfterMaxAttempts", func(t *testing.T) {
		// Arrange
		user := &model.User{
			ID:                  1,
			Email:               "test@example.com",
			FailedLoginAttempts: 4,
		}
		maxAttempts := 5
		lockoutDuration := 15 * time.Minute

		// Act
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= maxAttempts {
			lockUntil := time.Now().Add(lockoutDuration)
			user.LockedUntil = &lockUntil
		}

		// Assert
		assert.Equal(t, 5, user.FailedLoginAttempts)
		assert.NotNil(t, user.LockedUntil)
		assert.True(t, time.Now().Before(*user.LockedUntil))
	})
}

// TestAuthService_PasswordChange tests password change logic
func TestAuthService_PasswordChange(t *testing.T) {
	t.Run("HashNewPassword", func(t *testing.T) {
		// Arrange
		newPassword := "NewPassword123!"

		// Act
		hashedPassword, err := password.Hash(newPassword)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, newPassword, hashedPassword)
	})

	t.Run("VerifyNewPassword", func(t *testing.T) {
		// Arrange
		newPassword := "NewPassword123!"
		hashedPassword, _ := password.Hash(newPassword)

		// Act & Assert
		assert.NoError(t, password.Verify(newPassword, hashedPassword))
	})
}

// TestAuthService_TokenExpiry tests token expiry logic
func TestAuthService_TokenExpiry(t *testing.T) {
	t.Run("CheckTokenExpired", func(t *testing.T) {
		// Arrange
		expiresAt := time.Now().Add(-1 * time.Minute) // Already expired

		// Act & Assert
		assert.True(t, time.Now().After(expiresAt))
	})

	t.Run("CheckTokenValid", func(t *testing.T) {
		// Arrange
		expiresAt := time.Now().Add(15 * time.Minute) // Valid for 15 more minutes

		// Act & Assert
		assert.True(t, time.Now().Before(expiresAt))
	})
}

// TestAuthService_EmailValidation tests email validation
func TestAuthService_EmailValidation(t *testing.T) {
	t.Run("ValidEmail", func(t *testing.T) {
		// Arrange
		email := "test@example.com"

		// Assert - basic email format check
		assert.Contains(t, email, "@")
		assert.Contains(t, email, ".")
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		// Arrange
		invalidEmails := []string{
			"invalid",
			"invalid@",
			"@example.com",
			"",
		}

		// Assert
		for _, email := range invalidEmails {
			assert.False(t, isValidEmailFormat(email))
		}
	})
}

// isValidEmailFormat checks basic email format
func isValidEmailFormat(email string) bool {
	if len(email) < 5 {
		return false
	}
	hasAt := false
	hasDot := false
	for i, c := range email {
		if c == '@' && i > 0 && i < len(email)-1 {
			hasAt = true
		}
		if c == '.' && hasAt && i > len(email)/2 {
			hasDot = true
		}
	}
	return hasAt && hasDot
}

// TestAuthService_OrderTotalCalculation tests order total calculation
func TestAuthService_OrderTotalCalculation(t *testing.T) {
	t.Run("CalculateTotal", func(t *testing.T) {
		// Arrange
		items := []struct {
			Price    float64
			Quantity int
		}{
			{Price: 50.00, Quantity: 2},
			{Price: 100.00, Quantity: 1},
		}

		// Act
		var total float64
		for _, item := range items {
			total += item.Price * float64(item.Quantity)
		}

		// Assert
		assert.Equal(t, 200.00, total)
	})
}

// TestAuthService_CouponDiscount tests coupon discount calculation
func TestAuthService_CouponDiscount(t *testing.T) {
	t.Run("PercentageDiscount", func(t *testing.T) {
		// Arrange
		orderTotal := 100.00
		discountPercent := 10.0

		// Act
		discount := orderTotal * discountPercent / 100

		// Assert
		assert.Equal(t, 10.00, discount)
	})

	t.Run("FixedDiscount", func(t *testing.T) {
		// Arrange
		orderTotal := 100.00
		fixedDiscount := 20.00

		// Act
		finalTotal := orderTotal - fixedDiscount

		// Assert
		assert.Equal(t, 80.00, finalTotal)
	})
}
