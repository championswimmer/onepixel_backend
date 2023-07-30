package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=postgres user=postgres password=postgres dbname=onepixel port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	config := &gorm.Config{}
	db, err := gorm.Open(postgres.Open(dsn), config)

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
