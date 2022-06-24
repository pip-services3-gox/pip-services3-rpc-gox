package test_clients

import (
	"context"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	"github.com/stretchr/testify/assert"
)

func TestRetriesRestClient(t *testing.T) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "12345",

		"options.retries", "4",
		"options.timeout", "100",
		"options.connect_timeout", "100",
	)

	var client *DummyRestClient

	client = NewDummyRestClient()

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	client.Open(context.Background(), "")

	res, err := client.GetDummyById(context.Background(), "", "1")
	assert.NotNil(t, err)
	assert.Equal(t, tdata.Dummy{}, res)

}
