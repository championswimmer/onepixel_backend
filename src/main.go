package main

import (
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"

	"github.com/samber/lo"
)

func main() {
	// Initialize the database
	db := lo.Must(db.InitDB(false))

	// Create the app
	app := server.CreateApp(db)

	log.Fatal(app.Listen(":3000"))
}
