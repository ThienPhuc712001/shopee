package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// File validation constants
const (
	MaxImageSize      = 5 * 1024 * 1024  // 5MB
	MaxAvatarSize     = 2 * 1024 * 1024  // 2MB
	AllowedMimeTypes  = "image/jpeg,image/png,image/webp"
)

// AllowedExtensions defines permitted file extensions
var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// FileError represents file operation errors
type FileError struct {
	Code    string
	Message string
}

func (e *FileError) Error() string {
	return e.Message
}

// File validation errors
var (
	ErrFileTooLarge   = &FileError{Code: "FILE_TOO_LARGE", Message: "File size exceeds maximum allowed size"}
	ErrInvalidFileType = &FileError{Code: "INVALID_FILE_TYPE", Message: "File type is not allowed"}
	ErrNoFile         = &FileError{Code: "NO_FILE", Message: "No file provided"}
	ErrInvalidFile    = &FileError{Code: "INVALID_FILE", Message: "Invalid file"}
)

// ValidateFile checks if the uploaded file meets requirements
func ValidateFile(file *multipart.FileHeader, maxSize int64, allowedExts []string) error {
	// Check file size
	if file.Size > maxSize {
		return ErrFileTooLarge
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return ErrInvalidFileType
	}

	// Open file to check MIME type
	src, err := file.Open()
	if err != nil {
		return ErrInvalidFile
	}
	defer src.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return ErrInvalidFile
	}

	// Reset file reader
	src.Close()

	return nil
}

// ValidateImageMIMEType checks if the file is actually an image
func ValidateImageMIMEType(file multipart.File) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return ErrInvalidFile
	}

	// Reset file reader
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	allowedTypes := strings.Split(AllowedMimeTypes, ",")

	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return nil
		}
	}

	return ErrInvalidFileType
}

// GenerateUniqueFilename creates a unique filename using UUID
func GenerateUniqueFilename(originalFilename string) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	uuidStr := uuid.New().String()
	return uuidStr + ext
}

// GenerateTimestampFilename creates a filename using timestamp + random string
func GenerateTimestampFilename(originalFilename string) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return timestamp + "_" + randomStr + ext
}

// SanitizeFilename removes dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators to prevent directory traversal
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "..", "")

	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Limit filename length
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		filename = filename[:255-len(ext)] + ext
	}

	return filename
}

// EnsureDir creates directory if it doesn't exist
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// SaveFile saves uploaded file to specified path
func SaveFile(file multipart.File, filePath string) (int64, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return 0, err
	}

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	// Copy file content
	written, err := io.Copy(dst, file)
	if err != nil {
		return 0, err
	}

	return written, nil
}

// DeleteFile removes a file from storage
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetFileSize returns the size of a file
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileExtension returns the extension of a filename
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsAllowedExtension checks if file extension is allowed
func IsAllowedExtension(filename string, allowedExts []string) bool {
	ext := GetFileExtension(filename)
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
