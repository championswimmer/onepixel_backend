package main

import (
	"fmt"
	"github.com/samber/lo"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"os"
)

func main() {
	// Initialize the database
	db := lo.Must(db.GetProdDB())

	// Create the app
	app := server.CreateApp(db)

	httpPort, _ := lo.Coalesce(os.Getenv("PORT"), "3000")

	// TODO: move port to external YAML config
	utils.AppLogger.Fatal(app.Listen(fmt.Sprintf(":%s", httpPort)))
}
