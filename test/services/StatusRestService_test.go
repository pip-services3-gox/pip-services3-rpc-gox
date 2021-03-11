package test_rpc_services

import (
	"net/http"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
	"github.com/stretchr/testify/assert"
)

func TestStatusRestService(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var service *services.StatusRestService

	service = services.NewStatusRestService()
	service.Configure(restConfig)

	contextInfo := cinfo.NewContextInfo()
	contextInfo.Name = "Test"
	contextInfo.Description = "This is a test container"

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "status-service", "http", "default", "1.0"), service,
	)
	service.SetReferences(references)
	service.Open("")
	defer service.Close("")
	url := "http://localhost:3000"
	// Test "Status"
	getRes, getErr := http.Get(url + "/status")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
