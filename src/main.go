package main

import (
	"fmt"
	"github.com/samber/lo"
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"os"
)

// main godoc
//
// @title						onepixel API
// @version					0.1
// @description				1px.li URL Shortner API
// @termsOfService				https://github.com/championswimmer/onepixel_backend
// @contact.name				Arnav Gupta
// @contact.email				dev@championswimmer.in
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
// @host						127.0.0.1:3000
// @BasePath					/api/v1
// @securityDefinitions.apiKey	BearerToken
// @in							header
// @name						Authorization
func main() {
	// Initialize the database
	db := lo.Must(db.InitDBProd())

	// Create the app
	app := server.CreateApp(db)

	httpPort, _ := lo.Coalesce(os.Getenv("PORT"), "3000")

	// TODO: move port to external YAML config
	log.Fatal(app.Listen(fmt.Sprintf(":%s", httpPort)))
}
