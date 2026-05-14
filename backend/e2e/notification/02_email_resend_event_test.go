package notification

import (
	"encoding/json"
	"net/http"
	"testing"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

func TestEmailResendEventSuccess(t *testing.T) {
	client := helpers.NewClient()

	to := "rezaoda@gmail.com"
	appName := "SISKOR" // must be registered in GetApiResendKey
	provider := "Resend"

	url := helpers.BaseURL() +
		"/check-mail?to=" + to +
		"&appName=" + appName +
		"&provider=" + provider + "&type=event"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var apiRes helpers.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiRes)
	assert.NoError(t, err)

	assert.True(t, apiRes.Status)
	assert.Equal(t, "Send Mail On Event Success", apiRes.Message)
}
