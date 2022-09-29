package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
	rpccon "github.com/pip-services3-gox/pip-services3-rpc-gox/connect"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// RestClient is abstract client that calls remove endpoints using HTTP/REST protocol.
//
//	Configuration parameters:
//		- base_route:              base route for remote URI
//		- connection(s):
//			- discovery_key:         (optional) a key to retrieve the connection from IDiscovery
//			- protocol:              connection protocol: http or https
//			- host:                  host name or IP address
//			- port:                  port number
//			- uri:                   resource URI or connection string with all parameters in it
//		- options:
//			- retries:               number of retries (default: 3)
//			- connect_timeout:        connection timeout in milliseconds (default: 10 sec)
//			- timeout:               invocation timeout in milliseconds (default: 10 sec)
//			- correlation_id 	 place for adding correalationId, query - in query string, headers - in headers, both - in query and headers (default: query)
//
//	References:
//		- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0         (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0        (optional)  IDiscovery services to resolve connection
//
//	see services.RestService
//	see services.CommandableHttpService
//
//	Example:
//		type MyRestClient struct {
//			*RestClient
//		}
//		...
//		func (c *MyRestClient) GetData(ctx context.Context, correlationId string, id string) (result *tdata.MyDataPage[MyData], err error) {
//			timind := c.Instrument(ctx, correlationId, "myData.get_page_by_filter")
// 			defer timing.EndTiming(ctx)
//
//			params := cdata.NewEmptyStringValueMap()
//			params.Set("id", id)
//			response, calErr := c.Call(MyDataPageType, "get", "/data", correlationId, params, nil)
//			if calErr != nil {
//				return nil, calErr
//			}
//
//			return return clients.HandleHttpResponse[*tdata.MyDataPage[MyData]](response, correlationId)
//		}
//
//		client := NewMyRestClient();
//		client.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		));
//
//		result, err := client.GetData(context.Background(), "123", "1")
//		...
type RestClient struct {
	defaultConfig *cconf.ConfigParams
	//The HTTP client.
	Client *http.Client
	//The connection resolver.
	ConnectionResolver rpccon.HttpConnectionResolver
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
	//The configuration options.
	Options *cconf.ConfigParams
	//The base route.
	BaseRoute string
	//The number of retries.
	Retries int
	//The default headers to be added to every request.
	Headers *cdata.StringValueMap
	//The connection timeout in milliseconds.
	ConnectTimeout int
	//The invocation timeout in milliseconds.
	Timeout int
	//The remote service uri which is calculated on open.
	Uri string
	// add correlation id to headers
	passCorrelationId string
}

const (
	DefaultRequestMaxSize = 1024 * 1024
	DefaultConnectTimeout = 10000
	DefaultTimeout        = 10000
	DefaultRetriesCount   = 3
)

// NewRestClient creates new instance of RestClient
//	Returns: pointer on NewRestClient
func NewRestClient() *RestClient {
	rc := RestClient{}
	rc.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"options.request_max_size", DefaultRequestMaxSize,
		"options.connectTimeout", DefaultConnectTimeout,
		"options.timeout", DefaultTimeout,
		"options.retries", DefaultRetriesCount,
		"options.debug", true,
		"options.correlation_id", "query",
	)
	rc.ConnectionResolver = *rpccon.NewHttpConnectionResolver()
	rc.Logger = clog.NewCompositeLogger()
	rc.Counters = ccount.NewCompositeCounters()
	rc.Tracer = ctrace.NewCompositeTracer()
	rc.Options = cconf.NewEmptyConfigParams()
	rc.Retries = 1
	rc.Headers = cdata.NewEmptyStringValueMap()
	rc.ConnectTimeout = 10000
	rc.passCorrelationId = "query"
	return &rc
}

// Configure component by passing configuration parameters.
//	Parameters:
//		- ctx context.Context
//		- config *cconf.ConfigParams   configuration parameters to be set.
func (c *RestClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)
	c.Options = c.Options.Override(config.GetSection("options"))

	c.Retries = config.GetAsIntegerWithDefault("options.retries", c.Retries)
	c.ConnectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.ConnectTimeout)
	c.Timeout = config.GetAsIntegerWithDefault("options.timeout", c.Timeout)

	c.BaseRoute = config.GetAsStringWithDefault("base_route", c.BaseRoute)

	c.passCorrelationId = config.GetAsStringWithDefault("options.correlation_id_place", c.passCorrelationId)
	c.passCorrelationId = config.GetAsStringWithDefault("options.correlation_id", c.passCorrelationId)
}

