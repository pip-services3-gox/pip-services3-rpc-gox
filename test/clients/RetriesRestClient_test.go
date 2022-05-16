package test_clients

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestRetriesRestClient(t *testing.T) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "12345",

		"options.retries", "2",
		"options.timeout", "100",
		"options.connect_timeout", "100",
	)

	var client *DummyRestClient

	client = NewDummyRestClient()

	client.Configure(restConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")

	res, err := client.GetDummyById("", "1")
	assert.NotNil(t, err)
	assert.Nil(t, res)

}
