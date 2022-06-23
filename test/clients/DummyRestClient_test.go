package test_clients

import (
	"context"
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

	client.Configure(context.TODO(), restConfig)
	client.SetReferences(context.TODO(), cref.NewEmptyReferences())
	client.Open(context.TODO(), "")

	t.Run("DummyRestClient.CrudOperations", fixture.TestCrudOperations)
}
