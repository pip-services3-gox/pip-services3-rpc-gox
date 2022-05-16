package test_clients

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
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
