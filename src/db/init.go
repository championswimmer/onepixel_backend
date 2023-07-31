package db

import (
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=postgres user=postgres password=postgres dbname=onepixel port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	config := &gorm.Config{}
	config.Logger = logger.Default.LogMode(logger.Info)

	db := lo.Must(gorm.Open(postgres.Open(dsn), config))

	lo.Must0(db.AutoMigrate(&User{}))

	return db, nil
}
