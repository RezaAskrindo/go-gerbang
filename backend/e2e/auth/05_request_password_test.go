package auth

import (
	"net/http"
	"testing"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

func TestRequestResetPassword(t *testing.T) {
	client := helpers.NewClient()

	payload := map[string]string{
		"identity": "toshibareza@gmail.com",
	}

	resp, apiRes, err := helpers.DoJSON(
		client,
		http.MethodPost,
		helpers.BaseURL()+"/api/v1/auth/request-reset-password?baseUrl=http://localhost:3000&sender=GERBANG",
		payload,
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, apiRes.Status)
}
