package api

import (
	"bytes"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"onepixel_backend/src/server"
	"testing"
)

var app = server.CreateApp()

func TestUsersRoute_GetAllUsers(t *testing.T) {

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp := lo.Must(app.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

}

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(`{"email": "arnav@mail.com", "password": "123456"}`)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(app.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

}
