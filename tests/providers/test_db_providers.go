package providers

import (
	"github.com/championswimmer/duckdb-driver/duckdb"
	"github.com/samber/lo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"onepixel_backend/src/db"
	"onepixel_backend/src/utils/applogger"
)

func init() {
	db.InjectDBProvider("sqlite", ProvideSqliteDB)
	db.InjectDBProvider("duckdb", ProvideDuckDB)

}

func ProvideSqliteDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("Test: Using sqlite db")
	return lo.Must(gorm.Open(sqlite.Open(dbUrl), config))
}

func ProvideDuckDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("Test: Using duckdb db")
	return lo.Must(gorm.Open(duckdb.Open(dbUrl), config))
}
