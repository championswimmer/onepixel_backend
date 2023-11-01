package main

import (
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
)

func main() {
	// Initialize the database
	dbConnection, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Pass the dbConnection to CreateApp
	app := server.CreateApp(dbConnection)

	// Start the application
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
