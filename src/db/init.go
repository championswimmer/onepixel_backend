package db

import (
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"onepixel_backend/src/models"
	"onepixel_backend/src/utils/applogger"
	"os"
)

var dbSingleton *gorm.DB

func GetTestDB() (*gorm.DB, error) {
	applogger.Warn("Using test database")
	return getOrInitDB(true)
}

func GetProdDB() (*gorm.DB, error) {
	applogger.Warn("Using production database")
	return getOrInitDB(false)
}

func getOrInitDB(test bool) (*gorm.DB, error) {

	if dbSingleton != nil {
		return dbSingleton, nil
	}
	// TODO: move db config to external YAML config
	dsn, _ := lo.Coalesce(
		os.Getenv("DATABASE_PRIVATE_URL"),
		os.Getenv("DATABASE_URL"),
		"host=postgres user=postgres password=postgres dbname=onepixel port=5432 sslmode=disable TimeZone=UTC",
	)

	config := &gorm.Config{
		TranslateError: true,
	}

	if test {
		config.Logger = logger.Default.LogMode(logger.Info)
		dbSingleton = lo.Must(gorm.Open(sqlite.Open("file::memory:?cache=shared"), config))
	} else {
		config.Logger = logger.Default.LogMode(logger.Error)
		dbSingleton = lo.Must(gorm.Open(postgres.Open(dsn), config))
	}

	lo.Must0(dbSingleton.AutoMigrate(&models.User{}))
	lo.Must0(dbSingleton.AutoMigrate(&models.UrlGroup{}))
	lo.Must0(dbSingleton.AutoMigrate(&models.Url{}))

	return dbSingleton, nil
}
