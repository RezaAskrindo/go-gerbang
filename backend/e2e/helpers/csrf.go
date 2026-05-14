package helpers

import (
	"encoding/json"
	"net/http"
)

func GetCSRFToken(client *http.Client) (string, error) {
	req, _ := http.NewRequest(
		http.MethodGet,
		BaseURL()+"/secure-gateway-c",
		nil,
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Data.(string), nil
}
