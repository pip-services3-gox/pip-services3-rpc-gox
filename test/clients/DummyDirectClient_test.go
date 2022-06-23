package test_clients

import (
	"context"
	"testing"

	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	tlogic "github.com/pip-services3-gox/pip-services3-rpc-gox/test/logic"
)

func TestDummyDirectClient(t *testing.T) {

	client := NewDummyDirectClient()
	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor(
			"pip-services-dummies", "controller", "default",
			"default", "1.0",
		), tlogic.NewDummyController(),
	)
	client.SetReferences(context.Background(), references)
	client.Open(context.Background(), "")
	defer client.Close(context.Background(), "")

	fixture := NewDummyClientFixture(client)
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
