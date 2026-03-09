package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"regexp"
	"strings"
	"time"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	random, _ := rand.Int(rand.Reader, big.NewInt(100000))
	return "ORD-" + timestamp + "-" + random.String()
}

// GenerateTransactionID generates a unique transaction ID
func GenerateTransactionID() string {
	timestamp := time.Now().Format("20060102150405")
	random, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return "TXN-" + timestamp + "-" + random.String()
}

// IsValidEmail validates an email address
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidPhone validates a phone number
func IsValidPhone(phone string) bool {
	pattern := `^[\+]?[(]?[0-9]{1,3}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// TruncateString truncates a string to specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")

	// Remove multiple hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// FormatPrice formats a price for display
func FormatPrice(price float64, currency string) string {
	return currency + " " + FormatNumber(price)
}

// FormatNumber formats a number with thousand separators
func FormatNumber(num float64) string {
	str := strings.Split(strings.Split(strings.TrimRight(strings.TrimRight(
		string(rune(int(num))), "0"), "."), ".")[0], "")
	
	var result []string
	for i, v := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ",")
		}
		result = append(result, v)
	}
	return strings.Join(result, "")
}

// CalculateDiscount calculates the discount percentage
func CalculateDiscount(originalPrice, salePrice float64) int {
	if originalPrice <= 0 {
		return 0
	}
	discount := ((originalPrice - salePrice) / originalPrice) * 100
	return int(discount)
}

// IsInStock checks if a product is in stock
func IsInStock(stock, quantity int) bool {
	return stock >= quantity && stock > 0
}

// ParseDate parses a date string
func ParseDate(dateStr, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02"
	}
	return time.Parse(layout, dateStr)
}

// FormatDate formats a time to string
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02"
	}
	return t.Format(layout)
}

// GetTimeAgo returns a human-readable time ago string
func GetTimeAgo(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return formatDuration(diff.Minutes(), "minute")
	case diff < 24*time.Hour:
		return formatDuration(diff.Hours(), "hour")
	case diff < 30*24*time.Hour:
		return formatDuration(diff.Hours()/24, "day")
	case diff < 12*30*24*time.Hour:
		return formatDuration(diff.Hours()/(24*30), "month")
	default:
		return formatDuration(diff.Hours()/(24*365), "year")
	}
}

func formatDuration(value float64, unit string) string {
	v := int(value)
	if v == 1 {
		return "1 " + unit + " ago"
	}
	return string(rune(v)) + " " + unit + "s ago"
}
