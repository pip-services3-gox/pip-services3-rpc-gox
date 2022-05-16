package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	ccomands "github.com/pip-services3-go/pip-services3-commons-go/commands"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
)

/*
CommandableHttpService abstract service that receives remove calls via HTTP/REST protocol
to operations automatically generated for commands defined in ICommandable components.
Each command is exposed as POST operation that receives all parameters in body object.

Commandable services require only 3 lines of code to implement a robust external
HTTP-based remote interface.

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

References:

- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0            (optional) IDiscovery services to resolve connection
- *:endpoint:http:*:1.0          (optional) HttpEndpoint reference

See CommandableHttpClient
See RestService

Example:

	type MyCommandableHttpService struct {
	 *CommandableHttpService
	}

	func NewMyCommandableHttpService() *MyCommandableHttpService {
		c := MyCommandableHttpService{
			CommandableHttpService: services.NewCommandableHttpService("dummies"),
		}
		c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &c

    service := NewMyCommandableHttpService();
    service.Configure(cconf.NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));
    service.SetReferences(cref.NewReferencesFromTuples(
       cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller
    ));

	opnErr:=service.Open("123")
	if opnErr == nil {
       fmt.Println("The REST service is running on port 8080");
	}

*/
type CommandableHttpService struct {
	*RestService
	commandSet  *ccomands.CommandSet
	SwaggerAuto bool
}

// NewCommandableHttpService creates a new instance of the service.
// Parameters:
//   - baseRoute string a service base route.
// Returns: *CommandableHttpService
// pointer on new instance CommandableHttpService
// func NewCommandableHttpService(baseRoute string) *CommandableHttpService {
// 	c := &CommandableHttpService{}
// 	c.RestService = InheritRestService(c)
// 	c.BaseRoute = baseRoute
// 	c.SwaggerAuto = true
// 	c.DependencyResolver.Put("controller", "none")
// 	return c
// }

// InheritCommandableHttpService creates a new instance of the service.
// Parameters:
//   - overrides references to child class that overrides virtual methods
//   - baseRoute string a service base route.
// Returns: *CommandableHttpService
// pointer on new instance CommandableHttpService
func InheritCommandableHttpService(overrides IRestServiceOverrides, baseRoute string) *CommandableHttpService {
	c := &CommandableHttpService{}
	c.RestService = InheritRestService(overrides)
	c.BaseRoute = baseRoute
	c.SwaggerAuto = true
	c.DependencyResolver.Put("controller", "none")
	return c
}

//  Configure method configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *CommandableHttpService) Configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
	c.SwaggerAuto = config.GetAsBooleanWithDefault("swagger.auto", c.SwaggerAuto)
}

// Register method are registers all service routes in HTTP endpoint.
func (c *CommandableHttpService) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error("CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}
	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		route := command.Name()
		if route[0] != "/"[0] {
			route = "/" + route
		}

		c.RegisterRoute("post", route, nil, func(res http.ResponseWriter, req *http.Request) {

			// Make copy of request
			bodyBuf, bodyErr := ioutil.ReadAll(req.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(res, req, bodyErr)
				return
			}
			req.Body.Close()
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var params map[string]interface{} = make(map[string]interface{}, 0)
			json.Unmarshal(bodyBuf, &params)

			urlParams := req.URL.Query()
			for k, v := range urlParams {
				params[k] = v[0]
			}
			for k, v := range mux.Vars(req) {
				params[k] = v
			}

			correlationId := c.GetCorrelationId(req)
			args := crun.NewParametersFromValue(params)
			timing := c.Instrument(correlationId, c.BaseRoute+"."+command.Name())

			execRes, execErr := command.Execute(correlationId, args)
			timing.EndTiming(execErr)
			c.SendResult(res, req, execRes, execErr)

		})
	}

	if c.SwaggerAuto {
		var swaggerConfig = c.config.GetSection("swagger")
		var doc = NewCommandableSwaggerDocument(c.BaseRoute, swaggerConfig, commands)
		c.RegisterOpenApiSpec(doc.ToString())
	}
}
