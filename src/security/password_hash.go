package security

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
)

const HASH_COST_FACTOR = 10

func HashPassword(password string) string {
	if password == "" {
		log.Error(logger.RedBold, "Hashing empty password")
	}
	hashedPassword := lo.Must(bcrypt.GenerateFromPassword([]byte(password), HASH_COST_FACTOR))

	return string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		log.Error(logger.RedBold, "Comparing empty password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
