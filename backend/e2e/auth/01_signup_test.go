package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

type UserResponse struct {
	IdAccount string `json:"idAccount"`
}

func TestSignupSuccess(t *testing.T) {
	client := helpers.NewClient()

	username := "test_user"
	email := "toshibareza@gmail.com"

	req, _ := http.NewRequest(http.MethodGet, helpers.BaseURL()+"/users/by-identity/?username="+username, nil)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			Status  bool         `json:"status"`
			Message string       `json:"message"`
			Data    UserResponse `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		fmt.Printf("User already exists: %+v\n", result.Data)

		delResp, delApiRes, err := helpers.DoJSON(client, http.MethodDelete, helpers.BaseURL()+"/users/"+result.Data.IdAccount, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, delResp.StatusCode)
		assert.True(t, delApiRes.Status)
	}

	payload := map[string]interface{}{
		"identityNumber": "1112220101910001",
		"username":       username,
		"fullName":       "Test User",
		"email":          email,
		"password":       "Password123!",
		"phoneNumber":    "081234567890",
		"dateOfBirth":    time.Date(1991, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, apiRes, err := helpers.DoJSON(
		client,
		http.MethodPost,
		helpers.BaseURL()+"/api/v1/auth/sign-up?active=true",
		payload,
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, apiRes.Status)
}
