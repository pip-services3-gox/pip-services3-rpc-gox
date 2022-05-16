package test_services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusRestService(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", StatusRestServicePort)
	// Test "Status"
	getRes, getErr := http.Get(url + "/status")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