// SetReferences to dependent components.
//	Parameters:
//		- ctx context.Context
//		- references  crefer.IReferences	references to locate the component dependencies.
func (c *RestClient) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a services.InstrumentTiming object that is used to end the time measurement.
//	Parameters:
//		- ctx context.Context
//		- correlationId string (optional) transaction id to trace execution through call chain.
//		- name string a method name.
//	Returns: services.InstrumentTiming object to end the time measurement.
func (c *RestClient) Instrument(ctx context.Context, correlationId string, name string) *services.InstrumentTiming {
	c.Logger.Trace(ctx, correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, correlationId, name, "")
	return services.NewInstrumentTiming(correlationId, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// InstrumentError method are dds instrumentation to error handling.
//	Parameters:
//		- ctx context.Context
//		- correlationId string  (optional) transaction id to trace execution through call chain.
//		- name   string         a method name.
//		- err    error          an occured error
//		- result  any           (optional) an execution result
//	Returns: result any, err error an execution result and error
func (c *RestClient) InstrumentError(ctx context.Context, correlationId string, name string, inErr error, inRes any) (result any, err error) {
	if inErr != nil {
		c.Logger.Error(ctx, correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(ctx, name+".call_errors")
	}

	return inRes, inErr
}

// IsOpen are checks if the component is opened.
//	Returns: true if the component has been opened and false otherwise.
func (c *RestClient) IsOpen() bool {
	return c.Client != nil
}

// Open method are opens the component.
//	Parameters:
//		- ctx context.Context
//		- correlationId string	(optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *RestClient) Open(ctx context.Context, correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	connection, _, err := c.ConnectionResolver.Resolve(correlationId)
	if err != nil {
		return err
	}

	c.Uri = connection.Uri()
	c.Client = &http.Client{
		Timeout: time.Duration(c.Timeout+c.ConnectTimeout) * time.Millisecond,
	}
	if c.Client == nil {
		return cerr.NewConnectionError(
			correlationId,
			"CANNOT_CONNECT",
			"Connection to REST service failed",
		).WithDetails("url", c.Uri)
	}

	return nil
}

// Close method are closes component and frees used resources.
//	Parameters:
//		- ctx context.Context
//		- correlationId  string	(optional) transaction id to trace execution through call chain.
// Returns: error or nil no errors occured.
func (c *RestClient) Close(ctx context.Context, correlationId string) error {
	if c.Client != nil {
		c.Logger.Debug(ctx, correlationId, "Closed REST service at %s", c.Uri)
		c.Client = nil
		c.Uri = ""
	}
	return nil
}

// AddCorrelationId method are adds a correlation id (correlation_id) to invocation parameter map.
//	Parameters:
//		- params    *cdata.StringValueMap        invocation parameters.
//		- correlationId  string  (optional) a correlation id to be added.
//	Returns: invocation parameters with added correlation id.
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
//	Parameters:
//		- params  *cdata.StringValueMap      invocation parameters.
//		- filter  *cdata.FilterParams     (optional) filter parameters
//	Returns: invocation parameters with added filter parameters.
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
		if paging.Skip >= 0 {
			params.Put("skip", paging.Skip)
		}
		if paging.Take >= 0 {
			params.Put("take", paging.Take)
		}
	}

	return params
}

