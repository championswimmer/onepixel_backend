package config

import (
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"onepixel_backend/src/utils/applogger"
	"os"
	"path"
	"runtime"
)

func init() {
	if os.Getenv("ENV") == "test" {
		// for tests, chdir to the project root
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
		lo.Must0(os.Chdir(dir))
		if err := godotenv.Load("onepixel.test.env"); err != nil {
			applogger.Error(err)
		}
	}

	if os.Getenv("ENV") == "production" || os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
		if err := godotenv.Load("onepixel.production.env"); err != nil {
			applogger.Error(err)
		}
	}

	// Use defaults from local.env for all missing vars
	if err := godotenv.Load("onepixel.local.env"); err != nil {
		applogger.Error(err)
	}
}
