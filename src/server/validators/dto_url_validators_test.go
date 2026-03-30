package validators

import (
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/dtos"
	"testing"
)

func TestValidateCreateUrlGroupDtoRequest(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		err := ValidateCreateUrlGroupDtoRequest(&dtos.CreateUrlGroupRequest{
			ShortPath: "grp123",
			CreatorID: 1,
		})
		assert.Nil(t, err)
	})

	t.Run("MissingCreatorID", func(t *testing.T) {
		err := ValidateCreateUrlGroupDtoRequest(&dtos.CreateUrlGroupRequest{
			ShortPath: "grp123",
		})
		assert.Equal(t, mandatoryUrlGroupCreatorFieldError, err)
	})

	t.Run("InvalidGroupToken", func(t *testing.T) {
		err := ValidateCreateUrlGroupDtoRequest(&dtos.CreateUrlGroupRequest{
			ShortPath: "bad!",
			CreatorID: 1,
		})
		assert.Equal(t, invalidCreateUrlGroupError, err)
	})
}

func TestValidateCreateShortCodeRequest(t *testing.T) {
	t.Run("AcceptsRadix64", func(t *testing.T) {
		assert.Nil(t, ValidateCreateShortCodeRequest("Abc_-09"))
	})

	t.Run("RejectsEmpty", func(t *testing.T) {
		assert.Equal(t, invalidCreateShortCodeError, ValidateCreateShortCodeRequest(""))
	})

	t.Run("RejectsTooLong", func(t *testing.T) {
		assert.Equal(t, invalidCreateShortCodeError, ValidateCreateShortCodeRequest("abcdefghijk"))
	})

	t.Run("RejectsInvalidChars", func(t *testing.T) {
		assert.Equal(t, invalidCreateShortCodeError, ValidateCreateShortCodeRequest("bad!"))
	})
}
