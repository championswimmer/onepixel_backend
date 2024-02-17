package providers

import (
	"github.com/championswimmer/duckdb-driver/duckdb"
	_ "github.com/flashlabs/rootpath"
	"github.com/samber/lo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db"
	"onepixel_backend/src/utils/applogger"
	"os"
	"path"
	"sync"
	"time"
)

var removeDbOnce sync.Once

func init() {
	db.InjectDBProvider("sqlite", ProvideSqliteDB)
	db.InjectDBProvider("duckdb", ProvideDuckDB)

	// Remove existing test databases
	removeDbOnce.Do(func() {
		cwd := lo.Must(os.Getwd())
		appDbPath := path.Join(cwd, config.DBUrl)
		eventDbPath := path.Join(cwd, config.EventDBUrl)
		applogger.Debug("App: Removing existing test database file at", appDbPath)
		lo.Must0(os.RemoveAll(appDbPath))
		applogger.Debug("Events: Removing existing test database file at", eventDbPath)
		lo.Must0(os.RemoveAll(eventDbPath))
		lo.Must0(os.RemoveAll(eventDbPath + ".wal"))
		time.Sleep(1 * time.Second)
	})

}

func ProvideSqliteDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("Test: Using sqlite db")
	return lo.Must(gorm.Open(sqlite.Open(dbUrl), config))
}

func ProvideDuckDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("Test: Using duckdb db")
	return lo.Must(gorm.Open(duckdb.Open(dbUrl), config))
}
