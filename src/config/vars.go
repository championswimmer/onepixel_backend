package config

import (
	"github.com/samber/lo"
	"os"
	"strconv"
)

var Env string

var DBLogging string
var DBDialect string
var DBUrl string

var Port string
var MainHost string
var AdminHost string

var AdminApiKey string

var JwtSigningKey string
var JwtDurationDays int

// should run after env.go#init as this `vars` is alphabetically after `env`
func init() {
	Env, _ = lo.Coalesce(
		os.Getenv("RAILWAY_ENVIRONMENT"),
		os.Getenv("ENV"),
		"local",
	)
	DBLogging = os.Getenv("DB_LOGGING")
	DBDialect = os.Getenv("DB_DIALECT")
	DBUrl, _ = lo.Coalesce(
		os.Getenv("DATABASE_PRIVATE_URL"),
		os.Getenv("DATABASE_URL"),
	)
	Port = os.Getenv("PORT")
	MainHost = os.Getenv("MAIN_SITE_HOST")
	AdminHost = os.Getenv("ADMIN_SITE_HOST")
	AdminApiKey = os.Getenv("ADMIN_API_KEY")
	JwtSigningKey = os.Getenv("JWT_SIGNING_KEY")
	JwtDurationDays, _ = strconv.Atoi(os.Getenv("JWT_DURATION_DAYS"))
}
