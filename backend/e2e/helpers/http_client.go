package helpers

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

func NewClient() *http.Client {
	jar, _ := cookiejar.New(nil)

	return &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar, // ⭐ REQUIRED (credentials: include)
	}
}

func BaseURL() string {
	// if v := os.Getenv("BASE_URL"); v != "" {
	// 	return v
	// }
	return "http://localhost:9000"
}
