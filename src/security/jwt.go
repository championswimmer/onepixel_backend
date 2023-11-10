package security

import (
	"onepixel_backend/src/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

// TODO: pick JWT_KEY from config

var JWT_KEY = "this is a sample key, should change in prod"
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

func CreateJWTFromUser(u *models.User) string {
	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, jwt.MapClaims{
		"sub": strconv.Itoa(int(u.ID)),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
	})

	return lo.Must(token.SignedString([]byte(JWT_KEY)))
}

func ValidateJWT(t string) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	sub := uint(lo.Must(strconv.Atoi(lo.Must(claims.GetSubject()))))
	exp := lo.Must(claims.GetExpirationTime())
	if exp.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return &models.User{ID: sub}, nil
}
