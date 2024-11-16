package configs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupSQlite(path string) (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = db.AutoMigrate(
		entity.Order{},
		entity.User{},
		entity.Crowdfunding{},
		entity.Crowdfunding{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	err = db.Table("users").Create([]map[string]interface{}{
		{
			"role":       "admin",
			"address":    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
			"created_at": 0,
		},
		{
			"role":       "",
			"address":    common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4").String(),
			"created_at": 0,
		},
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create users: %w", err)
	}
	return db, nil
}

func SetupSQliteMemory(path string) (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = db.AutoMigrate(
		entity.Order{},
		entity.User{},
		entity.Crowdfunding{},
		entity.Crowdfunding{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	err = db.Table("users").Create([]map[string]interface{}{
		{
			"role":       "admin",
			"address":    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
			"created_at": 0,
		},
		{
			"role":       "",
			"address":    common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4").String(),
			"created_at": 0,
		},
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create users: %w", err)
	}
	return db, nil
}
