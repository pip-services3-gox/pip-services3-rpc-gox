package test_connect

import (
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/connect"
	"github.com/stretchr/testify/assert"
)

func TestHttpConnectionResolver(t *testing.T) {

	t.Run("HttpConnectionResolver.Resolve_URI", ResolveURI)
	t.Run("HttpConnectionResolver.Resolve_Parameters", ResolveParameters)
}

func ResolveURI(t *testing.T) {
	resolver := connect.NewHttpConnectionResolver()
	resolver.Configure(cconf.NewConfigParamsFromTuples(
		"connection.uri",
		"http://somewhere.com:777",
	))

	connection, _, _ := resolver.Resolve("")

	assert.Equal(t, "http", connection.Protocol())
	assert.Equal(t, "somewhere.com", connection.Host())
	assert.Equal(t, 777, connection.Port())
}

func ResolveParameters(t *testing.T) {
	resolver := connect.NewHttpConnectionResolver()
	resolver.Configure(cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "somewhere.com",
		"connection.port", "777",
	))

	connection, _, _ := resolver.Resolve("")
	assert.Equal(t, "http://somewhere.com:777", connection.Uri())

}
