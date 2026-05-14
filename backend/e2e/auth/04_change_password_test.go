package auth

import (
	"net/http"
	"testing"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

func TestChangePassword(t *testing.T) {
	client := helpers.NewClient()

	payload := map[string]string{
		"identity": "test_user",
		"password": "NewPassword123!",
	}

	resp, apiRes, err := helpers.DoJSON(
		client,
		http.MethodPost,
		helpers.BaseURL()+"/api/v1/auth/change-password",
		payload,
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, apiRes.Status)
}
