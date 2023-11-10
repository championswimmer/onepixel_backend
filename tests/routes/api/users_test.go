package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/db"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/models"
	"onepixel_backend/src/security"
	"onepixel_backend/src/server"
	"testing"

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

	var responseBody dtos.ErrorResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, 409, resp.StatusCode)
	assert.Equal(t, 409, responseBody.Status)
	assert.Equal(t, "User with this email already exists", responseBody.Message)
}

func TestUsersRoute_RegisterUserBodyParsingFail(t *testing.T) {
	reqBody := []byte(`{"email": "user1461134@test.com", "password": "123456"}`)

	// Not setting any content-type will generate a Body Parsing error
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))

	resp := lo.Must(app.Test(req))

	var responseBody dtos.ErrorResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, 400, responseBody.Status)
	assert.Equal(t, "The request body is not valid", responseBody.Message)
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
	jwt := security.CreateJWTFromUser(&models.User{ID: 1})
	req.Header.Set("Authorization", jwt)
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 200, resp.StatusCode)
}
func TestUsersRoute_ShouldNotRegisterUserWhenNoPassword(t *testing.T) {
	reqBody := []byte(`{"email": "arnav@mail.com"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 422, resp.StatusCode)
}
func TestUsersRoute_ShouldNotRegisterUserWhenNoEmail(t *testing.T) {
	reqBody := []byte(`{"password": "12345"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp := lo.Must(app.Test(req))
	assert.Equal(t, 422, resp.StatusCode)
}
