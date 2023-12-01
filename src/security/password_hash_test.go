package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	testPass := "kqj35h6l3jh6"
	hashedPass := HashPassword(testPass)
	assert.NotNil(t, hashedPass)
}

func TestCheckPasswordHash(t *testing.T) {
	testPass := "kqj35h6l3jh6"
	hashedPass := HashPassword(testPass)
	assert.NotNil(t, hashedPass)
	assert.True(t, CheckPasswordHash(testPass, hashedPass))
	assert.False(t, CheckPasswordHash(testPass, "wrong hash"))
	assert.False(t, CheckPasswordHash("wrong pass", hashedPass))
	assert.False(t, CheckPasswordHash("", ""))
}
