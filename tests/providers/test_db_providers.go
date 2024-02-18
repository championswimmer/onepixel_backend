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

var warnDbExistOnce sync.Once

func init() {
	db.InjectDBProvider("sqlite", ProvideSqliteDB)
	db.InjectDBProvider("duckdb", ProvideDuckDB)

	// Remove existing test databases
	warnDbExistOnce.Do(func() {
		cwd := lo.Must(os.Getwd())
		appDbPath := path.Join(cwd, config.DBUrl)
		eventDbPath := path.Join(cwd, config.EventDBUrl)
		if _, err := os.Stat(appDbPath); err == nil {
			applogger.Error("Test: app.db already exists")
		}
		if _, err := os.Stat(eventDbPath); err == nil {
			applogger.Error("Test: event.db already exists")
		}
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
