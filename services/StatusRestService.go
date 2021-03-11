package services

import (
	"net/http"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
)

/*
StatusRestService is a service that returns microservice status information via HTTP/REST protocol.

The service responds on /status route (can be changed) with a JSON object:
  {
    - "id":            unique container id (usually hostname)
    - "name":          container name (from ContextInfo)
    - "description":   container description (from ContextInfo)
    - "start_time":    time when container was started
    - "current_time":  current time in UTC
    - "uptime":        duration since container start time in milliseconds
    - "properties":    additional container properties (from ContextInfo)
    - "components":    descriptors of components registered in the container
  }

Configuration parameters:

  - baseroute:              base route for remote URI
  - route:                   status route (default: "status")
  - dependencies:
    - endpoint:              override for HTTP Endpoint dependency
    - controller:            override for Controller dependency
  - connection(s):
    - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
    - protocol:              connection protocol: http or https
    - host:                  host name or IP address
    - port:                  port number
    - uri:                   resource URI or connection string with all parameters in it

References:

- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0            (optional) IDiscovery services to resolve connection
- *:endpoint:http:*:1.0          (optional) HttpEndpoint reference

See: RestService
See: RestClient

Example:

    service = NewStatusService();
    service.Configure(cref.NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));

	opnErr:= service.Open("123")
	if opnErr == nil {
       fmt.Println("The Status service is accessible at http://localhost:8080/status");
    }
*/
type StatusRestService struct {
	*RestService
	startTime   time.Time
	references2 crefer.IReferences
	contextInfo *cinfo.ContextInfo
	route       string
}

// NewStatusRestService method are creates a new instance of this service.
func NewStatusRestService() *StatusRestService {
	c := StatusRestService{}
	c.RestService = NewRestService()
	c.RestService.IRegisterable = &c
	c.startTime = time.Now()
	c.route = "status"
	c.DependencyResolver.Put("context-info", crefer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	return &c
}

// Configure method are configures component by passing configuration parameters.
// Parameters:
//   - config  *cconf.ConfigParams  configuration parameters to be set.
func (c *StatusRestService) Configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
	c.route = config.GetAsStringWithDefault("route", c.route)
}

// SetReferences method are sets references to dependent components.
// Parameters:
//  - references crefer.IReferences	references to locate the component dependencies.
func (c *StatusRestService) SetReferences(references crefer.IReferences) {
	c.references2 = references
	c.RestService.SetReferences(references)

	depRes := c.DependencyResolver.GetOneOptional("context-info")
	if depRes != nil {
		c.contextInfo = depRes.(*cinfo.ContextInfo)
	}
}

// Register method are registers all service routes in HTTP endpoint.
func (c *StatusRestService) Register() {
	c.RegisterRoute("get", c.route, nil, c.status)
}

// Handles status requests
//   - req  *http.Request an HTTP request
//   - res  http.ResponseWriter  an HTTP response
func (c *StatusRestService) status(res http.ResponseWriter, req *http.Request) {

	id := ""
	if c.contextInfo != nil {
		id = c.contextInfo.ContextId
	}

	name := "Unknown"
	if c.contextInfo != nil {
		name = c.contextInfo.Name
	}

	description := ""
	if c.contextInfo != nil {
		description = c.contextInfo.Description
	}

	uptime := time.Now().Sub(c.startTime)

	properties := make(map[string]string, 0)
	if c.contextInfo != nil {
		properties = c.contextInfo.Properties
	}

	var components []string
	if c.references2 != nil {
		for _, locator := range c.references2.GetAllLocators() {
			components = append(components, cconv.StringConverter.ToString(locator))
		}
	}

	status := make(map[string]interface{})

	status["id"] = id
	status["name"] = name
	status["description"] = description
	status["start_time"] = cconv.StringConverter.ToString(c.startTime)
	status["current_time"] = cconv.StringConverter.ToString(time.Now())
	status["uptime"] = uptime
	status["properties"] = properties
	status["components"] = components

	c.SendResult(res, req, status, nil)
}
