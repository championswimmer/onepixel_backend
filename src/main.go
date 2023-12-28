package main

import (
	"fmt"
	"github.com/samber/lo"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"onepixel_backend/src/utils/applogger"
	"os"
)

func main() {
	// Initialize the database
	db := lo.Must(db.GetProdDB())

	// Create the app
	app := server.CreateApp(db)

	httpPort, _ := lo.Coalesce(os.Getenv("PORT"), "3000")

	// TODO: move port to external YAML config
	applogger.Fatal(app.Listen(fmt.Sprintf(":%s", httpPort)))
}
