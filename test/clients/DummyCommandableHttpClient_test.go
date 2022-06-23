package test_clients

import (
	"context"
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

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	_ = client.Open(context.Background(), "")
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
