package main

import (
	"github.com/samber/lo"
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
)

func main() {
	app := server.CreateApp()

	// Initialize the database
	lo.Must(db.InitDB())

	log.Fatal(app.Listen(":3000"))
}
