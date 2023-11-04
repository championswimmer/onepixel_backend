package main

import (
	"github.com/samber/lo"
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
)

func main() {
	// Initialize the database
	db := lo.Must(db.InitDB(false))

	// Create the app
	app := server.CreateApp(db)

	log.Fatal(app.Listen(":3000"))
}
