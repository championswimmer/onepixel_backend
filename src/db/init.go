package db

import (
	"onepixel_backend/src/models"

	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(test bool) (*gorm.DB, error) {
	dsn := "host=postgres user=postgres password=postgres dbname=onepixel port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	config := &gorm.Config{}
	config.Logger = logger.Default.LogMode(logger.Info)
	var db *gorm.DB

	if test {
		db = lo.Must(gorm.Open(sqlite.Open("file::memory:?cache=shared"), config))
	} else {
		db = lo.Must(gorm.Open(postgres.Open(dsn), config))
	}

	lo.Must0(db.AutoMigrate(&models.User{}))
	lo.Must0(db.AutoMigrate(&models.Url{}))

	return db, nil
}
