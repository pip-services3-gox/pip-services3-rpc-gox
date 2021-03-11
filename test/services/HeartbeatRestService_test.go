package test_rpc_services

import (
	"net/http"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
	"github.com/stretchr/testify/assert"
)

func TestHeartbeatRestService(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var service *services.HeartbeatRestService

	service = services.NewHeartbeatRestService()
	service.Configure(restConfig)

	service.Open("")
	defer service.Close("")
	url := "http://localhost:3000"
	// Test "Heartbeat"
	getRes, getErr := http.Get(url + "/heartbeat")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
