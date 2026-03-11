package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// HASH TESTS
// ============================================================================

func TestHash_Success(t *testing.T) {
	// Arrange
	password := "TestPassword123!"

	// Act
	hash, err := Hash(password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHash_LongPassword(t *testing.T) {
	// Arrange
	// bcrypt has a 72 byte limit
	password := string(make([]byte, 73))
	for i := range password {
		password = password[:i] + "a" + password[i+1:]
	}

	// Act
	hash, err := Hash(password)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Equal(t, ErrPasswordTooLong, err)
}

func TestHash_Consistency(t *testing.T) {
	// Arrange
	password := "TestPassword123!"

	// Act
	hash1, _ := Hash(password)
	hash2, _ := Hash(password)

	// Assert
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	// bcrypt includes salt, so hashes will be different but both valid
	assert.NotEqual(t, hash1, hash2)
}

// ============================================================================
// VERIFY TESTS
// ============================================================================

func TestVerify_Success(t *testing.T) {
	// Arrange
	password := "TestPassword123!"
	hash, _ := Hash(password)

	// Act
	err := Verify(password, hash)

	// Assert
	assert.NoError(t, err)
}

func TestVerify_WrongPassword(t *testing.T) {
	// Arrange
	password := "TestPassword123!"
	wrongPassword := "WrongPassword"
	hash, _ := Hash(password)

	// Act
	err := Verify(wrongPassword, hash)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrHashMismatch, err)
}

func TestVerify_EmptyHash(t *testing.T) {
	// Arrange
	password := "TestPassword123!"

	// Act
	err := Verify(password, "")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidHash, err)
}

func TestVerify_EmptyPassword(t *testing.T) {
	// Arrange
	hash, _ := Hash("TestPassword123!")

	// Act
	err := Verify("", hash)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrHashMismatch, err)
}

// ============================================================================
// VALIDATION TESTS
// ============================================================================

func TestValidate_Success(t *testing.T) {
	// Arrange
	password := "TestPassword123!"
	reqs := DefaultRequirements()

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.NoError(t, err)
}

func TestValidate_TooShort(t *testing.T) {
	// Arrange
	password := "Test1!"
	reqs := DefaultRequirements()

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooShort, err)
}

func TestValidate_NoUppercase(t *testing.T) {
	// Arrange
	password := "testpassword123!"
	reqs := DefaultRequirements()

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooWeak, err)
}

func TestValidate_NoLowercase(t *testing.T) {
	// Arrange
	password := "TESTPASSWORD123!"
	reqs := DefaultRequirements()

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooWeak, err)
}

func TestValidate_NoNumber(t *testing.T) {
	// Arrange
	password := "TestPassword!"
	reqs := DefaultRequirements()

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooWeak, err)
}

func TestValidate_WithSpecial(t *testing.T) {
	// Arrange - password with special char
	password := "Test123!"
	reqs := Requirements{
		MinLength:        8,
		MaxLength:        72,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   true,
	}

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.NoError(t, err)
}

func TestValidate_NoSpecial(t *testing.T) {
	// Arrange - password without special char
	password := "TestPass12"
	reqs := Requirements{
		MinLength:        8,
		MaxLength:        72,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   true,
	}

	// Act
	err := Validate(password, reqs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPasswordTooWeak, err)
}

func TestValidateDefault(t *testing.T) {
	// Test valid password
	err := ValidateDefault("TestPassword123!")
	assert.NoError(t, err)

	// Test invalid password
	err = ValidateDefault("weak")
	assert.Error(t, err)
}

func TestValidateStrict(t *testing.T) {
	// Test strict valid password
	err := ValidateStrict("TestPassword123!@#")
	assert.NoError(t, err)

	// Test too short for strict
	err = ValidateStrict("Test123!")
	assert.Error(t, err)
}

// ============================================================================
// IS HASHED TESTS
// ============================================================================

func TestIsHashed_BcryptHash(t *testing.T) {
	// Arrange
	hash, _ := Hash("TestPassword123!")

	// Act
	result := IsHashed(hash)

	// Assert
	assert.True(t, result)
}

func TestIsHashed_PlainText(t *testing.T) {
	// Arrange
	password := "TestPassword123!"

	// Act
	result := IsHashed(password)

	// Assert
	assert.False(t, result)
}

func TestIsHashed_EmptyString(t *testing.T) {
	// Act
	result := IsHashed("")

	// Assert
	assert.False(t, result)
}