// Call method are calls a remote method via HTTP/REST protocol.
//	Parameters:
//		- ctx context.Context
//		- prototype reflect.Type type for convert JSON result. Set nil for return raw JSON string
//		- method 	string           HTTP method: "get", "head", "post", "put", "delete"
//		- route   string          a command route. Base route will be added to this route
//		- correlationId  string    (optional) transaction id to trace execution through call chain.
//		- params  cdata.StringValueMap          (optional) query parameters.
//		- data   any           (optional) body object.
//	Returns: result any, err error result object or error.
func (c *RestClient) Call(ctx context.Context, method string, route string, correlationId string,
	params *cdata.StringValueMap, data any) (*http.Response, error) {

	method = strings.ToUpper(method)

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if c.passCorrelationId == "query" || c.passCorrelationId == "both" {
		params = c.AddCorrelationId(params, correlationId)
	}

	url := c.buildURL(route, params)

	if !c.IsOpen() {
		return nil, nil
	}

	var jsonStr string
	if data != nil {
		jsonStr, _ = convert.JsonConverter.ToJson(data)
	}

	retries := c.Retries
	var response *http.Response

	for retries > 0 {
		req, err := c.prepareRequest(ctx, correlationId, method, url, []byte(jsonStr))
		if err != nil {
			return nil, err
		}

		response, err = c.Client.Do(req)
		if err != nil {
			retries--
			if retries == 0 {
				return nil, cerr.NewUnknownError(
					correlationId,
					"COMMUNICATION_ERROR",
					"Unknown communication problem on REST client",
				).
					WithCause(err)
			}

			err = c.waitForRetry(ctx, correlationId, retries)
			if err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	if response.StatusCode == 204 {
		_ = response.Body.Close()
		return nil, nil
	}

	if response.StatusCode >= 400 {
		defer response.Body.Close()
		return nil, c.handleResponseError(response, correlationId)
	}

	return response, nil
}

func (c *RestClient) waitForRetry(ctx context.Context, correlationId string, retries int) error {
	waitTime := c.Timeout * int(math.Pow(float64(c.Retries-retries), 2))

	select {
	case <-time.After(time.Duration(waitTime) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return cerr.ApplicationErrorFactory.Create(
			&cerr.ErrorDescription{
				Type:          "Application",
				Category:      "Application",
				Code:          "CONTEXT_CANCELLED",
				Message:       "request canceled by parent context",
				CorrelationId: correlationId,
			},
		)
	}
}

func (c *RestClient) prepareRequest(ctx context.Context, correlationId string,
	method string, url string, body []byte) (*http.Request, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, cerr.NewUnknownError(
			correlationId,
			"UNSUPPORTED_METHOD",
			"Method is not supported by REST client",
		).
			WithDetails("verb", method).
			WithCause(err)
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.passCorrelationId == "headers" || c.passCorrelationId == "both" {
		req.Header.Set("correlation_id", correlationId)
	}
	for k, v := range c.Headers.Value() {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *RestClient) handleResponseError(response *http.Response, correlationId string) error {
	r, rErr := ioutil.ReadAll(response.Body)
	if rErr != nil {
		eDesct := cerr.ErrorDescription{
			Type:          "Application",
			Category:      "Application",
			Status:        response.StatusCode,
			Code:          "",
			Message:       rErr.Error(),
			CorrelationId: correlationId,
		}
		return cerr.ApplicationErrorFactory.Create(&eDesct).WithCause(rErr)
	}

	appErr := cerr.ApplicationError{}
	_ = json.Unmarshal(r, &appErr)
	if appErr.Status == 0 && len(r) > 0 { // not standart Pip.Services error
		values := make(map[string]any)
		decodeErr := json.Unmarshal(r, &values)
		if decodeErr != nil { // not json response
			appErr.Message = (string)(r)
		}
		appErr.Details = values
	}
	appErr.Status = response.StatusCode
	return &appErr
}

func (c *RestClient) buildURL(route string, params *cdata.StringValueMap) string {
	route = c.createRequestRoute(route)
	route = c.putParamsToRequestRoute(route, params)
	return c.Uri + route
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

func (c *RestClient) putParamsToRequestRoute(route string, params *cdata.StringValueMap) string {
	if params.Len() > 0 {
		var builder strings.Builder
		builder.Grow(1024)
		builder.WriteString(route)
		builder.WriteString("?")
		for k, v := range params.Value() {
			builder.WriteString(neturl.QueryEscape(k))
			builder.WriteString("=")
			builder.WriteString(neturl.QueryEscape(v))
			builder.WriteString("&")
		}
		route = builder.String()
		if strings.HasSuffix(route, "&") {
			route = strings.TrimRight(route, "&")
		}
	}

	return route
}
