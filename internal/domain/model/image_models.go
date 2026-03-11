package model

import (
	"time"
)

// ReviewImage represents a review image in the database
type ReviewImage struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID  int64     `gorm:"not null;index" json:"review_id"`
	URL       string    `gorm:"type:varchar(500);not null" json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Review *Review `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
}

// TableName specifies the table name for ReviewImage
func (ReviewImage) TableName() string {
	return "review_images"
}

// UserAvatar represents a user avatar in the database
type UserAvatar struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;uniqueIndex" json:"user_id"`
	URL       string    `gorm:"type:varchar(500);not null" json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for UserAvatar
func (UserAvatar) TableName() string {
	return "user_avatars"
}

// ImageUploadLog tracks upload history for auditing
type ImageUploadLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"index" json:"user_id"`
	ImageType    string    `gorm:"type:varchar(50)" json:"image_type"`
	OriginalName string    `gorm:"type:varchar(255)" json:"original_name"`
	StoredName   string    `gorm:"type:varchar(255)" json:"stored_name"`
	FileSize     int64     `json:"file_size"`
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName specifies the table name for ImageUploadLog
func (ImageUploadLog) TableName() string {
	return "image_upload_logs"
}
