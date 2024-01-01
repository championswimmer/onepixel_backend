package security

import (
	"onepixel_backend/src/config"
	"onepixel_backend/src/db/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

var SigningKey = []byte(config.JwtSigningKey)
var SigningMethod = jwt.SigningMethodHS256
var KeyDuration = config.JwtDurationDays

func CreateJWTFromUser(u *models.User) string {
	token := jwt.NewWithClaims(SigningMethod, jwt.MapClaims{
		"sub": strconv.Itoa(int(u.ID)),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().AddDate(0, 0, KeyDuration)), // 7 days
	})

	return lo.Must(token.SignedString(SigningKey))
}

func ValidateJWT(t string) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	sub := uint64(lo.Must(strconv.Atoi(lo.Must(claims.GetSubject()))))
	exp := lo.Must(claims.GetExpirationTime())
	if exp.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return &models.User{ID: sub}, nil
}
