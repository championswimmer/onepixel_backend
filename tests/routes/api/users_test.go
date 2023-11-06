package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/db"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/models"
	"onepixel_backend/src/server"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var app = server.CreateApp(lo.Must(db.InitDBTest()))
var USER_EMAIL = "test@mail.com"
var USER_PASSWORD = "test1234"

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, USER_EMAIL, USER_PASSWORD))

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(app.Test(req))

	assert.Equal(t, 201, resp.StatusCode)
}

func TestUsersRoute_LoginUser(t *testing.T) {
	var responseStruct dtos.LoginResponse
	reqBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, USER_EMAIL, USER_PASSWORD))
	req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	err = json.Unmarshal(body, &responseStruct)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, responseStruct.Token, "")

	user, err := auth.ValidateJWT(responseStruct.Token)
	assert.Equal(t, err, nil)
	assert.Equal(t, user.ID, uint(1))
}

func TestUsersRoute_RegisterUserDuplicateFail(t *testing.T) {
	reqBody := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, USER_EMAIL, USER_PASSWORD))
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 201, resp.StatusCode)

	resp = lo.Must(app.Test(req))
	assert.Equal(t, 409, resp.StatusCode)
}

func TestUsersRoute_GetUserInfoUnauthorized(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUsersRoute_GetUserInfoUnauthorizedInvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	req.Header.Set("Authorization", "xxxxxxxx")
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUsersRoute_GetUserInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	jwt := auth.CreateJWTFromUser(&models.User{ID: 1})
	req.Header.Set("Authorization", jwt)
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 200, resp.StatusCode)
}
