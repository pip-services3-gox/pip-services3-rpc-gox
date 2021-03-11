package test_rpc_clients

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	testrpc "github.com/pip-services3-gox/pip-services3-rpc-gox/test"
	testservices "github.com/pip-services3-gox/pip-services3-rpc-gox/test/services"
)

func TestDummyCommandableHttpClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var service *testservices.DummyCommandableHttpService
	var client *DummyCommandableHttpClient

	var fixture *DummyClientFixture

	ctrl := testrpc.NewDummyController()

	service = testservices.NewDummyCommandableHttpService()
	service.Configure(restConfig)

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
	)
	service.SetReferences(references)

	service.Open("")
	defer service.Close("")

	client = NewDummyCommandableHttpClient()
	fixture = NewDummyClientFixture(client)

	client.Configure(restConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
