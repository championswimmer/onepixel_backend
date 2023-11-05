package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/db"
	"onepixel_backend/src/models"
	"onepixel_backend/src/server"
	"strings"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var app = server.CreateApp(lo.Must(db.InitDBTest()))

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(`{"email": "user1461134@test.com", "password": "123456"}`)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(app.Test(req))

	assert.Equal(t, 201, resp.StatusCode)

}

func TestUsersRoute_RegisterUserDuplicateFail(t *testing.T) {
	reqBody := []byte(`{"email": "user14641522@test.com", "password": "123456"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 201, resp.StatusCode)

	resp = lo.Must(app.Test(req))
	assert.Equal(t, 409, resp.StatusCode)
}

func TestUsersController_CreateBadJSON(t *testing.T) {
    // Simulate a request with bad JSON
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader("{bad json"))
    req.Header.Set("Content-Type", "application/json")

    resp, _ := app.Test(req, -1)

    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

    var responseBody map[string]string
    json.NewDecoder(resp.Body).Decode(&responseBody)

    assert.Contains(t, responseBody["error"], "Cannot parse JSON")
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
