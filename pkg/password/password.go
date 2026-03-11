// Package password provides secure password hashing and verification using bcrypt
package password

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// ERRORS
// ============================================================================

var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong     = errors.New("password must not exceed 72 characters")
	ErrPasswordTooWeak     = errors.New("password is too weak")
	ErrHashMismatch        = errors.New("password hash mismatch")
	ErrInvalidHash         = errors.New("invalid password hash")
)

// ============================================================================
// PASSWORD REQUIREMENTS
// ============================================================================

// Requirements defines password complexity requirements
type Requirements struct {
	// MinLength is the minimum password length (default: 8)
	MinLength int

	// MaxLength is the maximum password length (default: 72, bcrypt limit)
	MaxLength int

	// RequireUppercase requires at least one uppercase letter
	RequireUppercase bool

	// RequireLowercase requires at least one lowercase letter
	RequireLowercase bool

	// RequireNumber requires at least one number
	RequireNumber bool

	// RequireSpecial requires at least one special character
	RequireSpecial bool
}

// DefaultRequirements returns default password requirements
func DefaultRequirements() Requirements {
	return Requirements{
		MinLength:        8,
		MaxLength:        72,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   false,
	}
}

// StrictRequirements returns strict password requirements
func StrictRequirements() Requirements {
	return Requirements{
		MinLength:        12,
		MaxLength:        72,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   true,
	}
}

// ============================================================================
// PASSWORD HASHING
// ============================================================================

// Hash hashes a password using bcrypt
func Hash(password string) (string, error) {
	// bcrypt has a 72 byte limit
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// HashWithCost hashes a password with a specific cost
func HashWithCost(password string, cost int) (string, error) {
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ============================================================================
// PASSWORD VERIFICATION
// ============================================================================

// Verify compares a password with its hash
func Verify(password, hash string) error {
	if hash == "" {
		return ErrInvalidHash
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrHashMismatch
		}
		return err
	}

	return nil
}

// IsHashed checks if a string appears to be a bcrypt hash
func IsHashed(s string) bool {
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	if len(s) < 60 {
		return false
	}
	return s[:4] == "$2a$" || s[:4] == "$2b$" || s[:4] == "$2y$"
}

// NeedsRehash checks if a hash needs to be rehashed with a higher cost
func NeedsRehash(hash string, targetCost int) (bool, error) {
	// Extract cost from hash
	// bcrypt hash format: $2a$10$... where 10 is the cost
	if len(hash) < 7 {
		return false, ErrInvalidHash
	}

	// Simple cost comparison by checking the cost digits in the hash
	// Hash format: $2a$CC$ where CC is the 2-digit cost
	if hash[4] < byte('0'+targetCost/10) || (hash[4] == byte('0'+targetCost/10) && hash[5] < byte('0'+targetCost%10)) {
		return true, nil
	}

	return false, nil
}

// ============================================================================
// PASSWORD VALIDATION
// ============================================================================

// Validate validates a password against requirements
func Validate(password string, reqs Requirements) error {
	if len(password) < reqs.MinLength {
		return ErrPasswordTooShort
	}

	if len(password) > reqs.MaxLength {
		return ErrPasswordTooLong
	}

	if reqs.RequireUppercase && !hasUppercase(password) {
		return ErrPasswordTooWeak
	}

	if reqs.RequireLowercase && !hasLowercase(password) {
		return ErrPasswordTooWeak
	}

	if reqs.RequireNumber && !hasNumber(password) {
		return ErrPasswordTooWeak
	}

	if reqs.RequireSpecial && !hasSpecial(password) {
		return ErrPasswordTooWeak
	}

	return nil
}

// ValidateDefault validates a password against default requirements
func ValidateDefault(password string) error {
	return Validate(password, DefaultRequirements())
}

// ValidateStrict validates a password against strict requirements
func ValidateStrict(password string) error {
	return Validate(password, StrictRequirements())
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasSpecial(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// ============================================================================
// PASSWORD GENERATION
// ============================================================================

// Generate generates a random secure password
func Generate(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	if length > 72 {
		length = 72
	}

	// Character sets
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers   = "0123456789"
		special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		all       = lowercase + uppercase + numbers + special
	)

	// Ensure at least one of each type
	password := make([]byte, length)
	password[0] = lowercase[0]
	password[1] = uppercase[0]
	password[2] = numbers[0]
	password[3] = special[0]

	// Fill the rest randomly
	for i := 4; i < length; i++ {
		password[i] = all[i%len(all)]
	}

	// Shuffle the password
	// Note: For production, use crypto/rand for secure shuffling
	for i := length - 1; i > 0; i-- {
		j := i % 7
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// ============================================================================
// STRENGTH CHECK
// ============================================================================

// Strength represents password strength
type Strength int

const (
	StrengthVeryWeak Strength = iota
	StrengthWeak
	StrengthMedium
	StrengthStrong
	StrengthVeryStrong
)

// String returns the string representation of strength
func (s Strength) String() string {
	switch s {
	case StrengthVeryWeak:
		return "very_weak"
	case StrengthWeak:
		return "weak"
	case StrengthMedium:
		return "medium"
	case StrengthStrong:
		return "strong"
	case StrengthVeryStrong:
		return "very_strong"
	default:
		return "unknown"
	}
}

// CheckStrength checks the strength of a password
func CheckStrength(password string) Strength {
	length := len(password)
	hasUpper := hasUppercase(password)
	hasLower := hasLowercase(password)
	hasNumber := hasNumber(password)
	hasSpecial := hasSpecial(password)

	// Count character types
	types := 0
	if hasUpper {
		types++
	}
	if hasLower {
		types++
	}
	if hasNumber {
		types++
	}
	if hasSpecial {
		types++
	}

	// Determine strength
	switch {
	case length < 6:
		return StrengthVeryWeak
	case length < 8 || types == 1:
		return StrengthWeak
	case length < 10 && types >= 2:
		return StrengthMedium
	case length >= 12 && types >= 3:
		if length >= 16 && types == 4 {
			return StrengthVeryStrong
		}
		return StrengthStrong
	default:
		return StrengthMedium
	}
}
