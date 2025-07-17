package database

import (
	"log"

	"github.com/latttchc/finding-forest-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect はデータベースに接続する
func Connect(dsn string) (*gorm.DB, error) {
	// GORM設定
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// データベース接続
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	// 接続プールの設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 接続プールの設定
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connected successfully")
	return db, nil
}

// Migrate はデータベースマイグレーションを実行する
func Migrate(db *gorm.DB) error {
	// model migrate
	err := db.AutoMigrate(
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}
