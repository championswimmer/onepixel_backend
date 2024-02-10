package db

import (
	"github.com/samber/lo"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"onepixel_backend/src/utils/applogger"
)

func ProvidePostgresDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("App: Using postgres db")
	return lo.Must(gorm.Open(postgres.Open(dbUrl), config))
}

func ProvideClickhouseDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("App: Using clickhouse db")
	return lo.Must(gorm.Open(clickhouse.Open(dbUrl), config))
}
