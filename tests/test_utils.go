package tests

import (
	"bytes"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/config"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/server"
	"testing"
)

var App = server.CreateAdminApp()

func TestUtil_CreateUser(t *testing.T, email string, password string) dtos.UserResponse {

	reqBody := []byte(`{"email": "` + email + `", "password": "` + password + `"}`)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-API-Key", config.AdminApiKey)

	resp := lo.Must(App.Test(req))

	assert.Equal(t, 201, resp.StatusCode)

	var responseBody dtos.UserResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.NotNil(t, responseBody.Token)

	return responseBody

}
