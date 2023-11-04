package api

import (
	"bytes"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"testing"
)

var app = server.CreateApp(lo.Must(db.InitDB(true)))

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(`{"email": "arnav@mail.com", "password": "123456"}`)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp := lo.Must(app.Test(req))

	assert.Equal(t, 201, resp.StatusCode)

}
