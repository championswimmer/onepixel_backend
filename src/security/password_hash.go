package security

import (
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"onepixel_backend/src/utils"
)

const HashCostFactor = 10

func HashPassword(password string) string {
	if password == "" {
		utils.AppLogger.Error("Hashing empty password")
	}
	hashedPassword := lo.Must(bcrypt.GenerateFromPassword([]byte(password), HashCostFactor))

	return string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		utils.AppLogger.Error("Comparing empty password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
