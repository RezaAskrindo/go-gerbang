package auth

import (
	"net/http"
	"testing"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoginWrongPassword(t *testing.T) {
	client := helpers.NewClient()

	payload := map[string]string{
		"identity": "test_user",
		"password": "wrongpassword",
	}

	resp, apiRes, err := helpers.DoJSON(
		client,
		http.MethodPost,
		helpers.BaseURL()+"/api/v1/auth/login",
		payload,
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, apiRes.Status)
}
