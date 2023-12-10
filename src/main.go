package main

import (
	"fmt"
	"github.com/samber/lo"
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"os"
)

func main() {
	// Initialize the database
	db := lo.Must(db.InitDBProd())

	// Create the app
	app := server.CreateApp(db)

	httpPort, _ := lo.Coalesce(os.Getenv("PORT"), "3000")

	// TODO: move port to external YAML config
	log.Fatal(app.Listen(fmt.Sprintf(":%s", httpPort)))
}
