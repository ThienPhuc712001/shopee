package database

import (
	"fmt"
	"log"
	"time"

	"ecommerce/pkg/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// NewConnection creates a new database connection using configuration
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dbConfig := Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	return Connect(dbConfig)
}

// Connect establishes a connection to SQL Server
func Connect(cfg Config) (*gorm.DB, error) {
	// Build connection string for SQL Server
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=false&trustservercertificate=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	// Configure GORM logger based on environment
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection with connection pool settings
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := testConnection(db); err != nil {
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	log.Println("✓ Database connection established successfully")
	log.Printf("  - Max idle connections: 10")
	log.Printf("  - Max open connections: 100")
	log.Printf("  - Connection lifetime: 1 hour")

	DB = db
	return db, nil
}

// testConnection verifies the database connection
func testConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("✓ Database ping successful")
	return nil
}

// AutoMigrate runs auto migration for given models
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("Starting database migration...")
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("✓ Database migration completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	log.Println("Closing database connection...")
	return sqlDB.Close()
}
