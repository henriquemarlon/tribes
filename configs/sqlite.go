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
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(
		entity.User{},
		entity.Order{},
		entity.Contract{},
		entity.Crowdfunding{},
	); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	var users []map[string]interface{}
	if path == ":memory:" {
		users = []map[string]interface{}{
			{
				"role":       "admin",
				"address":    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
				"created_at": 0,
			},
			{
				"role":       "creator",
				"address":    common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4").String(),
				"created_at": 0,
			},
		}
	} else {
		users = []map[string]interface{}{
			{
				"role":       "admin",
				"address":    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
				"created_at": 0,
			},
			{
				"role":       "user",
				"address":    common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4").String(),
				"created_at": 0,
			},
		}
	}

	if err := db.Table("users").Create(users).Error; err != nil {
		return nil, fmt.Errorf("failed to create users: %w", err)
	}

	return db, nil
}