func TestIsHashed_ShortString(t *testing.T) {
	// Act
	result := IsHashed("short")

	// Assert
	assert.False(t, result)
}

// ============================================================================
// STRENGTH CHECK TESTS
// ============================================================================

func TestCheckStrength_VeryWeak(t *testing.T) {
	// Arrange
	password := "abc"

	// Act
	strength := CheckStrength(password)

	// Assert
	assert.Equal(t, StrengthVeryWeak, strength)
}

func TestCheckStrength_Weak(t *testing.T) {
	// Arrange
	password := "abcdef"

	// Act
	strength := CheckStrength(password)

	// Assert
	assert.Equal(t, StrengthWeak, strength)
}

func TestCheckStrength_Medium(t *testing.T) {
	// Arrange
	password := "TestPass"

	// Act
	strength := CheckStrength(password)

	// Assert
	assert.Equal(t, StrengthMedium, strength)
}

func TestCheckStrength_Strong(t *testing.T) {
	// Arrange
	password := "TestPassword123"

	// Act
	strength := CheckStrength(password)

	// Assert
	assert.Equal(t, StrengthStrong, strength)
}

func TestCheckStrength_VeryStrong(t *testing.T) {
	// Arrange
	password := "TestPassword123!@#$%^&"

	// Act
	strength := CheckStrength(password)

	// Assert
	assert.Equal(t, StrengthVeryStrong, strength)
}

func TestStrength_String(t *testing.T) {
	// Test all strength levels
	assert.Equal(t, "very_weak", StrengthVeryWeak.String())
	assert.Equal(t, "weak", StrengthWeak.String())
	assert.Equal(t, "medium", StrengthMedium.String())
	assert.Equal(t, "strong", StrengthStrong.String())
	assert.Equal(t, "very_strong", StrengthVeryStrong.String())
}

// ============================================================================
// HELPER FUNCTION TESTS
// ============================================================================

func TestHasUppercase(t *testing.T) {
	assert.True(t, hasUppercase("Test"))
	assert.False(t, hasUppercase("test"))
	assert.True(t, hasUppercase("TEST"))
}

func TestHasLowercase(t *testing.T) {
	assert.True(t, hasLowercase("Test"))
	assert.False(t, hasLowercase("TEST"))
	assert.True(t, hasLowercase("test"))
}

func TestHasNumber(t *testing.T) {
	assert.True(t, hasNumber("Test123"))
	assert.False(t, hasNumber("Test"))
	assert.True(t, hasNumber("123"))
}

func TestHasSpecial(t *testing.T) {
	assert.True(t, hasSpecial("Test!"))
	assert.False(t, hasSpecial("Test"))
	assert.True(t, hasSpecial("!@#$"))
}

// ============================================================================
// PASSWORD GENERATION TESTS
// ============================================================================

func TestGenerate_Success(t *testing.T) {
	// Act
	password, err := Generate(12)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, password, 12)
}

func TestGenerate_MinLength(t *testing.T) {
	// Act
	password, err := Generate(4)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, password, 8) // Minimum is 8
}

func TestGenerate_MaxLength(t *testing.T) {
	// Act
	password, err := Generate(100)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, password, 72) // Maximum is 72 (bcrypt limit)
}

func TestGenerate_Validation(t *testing.T) {
	// Generate multiple passwords and validate
	for i := 0; i < 10; i++ {
		password, _ := Generate(12)
		err := ValidateDefault(password)
		assert.NoError(t, err, "Generated password should pass validation")
	}
}

// ============================================================================
// NEEDS REHASH TESTS
// ============================================================================

func TestNeedsRehash(t *testing.T) {
	// Arrange
	hash, _ := HashWithCost("TestPassword123!", 10)

	// Act - Should not need rehash for same cost
	needsRehash, err := NeedsRehash(hash, 10)

	// Assert
	assert.NoError(t, err)
	assert.False(t, needsRehash)
}

func TestNeedsRehash_HigherCost(t *testing.T) {
	// Arrange
	hash, _ := HashWithCost("TestPassword123!", 10)

	// Act - Should need rehash for higher cost
	needsRehash, err := NeedsRehash(hash, 12)

	// Assert
	assert.NoError(t, err)
	assert.True(t, needsRehash)
}

func TestNeedsRehash_InvalidHash(t *testing.T) {
	// Arrange
	invalidHash := "short"

	// Act
	needsRehash, err := NeedsRehash(invalidHash, 12)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidHash, err)
	assert.False(t, needsRehash)
}
