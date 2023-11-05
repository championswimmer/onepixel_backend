package main

import (
	"github.com/samber/lo"
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
)

func main() {
	// Initialize the database
	db := lo.Must(db.InitDBProd())

	// Create the app
	app := server.CreateApp(db)

	log.Fatal(app.Listen(":3000")) // TODO: move port to external YAML config
}
