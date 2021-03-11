package services

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
RestService Abstract service that receives remove calls via HTTP/REST protocol.

Configuration parameters:

  - base_route:              base route for remote URI
  - dependencies:
    - endpoint:              override for HTTP Endpoint dependency
    - controller:            override for Controller dependency
  - connection(s):
    - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
    - protocol:              connection protocol: http or https
    - host:                  host name or IP address
    - port:                  port number
    - uri:                   resource URI or connection string with all parameters in it
  - credential - the HTTPS credentials:
    - ssl_key_file:         the SSL private key in PEM
    - ssl_crt_file:         the SSL certificate in PEM
    - ssl_ca_file:          the certificate authorities (root cerfiticates) in PEM

References:

- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0            (optional) IDiscovery services to resolve connection
- *:endpoint:http:*:1.0          (optional) HttpEndpoint reference

See RestClient

Example:

    type MyRestService struct {
		*RestService
		controller IMyController
	}

	   ...
	func NewMyRestService() *MyRestService {
		c := MyRestService{}
		c.RestService = services.NewRestService()
		c.RestService.IRegisterable = &c
		c.numberOfCalls = 0
		c.DependencyResolver.Put("controller", crefer.NewDescriptor("mygroup", "controller", "*", "*", "1.0"))
		return &c
	}

    func (c * MyRestService) SetReferences(references IReferences) {
        c.RestService.SetReferences(references);
		resolv := c.DependencyResolver.GetRequired("controller");
		if resolv != nil {
			c.controller, _ = resolv.(IMyController)
		}
    }

	func (c *MyRestService) getOneById(res http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()
		vars := mux.Vars(req)

		mydataId := params.Get("mydata_id")
		if mydataId == "" {
			mydataId = vars["mydatay_id"]
		}
		result, err := c.controller.GetOneById(
			params.Get("correlation_id"),
			mydataId)
		c.SendResult(res, req, result, err)

	}
    func (c * MyRestService) Register() {

		c.RegisterRoute(
		"get", "get_mydata/{mydata_id}",
		 &cvalid.NewObjectSchema().
			WithRequiredProperty("mydata_id", cconv.String).Schema,
		c.getOneById)
           ...
    }


    service := NewMyRestService();
    service.Configure(cconf.NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));
    service.SetReferences(cref.NewReferencesFromTuples(
       cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller
    ));

	opnRes := service.Open("123")
	if opnErr == nil {
	   fmt.Println("The REST service is running on port 8080");
	}

*/
type RestService struct {
	IRegisterable
	defaultConfig *cconf.ConfigParams
	config        *cconf.ConfigParams
	references    crefer.IReferences
	localEndpoint bool
	opened        bool
	//The base route.
	BaseRoute string
	//The HTTP endpoint that exposes this service.
	Endpoint *HttpEndpoint
	//The dependency resolver.
	DependencyResolver *crefer.DependencyResolver
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters

	SwaggerService ISwaggerService
	SwaggerEnable  bool
	SwaggerRoute   string
}

// NewRestService is create new instance of RestService
func NewRestService() *RestService {
	rs := RestService{}
	rs.defaultConfig = cconf.NewConfigParamsFromTuples(
		"base_route", "",
		"dependencies.endpoint", "*:endpoint:http:*:1.0",
		"dependencies.swagger", "*:swagger-service:*:*:1.0",
	)
	rs.DependencyResolver = crefer.NewDependencyResolver()
	rs.DependencyResolver.Configure(rs.defaultConfig)
	rs.Logger = clog.NewCompositeLogger()
	rs.Counters = ccount.NewCompositeCounters()
	rs.SwaggerEnable = false
	rs.SwaggerRoute = "swagger"
	return &rs
}

// Configure method are configures component by passing configuration parameters.
// Parameters:
//   - config  *cconf.ConfigParams  configuration parameters to be set.
func (c *RestService) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config
	c.DependencyResolver.Configure(config)
	c.BaseRoute = config.GetAsStringWithDefault("base_route", c.BaseRoute)
	c.SwaggerEnable = config.GetAsBooleanWithDefault("swagger.enable", c.SwaggerEnable)
	c.SwaggerRoute = config.GetAsStringWithDefault("swagger.route", c.SwaggerRoute)
}

