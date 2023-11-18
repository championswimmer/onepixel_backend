package utils

import (
	"onepixel_backend/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashPassword(t *testing.T) {

	hashedPassword, err := utils.HashPassword("test")
	assert.Nil(t, err)
	assert.NotNil(t, hashedPassword)

	// If passwords are not same then error will be returned
	error := utils.VerifyPassword(hashedPassword, "test")

	assert.Nil(t, error)
}
