package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/security"
	"onepixel_backend/src/utils/applogger"
	"onepixel_backend/tests"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUsersRoute_LoginUser(t *testing.T) {
	// ---- REGISTER USER
	responseBody := tests.TestUtil_CreateUser(t, "user51341@test.com", "123456")

	// ---- LOGIN USER

	reqBody := []byte(`{"email": "user51341@test.com" , "password": "123456"}`)

	req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(tests.App.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	assert.NotNil(t, *responseBody.Token)
	applogger.Info(*responseBody.Token)

}

func TestUsersRoute_GetUserInfoUnauthorized(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	resp := lo.Must(tests.App.Test(req))
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUsersRoute_GetUserInfoUnauthorizedInvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	req.Header.Set("Authorization", "xxxxxxxx")
	resp := lo.Must(tests.App.Test(req))
	assert.Equal(t, 401, resp.StatusCode)
}

func TestUsersRoute_GetUserInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	jwt := security.CreateJWTFromUser(&models.User{ID: 1})
	req.Header.Set("Authorization", jwt)
	resp := lo.Must(tests.App.Test(req))
	assert.Equal(t, 200, resp.StatusCode)
}