// SetReferences method are sets references to dependent components.
// Parameters:
//   - references crefer.IReferences	references to locate the component dependencies.
func (c *RestService) SetReferences(references crefer.IReferences) {
	c.references = references

	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)

	// Get endpoint
	depRes := c.DependencyResolver.GetOneOptional("endpoint")
	if depRes != nil {
		c.Endpoint = depRes.(*HttpEndpoint)
	}

	// Or create a local one
	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.localEndpoint = true
	} else {
		c.localEndpoint = false
	}
	// Add registration callback to the endpoint
	c.Endpoint.Register(c)

	depRes = c.DependencyResolver.GetOneOptional("swagger")
	if depRes != nil {
		c.SwaggerService = depRes.(ISwaggerService)
	}
}

// UnsetReferences method are unsets (clears) previously set references to dependent components.
func (c *RestService) UnsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.Unregister(c)
		c.Endpoint = nil
	}
	c.SwaggerService = nil
}

func (c *RestService) createEndpoint() *HttpEndpoint {
	endpoint := NewHttpEndpoint()

	if c.config != nil {
		endpoint.Configure(c.config)
	}
	if c.references != nil {
		endpoint.SetReferences(c.references)
	}

	return endpoint
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
// Parameters:
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - name              a method name.
// Returns Timing object to end the time measurement.
func (c *RestService) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".exec_count")
	return c.Counters.BeginTiming(name + ".exec_time")
}

// InstrumentError method are adds instrumentation to error handling.
// Parameters:
//   - correlationId  string  (optional) transaction id to trace execution through call chain.
//   - name    string          a method name.
//   - err     error          an occured error
//   - result  interface{}    (optional) an execution result
// Returns:  result interface{}, err error
//        (optional) an execution callback
func (c *RestService) InstrumentError(correlationId string, name string, errIn error,
	resIn interface{}) (result interface{}, err error) {
	if errIn != nil {
		c.Logger.Error(correlationId, errIn, "Failed to execute %s method", name)
		c.Counters.IncrementOne(name + ".exec_errors")
	}
	return resIn, errIn
}

// IsOpen method checks if the component is opened.
// Returrns true if the component has been opened and false otherwise.
func (c *RestService) IsOpen() bool {
	return c.opened
}

// Open method are opens the component.
// Parameters:
// 	 - correlationId  string:	(optional) transaction id to trace execution through call chain.
//  Returns: error
// error or nil no errors occured.
func (c *RestService) Open(correlationId string) error {
	if c.opened {
		return nil
	}

	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.Endpoint.Register(c)
		c.localEndpoint = true
	}

	if c.localEndpoint {
		oErr := c.Endpoint.Open(correlationId)
		if oErr != nil {
			c.opened = false
			return oErr
		}
	}
	c.opened = true
	return nil
}

// Close method are closes component and frees used resources.
// Parameters:
//  - correlationId 	(optional) transaction id to trace execution through call chain.
// Returns: error
// error or nil no errors occured.
func (c *RestService) Close(correlationId string) error {
	if !c.opened {
		return nil
	}

	if c.Endpoint == nil {
		return cerr.NewInvalidStateError(correlationId, "NO_ENDPOINT", "HTTP endpoint is missing")
	}

	if c.localEndpoint {
		cErr := c.Endpoint.Close(correlationId)
		if cErr != nil {
			c.opened = false
			return cErr
		}
	}
	c.opened = false
	return nil
}

// SendResult method method are sends result as JSON object.
// That function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
// Parameters:
//   - req       a HTTP request object.
//   - res       a HTTP response object.
//   - result    (optional) result object to send
//   - err error (optional) error objrct to send
func (c *RestService) SendResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendResult(res, req, result, err)
}

// SendCreatedResult method are sends newly created object as JSON.
// That callack function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 201 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
// Parameters:
//  - req       a HTTP request object.
//  - res       a HTTP response object.
//  - result    (optional) result object to send
//  - err error (optional) error objrct to send
func (c *RestService) SendCreatedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendCreatedResult(res, req, result, err)
}

// SendDeletedResult method are sends deleted object as JSON.
// That callack function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
// Parameters:
//   - req       a HTTP request object.
//   - res       a HTTP response object.
// 	 - result    (optional) result object to send
// 	 - err error (optional) error objrct to send
func (c *RestService) SendDeletedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendDeletedResult(res, req, result, err)
}

// SendError method are sends error serialized as ErrorDescription object
// and appropriate HTTP status code.
// If status code is not defined, it uses 500 status code.
// Parameters:
//   - req       a HTTP request object.
//   - res       a HTTP response object.
//   - error     an error object to be sent.
func (c *RestService) SendError(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendError(res, req, err)
}

func (c *RestService) appendBaseRoute(route string) string {

	if route == "" {
		route = "/"
	}

	if c.BaseRoute != "" && len(c.BaseRoute) > 0 {
		baseRoute := c.BaseRoute
		if len(route) == 0 {
			route = "/"
		}
		if route[0] != '/' {
			route = "/" + route
		}
		if baseRoute[0] != '/' {
			baseRoute = "/" + baseRoute
		}
		route = baseRoute + route
	}
	return route

}

