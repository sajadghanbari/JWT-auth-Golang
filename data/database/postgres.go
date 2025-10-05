package db

import (
	"JWT-Authentication-go/config"

	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	cnn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DbName,
		cfg.Postgres.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	dbClient = db

	sqlDb, _ := dbClient.DB()
	if err := sqlDb.Ping(); err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Minute)

	log.Println("✅ Database connected successfully")

	if cfg.Postgres.AutoMigrate {
		if err := runMigrations(); err != nil {
			return fmt.Errorf("migration error: %v", err)
		}
	}

	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	sqlDb, _ := dbClient.DB()
	sqlDb.Close()
}


//     AutoMigrate=true   

func runMigrations() error {
	log.Println("🛠 Running migrations...")
	return dbClient.AutoMigrate(Tables()...)
}