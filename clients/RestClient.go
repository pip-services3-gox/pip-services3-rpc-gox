package clients

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	rpccon "github.com/pip-services3-gox/pip-services3-rpc-gox/connect"
)

/*
RestClient is abstract client that calls remove endpoints using HTTP/REST protocol.

Configuration parameters:
  - base_route:              base route for remote URI
  - connection(s):
    - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
    - protocol:              connection protocol: http or https
    - host:                  host name or IP address
    - port:                  port number
    - uri:                   resource URI or connection string with all parameters in it
  - options:
    - retries:               number of retries (default: 3)
    - connectTimeout:        connection timeout in milliseconds (default: 10 sec)
	- timeout:               invocation timeout in milliseconds (default: 10 sec)
	- correlation_id_place 	 place for adding correalationId, query - in query string, headers - in headers, both - in query and headers (default: query)

References:
- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
- *:counters:*:*:1.0         (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0        (optional)  IDiscovery services to resolve connection

See RestService
See CommandableHttpService

 Example:
    type MyRestClient struct {
		*RestClient
	}
    ...
    func (c *MyRestClient) GetData(correlationId string, id string) (result *testrpc.MyDataPage, err error) {

		params := cdata.NewEmptyStringValueMap()
		params.Set("id", id)

		calValue, calErr := c.Call(MyDataPageType, "get", "/data", correlationId, params, nil)
		if calErr != nil {
			return nil, calErr
		}

		result, _ = calValue.(*testrpc.MyDataPage)
		c.Instrument(correlationId, "myData.get_page_by_filter")
		return result, nil
	}

    client := NewMyRestClient();
    client.Configure(cconf.NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));

	result, err := client.GetData("123", "1")
	 ...

*/
type RestClient struct {
	defaultConfig cconf.ConfigParams
	//The HTTP client.
	Client *http.Client
	//The connection resolver.
	ConnectionResolver rpccon.HttpConnectionResolver
	//The logger.
	Logger clog.CompositeLogger
	//The performance counters.
	Counters ccount.CompositeCounters
	//The configuration options.
	Options cconf.ConfigParams
	//The base route.
	BaseRoute string
	//The number of retries.
	Retries int
	//The default headers to be added to every request.
	Headers cdata.StringValueMap
	//The connection timeout in milliseconds.
	ConnectTimeout int
	//The invocation timeout in milliseconds.
	Timeout int
	//The remote service uri which is calculated on open.
	Uri string
	// add correlation id to headers
	correlationIdPlace string
}

// NewRestClient creates new instance of RestClient
// Retruns pointer on NewRestClient
func NewRestClient() *RestClient {
	rc := RestClient{}
	rc.defaultConfig = *cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"options.request_max_size", 1024*1024,
		"options.connectTimeout", 10000,
		"options.timeout", 10000,
		"options.retries", 3,
		"options.debug", true,
		"options.correlation_id_place", "query",
	)
	rc.ConnectionResolver = *rpccon.NewHttpConnectionResolver()
	rc.Logger = *clog.NewCompositeLogger()
	rc.Counters = *ccount.NewCompositeCounters()
	rc.Options = *cconf.NewEmptyConfigParams()
	rc.Retries = 1
	rc.Headers = *cdata.NewEmptyStringValueMap()
	rc.ConnectTimeout = 10000
	rc.correlationIdPlace = "query"
	return &rc
}

// Configures component by passing configuration parameters.
// Parameters:
// - config *cconf.ConfigParams   configuration parameters to be set.
func (c *RestClient) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(&c.defaultConfig)
	c.ConnectionResolver.Configure(config)
	c.Options = *c.Options.Override(config.GetSection("options"))
	c.Retries = config.GetAsIntegerWithDefault("options.retries", c.Retries)
	c.ConnectTimeout = config.GetAsIntegerWithDefault("options.connectTimeout", c.ConnectTimeout)
	c.Timeout = config.GetAsIntegerWithDefault("options.timeout", c.Timeout)
	c.BaseRoute = config.GetAsStringWithDefault("base_route", c.BaseRoute)
	c.correlationIdPlace = config.GetAsStringWithDefault("options.correlation_id_place", c.correlationIdPlace)
}

// Sets references to dependent components.
// Parameters:
// - references  crefer.IReferences	references to locate the component dependencies.
func (c *RestClient) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.ConnectionResolver.SetReferences(references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
// Parameters:
// - correlationId  string   (optional) transaction id to trace execution through call chain.
// - name    string          a method name.
// Return Timing object to end the time measurement.
func (c *RestClient) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	return c.Counters.BeginTiming(name + ".call_time")
}

// InstrumentError method are dds instrumentation to error handling.
//    - correlationId   string  (optional) transaction id to trace execution through call chain.
//    - name   string           a method name.
//    - err    error           an occured error
//    - result  interface{}           (optional) an execution result
// Returns: result interface{}, err error
// an execution result and error
func (c *RestClient) InstrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
	if inErr != nil {
		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(name + ".call_errors")
	}

	return inRes, inErr
}

// IsOpen are checks if the component is opened.
// Return true if the component has been opened and false otherwise.
func (c *RestClient) IsOpen() bool {
	return c.Client != nil
}

// Open method are opens the component.
// 	- correlationId string	(optional) transaction id to trace execution through call chain.
// Returns: error
// error or nil no errors occured.
func (c *RestClient) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	connection, _, conErr := c.ConnectionResolver.Resolve(correlationId)
	if conErr != nil {
		return conErr
	}

	c.Uri = connection.Uri()
	localClient := http.Client{}
	localClient.Timeout = (time.Duration)(c.Timeout) * time.Millisecond
	c.Client = &localClient
	if c.Client == nil {
		ex := cerr.NewConnectionError(correlationId, "CANNOT_CONNECT", "Connection to REST service failed").WithDetails("url", c.Uri)
		return ex
	}

	return nil
}