// RegisterRoute method are registers a route in HTTP endpoint.
// Parameters:
//   - method        HTTP method: "get", "head", "post", "put", "delete"
//   - route         a command route. Base route will be added to this route
//   - schema        a validation schema to validate received parameters.
//   - action        an action function that is called when operation is invoked.
func (c *RestService) RegisterRoute(method string, route string, schema *cvalid.Schema,
	action func(res http.ResponseWriter, req *http.Request)) {
	if c.Endpoint == nil {
		return
	}
	route = c.appendBaseRoute(route)
	c.Endpoint.RegisterRoute(method, route, schema, action)
}

// RegisterRouteWithAuth method are registers a route with authorization in HTTP endpoint.
// Parameters:
//   - method        HTTP method: "get", "head", "post", "put", "delete"
//   - route         a command route. Base route will be added to this route
//   - schema        a validation schema to validate received parameters.
//   - authorize     an authorization interceptor
//   - action        an action function that is called when operation is invoked.
func (c *RestService) RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
	authorize func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc),
	action func(res http.ResponseWriter, req *http.Request)) {
	if c.Endpoint == nil {
		return
	}
	route = c.appendBaseRoute(route)
	c.Endpoint.RegisterRouteWithAuth(
		method, route, schema,
		func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
			if authorize != nil {
				authorize(res, req, next)
			} else {
				next.ServeHTTP(res, req)
			}
		}, action)
}

// RegisterInterceptor method are registers a middleware for a given route in HTTP endpoint.
// Parameters:
//   - route         a command route. Base route will be added to this route
//   - action        an action function that is called when middleware is invoked.
func (c *RestService) RegisterInterceptor(route string,
	action func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)) {
	if c.Endpoint == nil {
		return
	}

	route = c.appendBaseRoute(route)

	c.Endpoint.RegisterInterceptor(
		route, action)
}

// GetParam methods helps get all params from query
//   - req   - incoming request
//   - name  - parameter name
// Returns value or empty string if param not exists
func (c *RestService) GetParam(req *http.Request, name string) string {
	param := req.URL.Query().Get(name)
	if param == "" {
		param = mux.Vars(req)[name]
	}
	return param
}

// DecodeBody methods helps decode body
//   - req   	- incoming request
//   - target  	- pointer on target variable for decode
// Returns error
func (c *RestService) DecodeBody(req *http.Request, target interface{}) error {

	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	err = json.Unmarshal(bytes, target)
	if err != nil {
		return err
	}
	return nil
}

// GetPagingParams methods helps decode paging params
//   - req   	- incoming request
// Returns paging params
func (c *RestService) GetPagingParams(req *http.Request) *cdata.PagingParams {

	pagingParams := make(map[string]string, 0)
	pagingParams["skip"] = c.GetParam(req, "skip")
	pagingParams["take"] = c.GetParam(req, "take")
	pagingParams["total"] = c.GetParam(req, "total")

	return cdata.NewPagingParamsFromValue(pagingParams)
}

// GetFilterParams methods helps decode filter params
//   - req   	- incoming request
// Returns filter params
func (c *RestService) GetFilterParams(req *http.Request) *cdata.FilterParams {

	params := req.URL.Query()

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")
	delete(params, "correlation_id")

	return cdata.NewFilterParamsFromValue(params)
}

// GetCorrelationId method returns CorrelationId from request
// Parameters:
//   req *http.Request  request
// Returns: string
// retrun correlation_id or empty string
func (c *RestService) GetCorrelationId(req *http.Request) string {
	correlationId := req.URL.Query().Get("correlation_id")
	if correlationId == "" {
		correlationId = req.Header.Get("correlation_id")
	}
	return correlationId
}

func (c *RestService) RegisterOpenApiSpecFromFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.Logger.Error("RestService", err, "Can't read swagger file by path %s", path)
		return
	}
	c.RegisterOpenApiSpec((string)(content))
}

func (c *RestService) RegisterOpenApiSpec(content string) {
	if c.SwaggerEnable {
		c.RegisterRoute("get", c.SwaggerRoute, nil, func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add("Content-Length", cconv.StringConverter.ToString(len(content)))
			res.Header().Add("Content-Type", "application/x-yaml")
			res.WriteHeader(200)
			io.WriteString(res, content)
		})

		if c.SwaggerService != nil {
			c.SwaggerService.RegisterOpenApiSpec(c.BaseRoute, c.SwaggerRoute)
		}
	}
}
