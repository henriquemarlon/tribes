package configs

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
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
		entity.SocialAccount{},
	); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	var users []map[string]interface{}
	if path == ":memory:" {
		users = []map[string]interface{}{
			{
				"role":                "admin",
				"address":             common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
				"investment_limit":    uint256.NewInt(0),
				"debt_issuance_limit": uint256.NewInt(0),
				"created_at":          0,
			},
		}
	} else {
		users = []map[string]interface{}{
			{
				"role":                "admin",
				"address":             common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").String(),
				"investment_limit":    uint256.NewInt(0),
				"debt_issuance_limit": uint256.NewInt(0),
				"created_at":          0,
			},
		}
	}

	if err := db.Table("users").Create(users).Error; err != nil {
		return nil, fmt.Errorf("failed to create users: %w", err)
	}
	return db, nil
}
