package test_clients

import (
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

func TestDummyCommandableHttpClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpServicePort,
	)

	client := NewDummyCommandableHttpClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(restConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
