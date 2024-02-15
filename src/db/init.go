package db

import (
	"errors"
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

type DatabaseProvider func(dbUrl string, config *gorm.Config) *gorm.DB

var dbProviders map[string]DatabaseProvider = map[string]DatabaseProvider{}

func InjectDBProvider(name string, provider DatabaseProvider) {
	dbProviders[name] = provider
}

func init() {
	InjectDBProvider("postgres", ProvidePostgresDB)
	InjectDBProvider("clickhouse", ProvideClickhouseDB)
}

func GetAppDB() (*gorm.DB, error) {

	createAppDbOnce.Do(func() {
		applogger.Warn("App: Initialising database")
		switch config.DBDialect {
		case "sqlite":
			appDb = dbProviders["sqlite"](config.DBUrl, getGormConfig())
			break
		case "postgres":
			appDb = dbProviders["postgres"](config.DBUrl, getGormConfig())
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
		applogger.Warn("Events: Initialising database")

		switch config.EventDBDialect {
		case "clickhouse":
			eventsDb = dbProviders["clickhouse"](config.EventDBUrl, getGormConfig())
			break
		case "duckdb":
			eventsDb = dbProviders["duckdb"](config.EventDBUrl, getGormConfig())
			break
		default:
			panic("EventDB config incorrect")
		}

		// automigrate table if we cannot get column types
		lo.TryCatchWithErrorValue(func() error {
			if eventsDb.Migrator().HasTable((&models.EventRedirect{}).TableName()) {
				applogger.Info("Events: table exists")
				return nil
			} else {
				return errors.New("table not found")
			}
			//_, err := eventsDb.Migrator().ColumnTypes((&models.EventRedirect{}).TableName())
			//return err
		}, func(e any) {
			applogger.Error("Error reading column types of eventsdb: " + e.(error).Error())
			lo.Must0(eventsDb.AutoMigrate(&models.EventRedirect{}))
			applogger.Info("Events: table automigrated")

		})

	})

	return eventsDb, nil
}
