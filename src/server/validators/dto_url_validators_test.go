package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateShortCode(t *testing.T) {
	tests := []struct {
		name      string
		shortcode string
		wantErr   bool
		errMsg    string
	}{
		{"valid shortcode", "abc123", false, ""},
		{"valid max length", "abcdefghij", false, ""},
		{"empty shortcode", "", true, "Invalid short code"},
		{"too long shortcode", "abcdefghijk", true, "shortcode must be at most 10 characters long"},
		{"way too long shortcode", "abcdefghijklmnop", true, "shortcode must be at most 10 characters long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShortCode(tt.shortcode)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMsg, err.message)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
