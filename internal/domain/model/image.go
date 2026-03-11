package model

// ImageType represents the type of image upload
type ImageType string

const (
	ImageTypeProduct  ImageType = "product"
	ImageTypeReview   ImageType = "review"
	ImageTypeAvatar   ImageType = "avatar"
)

// ImageUpload represents the upload configuration
type ImageUpload struct {
	MaxSize     int64    // Maximum file size in bytes
	AllowedExts []string // Allowed file extensions
	UploadPath  string   // Upload directory
}

// GetUploadConfig returns configuration based on image type
func GetUploadConfig(imageType ImageType) ImageUpload {
	switch imageType {
	case ImageTypeProduct:
		return ImageUpload{
			MaxSize:     5 * 1024 * 1024, // 5MB
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			UploadPath:  "uploads/products",
		}
	case ImageTypeReview:
		return ImageUpload{
			MaxSize:     5 * 1024 * 1024, // 5MB
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			UploadPath:  "uploads/reviews",
		}
	case ImageTypeAvatar:
		return ImageUpload{
			MaxSize:     2 * 1024 * 1024, // 2MB for avatars
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			UploadPath:  "uploads/avatars",
		}
	default:
		return ImageUpload{
			MaxSize:     5 * 1024 * 1024,
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			UploadPath:  "uploads/other",
		}
	}
}
