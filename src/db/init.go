package db

import (
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm/logger"
	"onepixel_backend/src/models"
	"os"

	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDBTest() (*gorm.DB, error) {
	log.Warn(logger.YellowBold, "Using test database")
	return initDB(true)
}

func InitDBProd() (*gorm.DB, error) {
	log.Warn(logger.YellowBold, "Using production database")
	return initDB(false)
}

func initDB(test bool) (*gorm.DB, error) {
	// TODO: move db config to external YAML config
	dsn, _ := lo.Coalesce(
		os.Getenv("DATABASE_PRIVATE_URL"),
		os.Getenv("DATABASE_URL"),
		"host=postgres user=postgres password=postgres dbname=onepixel port=5432 sslmode=disable TimeZone=UTC",
	)

	config := &gorm.Config{
		TranslateError: true,
	}
	var db *gorm.DB

	if test {
		db = lo.Must(gorm.Open(sqlite.Open("file::memory:?cache=shared"), config))
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		db = lo.Must(gorm.Open(postgres.Open(dsn), config))
		config.Logger = logger.Default.LogMode(logger.Error)
	}

	lo.Must0(db.AutoMigrate(&models.User{}))
	lo.Must0(db.AutoMigrate(&models.Url{}))

	return db, nil
}
