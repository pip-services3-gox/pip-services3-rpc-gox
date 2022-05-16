package test_services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatRestService(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", HeartbeatRestServicePort)
	// Test "Heartbeat"
	getRes, getErr := http.Get(url + "/heartbeat")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
