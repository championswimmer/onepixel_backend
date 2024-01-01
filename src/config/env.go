package config

import "github.com/joho/godotenv"

func init() {
	godotenv.Load("onepixel.local.env")
}
