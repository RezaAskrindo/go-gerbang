package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrCSRF = errors.New("CSRF validation failed")

func DoJSON(
	client *http.Client,
	method, url string,
	body any,
	headers map[string]string,
) (*http.Response, *APIResponse, error) {

	return doJSON(client, method, url, body, headers, true)
}

func doJSON(
	client *http.Client,
	method, url string,
	body any,
	headers map[string]string,
	allowRetry bool,
) (*http.Response, *APIResponse, error) {

	// 1️⃣ Ensure CSRF token exists
	csrfToken, err := GetCSRFToken(client)
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SGCsrf-Token", csrfToken)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	var apiRes APIResponse
	_ = json.Unmarshal(raw, &apiRes)

	// 2️⃣ Detect CSRF failure
	if !apiRes.Status && apiRes.Message == "CSRF validation failed" {
		if allowRetry {
			// retry once (like JS)
			return doJSON(client, method, url, body, headers, false)
		}
		return resp, &apiRes, ErrCSRF
	}

	return resp, &apiRes, nil
}
