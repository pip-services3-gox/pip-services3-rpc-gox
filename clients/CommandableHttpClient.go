package clients

import (
	"context"
	"reflect"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

// CommandableHttpClient is abstract client that calls commandable HTTP service.
//
// Commandable services are generated automatically for ICommandable objects.
// Each command is exposed as POST operation that receives all parameters
// in body object.
//
//	Configuration parameters:
//		- base_route:              base route for remote URI
//		- connection(s):
//			- discovery_key:         (optional) a key to retrieve the connection from connect.idiscovery.html IDiscovery]]
//			- protocol:              connection protocol: http or https
//			- host:                  host name or IP address
//			- port:                  port number
//			- uri:                   resource URI or connection string with all parameters in it
//		- options:
//			- retries:               number of retries (default: 3)
//			- connect_timeout:       connection timeout in milliseconds (default: 10 sec)
//			- timeout:               invocation timeout in milliseconds (default: 10 sec)
//
//	References:
//		- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0       (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0      (optional) IDiscovery services to resolve connection
//
//	Example:
//		type MyCommandableHttpClient struct{
//			*CommandableHttpClient
//			prototype reflect.Type // type of operation data
//			...
//		}
//		func (c * MyCommandableHttpClient) GetData(correlationId string, id string)(result MyData, err error){
//			params:= cdata.NewEmptyStringValueMap()
//			params.Set("id",id)
//			res, err := c.CallCommand(
//				prototype
//				"get_data",
//				correlationId,
//				params,
//			)
//			...
//			// convert response to MyData
//			...
//			return result, err
//		}
//
//
//		client = NewMyCommandableHttpClient();
//		client.Configure(cconf.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080
//		));
//
//		result, err := client.GetData("123", "1")
//		...
type CommandableHttpClient struct {
	*RestClient
}

// NewCommandableHttpClient is creates a new instance of the client.
//	Parameters:
//		- baseRoute string a base route for remote service.
//	Returns: *CommandableHttpClient pointer on new instance
func NewCommandableHttpClient(baseRoute string) *CommandableHttpClient {
	c := CommandableHttpClient{}
	c.RestClient = NewRestClient()
	c.BaseRoute = baseRoute
	return &c
}

// CallCommand is calls a remote method via HTTP commadable protocol.
// The call is made via POST operation and all parameters are sent in body object.
// The complete route to remote method is defined as baseRoute + "/" + name.
//	Parameters:
//		- ctx context.Context
//		- prototype reflect.Type type of returned data
//		- name        string      a name of the command to call.
//		- correlationId  string   (optional) transaction id to trace execution through call chain.
//		- params     cdata.StringValueMap       command parameters.
//	Returns: result any, err error result or error.
func (c *CommandableHttpClient) CallCommand(ctx context.Context, prototype reflect.Type, name string, correlationId string, params *cdata.AnyValueMap) (result any, err error) {
	timing := c.Instrument(ctx, correlationId, c.BaseRoute+"."+name)
	cRes, cErr := c.Call(ctx, prototype, "post", name, correlationId, nil, params.Value())
	timing.EndTiming(ctx, cErr)
	return cRes, cErr
}
