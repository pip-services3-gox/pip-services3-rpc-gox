package test_rpc_clients

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	testrpc "github.com/pip-services3-gox/pip-services3-rpc-gox/test"
	testservices "github.com/pip-services3-gox/pip-services3-rpc-gox/test/services"
)

func TestDummyRestClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
		"options.correlation_id_place", "headers",
	)

	var service *testservices.DummyRestService
	var client *DummyRestClient

	var fixture *DummyClientFixture

	ctrl := testrpc.NewDummyController()

	service = testservices.NewDummyRestService()
	service.Configure(restConfig)

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(references)
	service.Open("")
	defer service.Close("")

	client = NewDummyRestClient()
	fixture = NewDummyClientFixture(client)

	client.Configure(restConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")

	t.Run("DummyRestClient.CrudOperations", fixture.TestCrudOperations)
}
