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
