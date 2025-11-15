package config

import (
	"log"

	"github.com/mhcvintar/repo-gen/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        return nil, err
    }

    log.Println("In-memory SQLite database initialized successfully")
    return db, nil
}