// Close method are closes component and frees used resources.
// Parameters:
// 	- correlationId  string	(optional) transaction id to trace execution through call chain.
// Retruns: error
// error or nil no errors occured.
func (c *RestClient) Close(correlationId string) error {
	if c.Client != nil {
		c.Logger.Debug(correlationId, "Closed REST service at %s", c.Uri)
		c.Client = nil
		c.Uri = ""
	}
	return nil
}

// AddCorrelationId method are adds a correlation id (correlation_id) to invocation parameter map.
// Parameters:
//    - params    *cdata.StringValueMap        invocation parameters.
//    - correlationId  string  (optional) a correlation id to be added.
//  Return invocation parameters with added correlation id.
func (c *RestClient) AddCorrelationId(params *cdata.StringValueMap, correlationId string) *cdata.StringValueMap {
	// Automatically generate short ids for now
	if correlationId == "" {
		//correlationId = IdGenerator.NextShort()
		return params
	}

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	params.Put("correlation_id", correlationId)
	return params
}

// AddFilterParams method are adds filter parameters (with the same name as they defined)
// to invocation parameter map.
// Parameters:
//    - params  *cdata.StringValueMap      invocation parameters.
//    - filter  *cdata.FilterParams     (optional) filter parameters
// Return invocation parameters with added filter parameters.
func (c *RestClient) AddFilterParams(params *cdata.StringValueMap, filter *cdata.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
	}
	return params
}

// AddPagingParams method are adds paging parameters (skip, take, total) to invocation parameter map.
// Parameters:
//    - params        invocation parameters.
//    - paging        (optional) paging parameters
// Return invocation parameters with added paging parameters.
func (c *RestClient) AddPagingParams(params *cdata.StringValueMap, paging *cdata.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		if paging.Skip != nil {
			params.Put("skip", *paging.Skip)
		}
		if paging.Take != nil {
			params.Put("take", *paging.Take)
		}
	}

	return params
}

func (c *RestClient) createRequestRoute(route string) string {
	builder := ""

	if c.BaseRoute != "" && len(c.BaseRoute) > 0 {
		if c.BaseRoute[0] != "/"[0] {
			builder += "/"
		}
		builder += c.BaseRoute
	}

	if route != "" && route[0] != "/"[0] {
		builder += "/"
	}
	builder += route

	return builder
}

// Call method are calls a remote method via HTTP/REST protocol.
// Parameters:
// 	- prototype reflect.Type type for convert JSON result. Set nil for return raw JSON string
//  - method 	string           HTTP method: "get", "head", "post", "put", "delete"
//  - route   string          a command route. Base route will be added to this route
//  - correlationId  string    (optional) transaction id to trace execution through call chain.
//  - params  cdata.StringValueMap          (optional) query parameters.
//  - data   interface{}           (optional) body object.
// Returns:  result interface{}, err error
// result object or error.
func (c *RestClient) Call(prototype reflect.Type, method string, route string, correlationId string, params *cdata.StringValueMap,
	data interface{}) (result interface{}, err error) {

	method = strings.ToUpper(method)
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	route = c.createRequestRoute(route)
	if c.correlationIdPlace == "query" || c.correlationIdPlace == "both" {
		params = c.AddCorrelationId(params, correlationId)
	}
	if params.Len() > 0 {
		route += "?"
		for k, v := range params.Value() {
			route += (k + "=" + v + "&")
		}
		if strings.HasSuffix(route, "&") {
			route = strings.TrimRight(route, "&")
		}
	}

	url := c.Uri + route

	if !c.IsOpen() {
		return nil, nil
	}
	var jsonStr []byte
	if data != nil {
		jsonStr, _ = json.Marshal(data)
	} else {
		jsonStr = make([]byte, 0, 0)
	}
	req, reqErr := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))

	if reqErr != nil {
		err = cerr.NewUnknownError(correlationId, "UNSUPPORTED_METHOD", "Method is not supported by REST client").WithDetails("verb", method).WithCause(reqErr)
		return nil, err
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.correlationIdPlace == "headers" || c.correlationIdPlace == "both" {
		req.Header.Set("correlation_id", correlationId)
	}
	//req.Header.Set("User-Agent", c.UserAgent)
	for k, v := range c.Headers.Value() {
		req.Header.Set(k, v)
	}

	retries := c.Retries
	var resp *http.Response
	var respErr error

	for retries > 0 {
		// Try send request
		resp, respErr = c.Client.Do(req)
		if respErr != nil {

			retries--
			if retries == 0 {
				err = cerr.NewUnknownError(correlationId, "COMMUNICATION_ERROR", "Unknown communication problem on REST client").WithCause(respErr)
				return nil, err
			}
			continue
		}
		break
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == 204 {
		return nil, nil
	}

	r, rErr := ioutil.ReadAll(resp.Body)
	if rErr != nil {
		eDesct := cerr.ErrorDescription{
			Type:          "Application",
			Category:      "Application",
			Status:        resp.StatusCode,
			Code:          "",
			Message:       rErr.Error(),
			CorrelationId: correlationId,
		}
		return nil, cerr.ApplicationErrorFactory.Create(&eDesct).WithCause(rErr)
	}

	if resp.StatusCode >= 400 {
		appErr := cerr.ApplicationError{}
		json.Unmarshal(r, &appErr)
		return nil, &appErr
	}

	if prototype != nil {
		return ConvertComandResult(r, prototype)
	}

	return r, rErr

}
