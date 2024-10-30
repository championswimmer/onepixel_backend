package config

import (
	"os"
	"strconv"

	"github.com/samber/lo"
)

var Env string

var DBLogging string
var DBDialect string
var DBUrl string
var UseFileDB bool

var EventDBUrl string
var EventDBDialect string

var Port string
var MainHost string
var AdminHost string
var RedirUrlBase string

var AdminApiKey string
var AdminUserEmail string

var JwtSigningKey string
var JwtDurationDays int

var PosthogApiKey string

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
	UseFileDB, _ = strconv.ParseBool(os.Getenv("USE_FILE_DB"))
	EventDBDialect = os.Getenv("EVENTDB_DIALECT")
	EventDBUrl, _ = lo.Coalesce(
		os.Getenv("EVENTDB_PRIVATE_URL"),
		os.Getenv("EVENTDB_URL"),
	)
	Port = os.Getenv("PORT")
	MainHost = os.Getenv("MAIN_SITE_HOST")
	AdminHost = os.Getenv("ADMIN_SITE_HOST")
	RedirUrlBase = "http://" + MainHost + ":" + Port + "/"
	if Env == "production" {
		RedirUrlBase = "https://" + MainHost + "/"
	}
	AdminApiKey = os.Getenv("ADMIN_API_KEY")
	AdminUserEmail = os.Getenv("ADMIN_USER_EMAIL")
	JwtSigningKey = os.Getenv("JWT_SIGNING_KEY")
	JwtDurationDays, _ = strconv.Atoi(os.Getenv("JWT_DURATION_DAYS"))
	PosthogApiKey = os.Getenv("POSTHOG_API_KEY")
}
