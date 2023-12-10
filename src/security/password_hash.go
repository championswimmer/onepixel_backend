package security

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
)

const HashCostFactor = 10

func HashPassword(password string) string {
	if password == "" {
		log.Error(logger.RedBold, "Hashing empty password", logger.Reset)
	}
	hashedPassword := lo.Must(bcrypt.GenerateFromPassword([]byte(password), HashCostFactor))

	return string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		log.Error(logger.RedBold, "Comparing empty password", logger.Reset)
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
