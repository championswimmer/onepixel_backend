package db

import (
	"github.com/oschwald/geoip2-golang"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils"
	"onepixel_backend/src/utils/applogger"
	"os"
	"sync"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var appDb *gorm.DB        // singleton
var eventsDb *gorm.DB     // singleton
var reader *geoip2.Reader // singleton
var createAppDbOnce sync.Once
var createEventsDbOnce sync.Once
var createGeoIPDbOnce sync.Once

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

func GetAppDB() *gorm.DB {

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

	return appDb
}

func GetEventsDB() *gorm.DB {
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

		// try to automigrate
		err, success := lo.TryWithErrorValue(func() error {
			return eventsDb.AutoMigrate(&models.EventRedirect{})
		})
		if !success {
			applogger.Error("Events: AutoMigrate failed", err)
		}

	})

	return eventsDb
}

func GetGeoIPDB() *geoip2.Reader {

	// download file : https://git.io/GeoLite2-City.mmdb
	createGeoIPDbOnce.Do(func() {
		applogger.Warn("GeoIP: Initialising database")
		fresh := utils.IsFileFresh(30, "GeoLite2-City.mmdb")
		if !fresh {
			applogger.Error("GeoIP: GeoLite2-City.mmdb is not fresh; downloading again")
			lo.Try(func() error {
				return os.Remove("GeoLite2-City.mmdb")
			})
			lo.Must0(utils.DownloadFile("https://git.io/GeoLite2-City.mmdb", "GeoLite2-City.mmdb"))
			applogger.Info("GeoIP: GeoLite2-City.mmdb downloaded")
		}

		reader = lo.Must(geoip2.Open("GeoLite2-City.mmdb"))

	})

	return reader

}
