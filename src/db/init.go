package db

import (
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
)

var dbSingleton *gorm.DB

func GetDB() (*gorm.DB, error) {

	// TODO: thread safety?
	if dbSingleton != nil {
		return dbSingleton, nil
	}

	dbConfig := &gorm.Config{
		TranslateError: true,
	}
	var dbLogLevel logger.LogLevel = lo.Switch[string, logger.LogLevel](config.DBLogging).
		Case("info", logger.Info).
		Case("warn", logger.Warn).
		Case("error", logger.Error).
		Default(logger.Error)

	dbConfig.Logger = logger.Default.LogMode(dbLogLevel)
	switch config.DBDialect {
	case "sqlite":
		applogger.Warn("Using sqlite db")
		dbSingleton = lo.Must(gorm.Open(sqlite.Open(config.DBUrl), dbConfig))
		break
	case "postgres":
		applogger.Warn("Using postgres db")
		dbSingleton = lo.Must(gorm.Open(postgres.Open(config.DBUrl), dbConfig))
		break
	default:
		panic("Database config incorrect")
	}

	lo.Must0(dbSingleton.AutoMigrate(&models.User{}))
	lo.Must0(dbSingleton.AutoMigrate(&models.UrlGroup{}))
	lo.Must0(dbSingleton.AutoMigrate(&models.Url{}))

	return dbSingleton, nil
}
