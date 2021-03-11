package services

import (
	"net/http"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

/*
HeartbeatRestService service returns heartbeat via HTTP/REST protocol.

The service responds on /heartbeat route (can be changed)
with a string with the current time in UTC.

This service route can be used to health checks by loadbalancers and
container orchestrators.

Configuration parameters:

  - baseroute:              base route for remote URI (default: "")
  - route:                   route to heartbeat operation (default: "heartbeat")
  - dependencies:
    - endpoint:              override for HTTP Endpoint dependency
  - connection(s):
    - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
    - protocol:              connection protocol: http or https
    - host:                  host name or IP address
    - port:                  port number
    - uri:                   resource URI or connection string with all parameters in it

References:

- *:logger:*:*:1.0               (optional)  ILogger components to pass log messages
- *:counters:*:*:1.0             (optional)  ICounters components to pass collected measurements
- *:discovery:*:*:1.0            (optional)  IDiscovery services to resolve connection
- *:endpoint:http:*:1.0          (optional) HttpEndpoint reference

See RestService
See RestClient

Example:

    service := NewHeartbeatService();
    service.Configure(cconf.NewConfigParamsFromTuples(
        "route", "ping",
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));

    opnErr := service.Open("123")
    if opnErr == nil {
       fmt.Println("The Heartbeat service is accessible at http://+:8080/ping");
    }
*/
type HeartbeatRestService struct {
	*RestService
	route string
}

/**
  Creates a new instance of c service.
*/
func NewHeartbeatRestService() *HeartbeatRestService {
	c := HeartbeatRestService{}
	c.RestService = NewRestService()
	c.RestService.IRegisterable = &c
	c.route = "heartbeat"
	return &c
}

/**
  Configures component by passing configuration parameters.

  @param config    configuration parameters to be set.
*/
func (c *HeartbeatRestService) Configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
	c.route = config.GetAsStringWithDefault("route", c.route)
}

/**
  Registers all service routes in HTTP endpoint.
*/
func (c *HeartbeatRestService) Register() {
	c.RegisterRoute("get", c.route, nil, func(res http.ResponseWriter, req *http.Request) { c.heartbeat(req, res) })
}

/**
  Handles heartbeat requests

  @param req   an HTTP request
  @param res   an HTTP response
*/
func (c *HeartbeatRestService) heartbeat(req *http.Request, res http.ResponseWriter) {
	c.SendResult(res, req, time.Now(), nil)
}
