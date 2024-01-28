package db

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
	"sync"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var appDb *gorm.DB    // singleton
var eventsDb *gorm.DB // singleton
var createAppDbOnce sync.Once
var createEventsDbOnce sync.Once

func getGormConfig() (dbConfig *gorm.Config) {
	dbConfig = &gorm.Config{
		TranslateError: true,
	}
	var dbLogLevel logger.LogLevel = lo.Switch[string, logger.LogLevel](config.DBLogging).
		Case("info", logger.Info).
		Case("warn", logger.Warn).
		Case("error", logger.Error).
		Default(logger.Error)

	dbConfig.Logger = logger.Default.LogMode(dbLogLevel)
	return
}

func GetAppDB() (*gorm.DB, error) {

	createAppDbOnce.Do(func() {
		switch config.DBDialect {
		case "sqlite":
			applogger.Warn("App: Using sqlite db")
			appDb = lo.Must(gorm.Open(sqlite.Open(config.DBUrl), getGormConfig()))
			break
		case "postgres":
			applogger.Warn("App: Using postgres db")
			appDb = lo.Must(gorm.Open(postgres.Open(config.DBUrl), getGormConfig()))
			break
		default:
			panic("Database config incorrect")
		}

		lo.Must0(appDb.AutoMigrate(&models.User{}))
		lo.Must0(appDb.AutoMigrate(&models.UrlGroup{}))
		lo.Must0(appDb.AutoMigrate(&models.Url{}))
	})

	return appDb, nil
}

func GetEventsDB() (*gorm.DB, error) {
	createEventsDbOnce.Do(func() {

		applogger.Warn("Events: Using clickhouse db")

		eventsDb = lo.Must(gorm.Open(clickhouse.Open(config.EventDBUrl), getGormConfig()))

		// create table if not exists
		if _, err := eventsDb.Migrator().ColumnTypes((&models.EventRedirect{}).TableName()); err != nil {
			lo.Must0(eventsDb.AutoMigrate(&models.EventRedirect{}))
		}
	})

	return eventsDb, nil
}
