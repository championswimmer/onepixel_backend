package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"onepixel_backend/src/db"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/server"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var app = server.CreateApp(lo.Must(db.InitDB(true)))

func TestUsersRoute_RegisterUser(t *testing.T) {
	reqBody := []byte(`{"email": "arnav@mail.com", "password": "123456"}`)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(app.Test(req))

	assert.Equal(t, 201, resp.StatusCode)
}

func TestUsersRoute_LoginUser(t *testing.T) {
	var responseStruct dtos.LoginResponse

	reqBody := []byte(`{"email": "arnav@mail.com", "password": "123456"}`)
	req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	err = json.Unmarshal(body, &responseStruct)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, responseStruct.Token, "")

	// Validate the JWT token
	token, err := jwt.Parse(responseStruct.Token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		assert.Equal(t, ok, true)
		return []byte("some_random_secret_string"), nil
	})
	assert.Equal(t, err, nil)
	assert.NotEqual(t, token, "")
}
