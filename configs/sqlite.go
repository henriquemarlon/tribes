package configs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupSQlite() (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open("tribes.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = db.AutoMigrate(
		&entity.Bid{},
		&entity.User{},
		&entity.Auction{},
		&entity.Contract{},
	)

	db.Create(&entity.User{
		Role:      "admin",
		Username:  "admin",
		Address:   custom_type.NewAddress(common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")),
		CreatedAt: 0,
	})

	db.Create(&entity.User{
		Role:      "auctioneer",
		Username:  "auctioneer",
		Address:   custom_type.NewAddress(common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4")),
		CreatedAt: 0,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	return db, nil
}

func SetupSQliteMemory() (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	err = db.AutoMigrate(
		&entity.Bid{},
		&entity.User{},
		&entity.Auction{},
		&entity.Contract{},
	)

	db.Create(&entity.User{
		Role:      "admin",
		Username:  "admin",
		Address:   custom_type.NewAddress(common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")),
		CreatedAt: 0,
	})

	db.Create(&entity.User{
		Role:      "auctioneer",
		Username:  "auctioneer",
		Address:   custom_type.NewAddress(common.HexToAddress("0xf49Fc2E6478982F125c0F38d38f67B32772604B4")),
		CreatedAt: 0,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	return db, nil
}
