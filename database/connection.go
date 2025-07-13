package database

import (
	"JWT-Authentication-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB  *gorm.DB


func Connect() (*gorm.DB, error){
	dsn := "root:@tcp(localhost:8500)/jwt_go"

	db,err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	db.AutoMigrate(&models.User{})
	return db, nil
}
