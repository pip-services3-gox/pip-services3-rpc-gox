package test_clients

import (
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

func TestDummyRestClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestServicePort,
		"options.correlation_id_place", "headers",
	)

	client := NewDummyRestClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(restConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")

	t.Run("DummyRestClient.CrudOperations", fixture.TestCrudOperations)
}
