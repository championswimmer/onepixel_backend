package config

import (
	"github.com/samber/lo"
	"os"
)

var DBLogging string
var DBDialect string
var DBUrl string
var Port string
var MainHost string
var AdminHost string

// should run after env.go#init as this `vars` is alphabetically after `env`
func init() {
	DBLogging = os.Getenv("DB_LOGGING")
	DBDialect = os.Getenv("DB_DIALECT")
	DBUrl, _ = lo.Coalesce(
		os.Getenv("DATABASE_PRIVATE_URL"),
		os.Getenv("DATABASE_URL"),
	)
	Port = os.Getenv("PORT")
	MainHost = os.Getenv("MAIN_SITE_HOST")
	AdminHost = os.Getenv("ADMIN_SITE_HOST")
}
