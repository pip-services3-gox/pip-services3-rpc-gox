package test_clients

import (
	"testing"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	tlogic "github.com/pip-services3-go/pip-services3-rpc-go/test/logic"
	"github.com/stretchr/testify/assert"
)

func TestDummyDirectClient(t *testing.T) {

	var _dummy1 tdata.Dummy
	var _dummy2 tdata.Dummy

	var client *DummyDirectClient

	ctrl := tlogic.NewDummyController()
	client = NewDummyDirectClient()
	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
	)
	client.SetReferences(references)
	client.Open("")
	defer client.Close("")

	_dummy1 = tdata.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 = tdata.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	var dummy1 tdata.Dummy

	// Create one dummy
	dummy, err := client.CreateDummy("", _dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)
	dummy1 = *dummy

	// Create another dummy
	dummy, err = client.CreateDummy("", _dummy2)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)

	// Get all dummies
	dummies, err := client.GetDummies("", cdata.NewEmptyFilterParams(), cdata.NewPagingParams(0, 5, false))
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	dummy, err = client.UpdateDummy("", dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)
	dummy1 = *dummy

	// Delete dummy
	dummy, err = client.DeleteDummy("", dummy1.Id)
	assert.Nil(t, err)

	// Try to get delete dummy
	dummy, err = client.GetDummyById("", dummy1.Id)
	assert.Nil(t, err)
	assert.Nil(t, dummy)
}
