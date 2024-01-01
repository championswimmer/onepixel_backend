package main

import (
	"fmt"
	"github.com/samber/lo"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"onepixel_backend/src/utils/applogger"
)

func main() {
	// Initialize the database
	db := lo.Must(db.GetProdDB())

	// Create the app
	app := server.CreateApp(db)

	// TODO: move port to external YAML config
	applogger.Fatal(app.Listen(fmt.Sprintf(":%s", config.Port)))
}
