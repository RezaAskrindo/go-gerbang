package auth

import (
	"net/http"
	"testing"

	"go-gerbang/e2e/helpers"

	"github.com/stretchr/testify/assert"
)

func TestCheckMigration(t *testing.T) {
	client := helpers.NewClient()

	req, _ := http.NewRequest(http.MethodGet, helpers.BaseURL()+"/check-migration", nil)
	_, err := client.Do(req)
	assert.NoError(t, err)
}
