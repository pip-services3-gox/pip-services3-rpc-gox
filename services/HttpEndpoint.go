package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	"github.com/pip-services3-go/pip-services3-rpc-go/connect"
)

/*
 HttpEndpoint used for creating HTTP endpoints. An endpoint is a URL, at which a given service can be accessed by a client.

Configuration parameters:

Parameters to pass to the configure method for component configuration:

  - connection(s) - the connection resolver"s connections:
    - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
    - "connection.protocol" - the connection"s protocol;
    - "connection.host" - the target host;
    - "connection.port" - the target port;
    - "connection.uri" - the target URI.
  - credential - the HTTPS credentials:
    - "credential.ssl_key_file" - the SSL func (c *HttpEndpoint )key in PEM
    - "credential.ssl_crt_file" - the SSL certificate in PEM
    - "credential.ssl_ca_file" - the certificate authorities (root cerfiticates) in PEM

  - cors-headers - pair CORS headers: origin. Example: MyHeader1: \*.\*

References:

A logger, counters, and a connection resolver can be referenced by passing the
following references to the object"s setReferences method:

  - logger: "*:logger:*:*:1.0";
  - counters: "*:counters:*:*:1.0";
  - discovery: "*:discovery:*:*:1.0" (for the connection resolver).

Examples:

    endpoint := NewHttpEndpoint();
    endpoint.Configure(config);
    endpoint.SetReferences(references);
    ...
	endpoint.Open(correlationId)
*/
type HttpEndpoint struct {
	defaultConfig          *cconf.ConfigParams
	server                 *http.Server
	router                 *mux.Router
	connectionResolver     *connect.HttpConnectionResolver
	logger                 *clog.CompositeLogger
	counters               *ccount.CompositeCounters
	maintenanceEnabled     bool
	fileMaxSize            int64
	protocolUpgradeEnabled bool
	uri                    string
	registrations          []IRegisterable
	allowedHeaders         []string
	allowedOrigins         []string
}

// NewHttpEndpoint creates new HttpEndpoint
func NewHttpEndpoint() *HttpEndpoint {
	c := HttpEndpoint{}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", "3000",

		"credential.ssl_key_file", nil,
		"credential.ssl_crt_file", nil,
		"credential.ssl_ca_file", nil,

		"options.maintenance_enabled", false,
		"options.request_max_size", 1024*1024,
		"options.file_max_size", 200*1024*1024,
		"options.connect_timeout", "60000",
		"options.debug", "true",
	)
	c.connectionResolver = connect.NewHttpConnectionResolver()
	c.logger = clog.NewCompositeLogger()
	c.counters = ccount.NewCompositeCounters()
	c.maintenanceEnabled = false
	c.fileMaxSize = 200 * 1024 * 1024
	c.protocolUpgradeEnabled = false
	c.registrations = make([]IRegisterable, 0, 0)
	c.allowedHeaders = []string{
		//"Accept",
		//"Content-Type",
		//"Content-Length",
		//"Accept-Encoding",
		//"X-CSRF-Token",
		//"Authorization",
		"correlation_id",
		//"access_token",
	}
	c.allowedOrigins = make([]string, 0)
	return &c
}

// Configure method are configures this HttpEndpoint using the given configuration parameters.
// Configuration parameters:
//    - connection(s) - the connection resolver"s connections;
//        - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
//        - "connection.protocol" - the connection"s protocol;
//        - "connection.host" - the target host;
//        - "connection.port" - the target port;
//        - "connection.uri" - the target URI.
//        - "credential.ssl_key_file" - SSL func (c *HttpEndpoint )key in PEM
//        - "credential.ssl_crt_file" - SSL certificate in PEM
//        - "credential.ssl_ca_file" - Certificate authority (root certificate) in PEM
//  - config    configuration parameters, containing a "connection(s)" section.
func (c *HttpEndpoint) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.connectionResolver.Configure(config)

	c.maintenanceEnabled = config.GetAsBooleanWithDefault("options.maintenance_enabled", c.maintenanceEnabled)
	c.fileMaxSize = config.GetAsLongWithDefault("options.file_max_size", c.fileMaxSize)
	c.protocolUpgradeEnabled = config.GetAsBooleanWithDefault("options.protocol_upgrade_enabled", c.protocolUpgradeEnabled)

	headers := strings.Split(config.GetAsStringWithDefault("cors_headers", ""), ",")
	if headers != nil && len(headers) > 0 {
		for _, header := range headers {
			c.AddCorsHeader(strings.TrimSpace(header), "")
		}
	}

	origins := strings.Split(config.GetAsStringWithDefault("cors_origins", ""), ",")
	if origins != nil && len(origins) > 0 {
		for _, origin := range origins {
			c.AddCorsHeader("", strings.TrimSpace(origin))
		}
	}
}

// SetReferences method are sets references to this endpoint"s logger, counters, and connection resolver.
//    References:
//    - logger: "*:logger:*:*:1.0"
//    - counters: "*:counters:*:*:1.0"
//    - discovery: "*:discovery:*:*:1.0" (for the connection resolver)
// Parameters:
//    - references    an IReferences object, containing references to a logger, counters,
//     and a connection resolver.
func (c *HttpEndpoint) SetReferences(references crefer.IReferences) {
	c.logger.SetReferences(references)
	c.counters.SetReferences(references)
	c.connectionResolver.SetReferences(references)
}

// IsOpen method is  whether or not this endpoint is open with an actively listening REST server.
func (c *HttpEndpoint) IsOpen() bool {
	return c.server != nil
}

// Opens a connection using the parameters resolved by the referenced connection
// resolver and creates a REST server (service) using the set options and parameters.
// Parameters:
//   - correlationId   string  (optional) transaction id to trace execution through call chain.
// Returns : error
// an error if one is raised.
func (c *HttpEndpoint) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}
	connection, credential, err := c.connectionResolver.Resolve(correlationId)
	if err != nil {
		return err
	}

	c.uri = connection.Uri()
	url := connection.Host() + ":" + strconv.Itoa(connection.Port())
	c.server = &http.Server{Addr: url}
	c.router = mux.NewRouter()

	// Add default origins
	// if len(c.allowedOrigins) == 0 {
	// 	c.allowedOrigins = []string{"*"}
	// }

	allowedOrigins := handlers.AllowedOrigins(c.allowedOrigins)
	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"DELETE",
		"PATCH",
	})
	allowedHeaders := handlers.AllowedHeaders(c.allowedHeaders)
	c.server.Handler = handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(c.router)

	c.router.Use(c.noCache)
	c.router.Use(c.doMaintenance)

	c.performRegistrations()

	if connection.Protocol() == "https" {
		sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")

		go func() {
			servErr := c.server.ListenAndServeTLS(sslKeyFile, sslCrtFile)
			if servErr != nil {
				//fmt.Println("Server stoped:", servErr.Error())
			}
		}()

	} else {
		go func() {
			servErr := c.server.ListenAndServe()
			if servErr != nil {
				//fmt.Println("Server stoped:", servErr.Error())
			}
		}()
	}

	regErr := c.connectionResolver.Register(correlationId)
	if regErr != nil {
		c.logger.Error(correlationId, regErr, "ERROR_REG_SRV", "Can't register REST service at %s", c.uri)
	}
	c.logger.Debug(correlationId, "Opened REST service at %s", c.uri)
	return regErr
}

// Prevents IE from caching REST requests
func (c *HttpEndpoint) noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// Returns maintenance error code
func (c *HttpEndpoint) doMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make this more sophisticated
		if c.maintenanceEnabled {
			w.Header().Add("Retry-After", "3600")
			jsonStr, _ := json.Marshal(503)
			w.Write(jsonStr)
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Close method are closes this endpoint and the REST server (service) that was opened earlier.
// Parameters:
//   - correlationId  string   (optional) transaction id to trace execution through call chain.
// Returns: error
// an error if one is raised.
func (c *HttpEndpoint) Close(correlationId string) error {
	if c.server != nil {
		// Attempt a graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		clErr := c.server.Shutdown(ctx)
		if clErr != nil {
			c.logger.Warn(correlationId, "Failed while closing REST service: %s", clErr.Error())
			return clErr
		}
		c.logger.Debug(correlationId, "Closed REST service at %s", c.uri)
		c.server = nil
		c.uri = ""
	}
	return nil
}

// Registers a registerable object for dynamic endpoint discovery.
// Parameters:
//   - registration  IRegisterable   implements of IRegisterable interface.
// See IRegisterable
func (c *HttpEndpoint) Register(registration IRegisterable) {
	c.registrations = append(c.registrations, registration)
}

// Unregisters a registerable object, so that it is no longer used in dynamic
// endpoint discovery.
// Parameters:
//   - registration  IRegisterable  the registration to remove.
// See IRegisterable
func (c *HttpEndpoint) Unregister(registration IRegisterable) {
	for i := 0; i < len(c.registrations); {
		if c.registrations[i] == registration {
			if i == len(c.registrations)-1 {
				c.registrations = c.registrations[:i]
			} else {
				c.registrations = append(c.registrations[:i], c.registrations[i+1:]...)
			}
		} else {
			i++
		}
	}
}

func (c *HttpEndpoint) performRegistrations() {
	for _, registration := range c.registrations {
		registration.Register()
	}
}

func (c *HttpEndpoint) fixRoute(route string) string {
	if len(route) > 0 && !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	return route
}

// GetCorrelationId method returns CorrelationId from request
// Parameters:
//   req *http.Request  request
// Returns: string
// retrun correlation_id or empty string
func (c *HttpEndpoint) GetCorrelationId(req *http.Request) string {
	correlationId := req.URL.Query().Get("correlation_id")
	if correlationId == "" {
		correlationId = req.Header.Get("correlation_id")
	}
	return correlationId
}

// RegisterRoute method are registers an action in this objects REST server (service) by the given method and route.
//   - method   string     the HTTP method of the route.
//   - route    string     the route to register in this object"s REST server (service).
//   - schema   *cvalid.Schema     the schema to use for parameter validation.
//   - action   http.HandlerFunc     the action to perform at the given route.
func (c *HttpEndpoint) RegisterRoute(method string, route string, schema *cvalid.Schema,
	action http.HandlerFunc) {

	method = strings.ToLower(method)
	if method == "del" {
		method = "delete"
	}
	route = c.fixRoute(route)
	actionCurl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  Perform validation
		if schema != nil {
			var params map[string]interface{} = make(map[string]interface{}, 0)
			for k, v := range r.URL.Query() {
				params[k] = v[0]
			}

			for k, v := range mux.Vars(r) {
				params[k] = v
			}

			// Make copy of request
			bodyBuf, bodyErr := ioutil.ReadAll(r.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(w, r, bodyErr)
				return
			}
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var body interface{}
			json.Unmarshal(bodyBuf, &body)
			params["body"] = body

			correlationId := c.GetCorrelationId(r)
			err := schema.ValidateAndReturnError(correlationId, params, false)
			if err != nil {
				HttpResponseSender.SendError(w, r, err)
				return
			}
		}
		action(w, r)
	})
	c.router.Handle(route, actionCurl).Methods(strings.ToUpper(method))
}

// RegisterRouteWithAuth method are registers an action with authorization in this objects REST server (service)
// by the given method and route.
// Parameters:
//   - method    string    the HTTP method of the route.
//   - route     string    the route to register in this object"s REST server (service).
//   - schema    *cvalid.Schema    the schema to use for parameter validation.
//   - authorize     the authorization interceptor
//   - action        the action to perform at the given route.
func (c *HttpEndpoint) RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
	authorize func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc),
	action http.HandlerFunc) {

	if authorize != nil {
		nextAction := action
		action = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorize(w, r, nextAction)
		})
	}

	c.RegisterRoute(method, route, schema, action)
}

// RegisterInterceptor method are registers a middleware action for the given route.
//   - route         the route to register in this object"s REST server (service).
//   - action        the middleware action to perform at the given route.
func (c *HttpEndpoint) RegisterInterceptor(route string, action func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {

	route = c.fixRoute(route)
	interceptorFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched, _ := regexp.MatchString(route, r.URL.Path)
			if route != "" && !matched {
				next.ServeHTTP(w, r)
			} else {
				action(w, r, next.ServeHTTP)
			}
		})
	}
	c.router.Use(interceptorFunc)
}

// AddCORSHeader method adds allowed header, ignore if it already exist
// must be call before to opening endpoint
func (c *HttpEndpoint) AddCorsHeader(header string, origin string) {

	if len(header) > 0 {
		contain := false
		for _, allowedHeader := range c.allowedHeaders {
			if allowedHeader == header {
				contain = true
				break
			}
		}
		if !contain {
			c.allowedHeaders = append(c.allowedHeaders, header)
		}
	}
	if len(origin) > 0 {
		contain := false
		for _, allowedOrigin := range c.allowedOrigins {
			if allowedOrigin == origin {
				contain = true
				break
			}
		}
		if !contain {
			c.allowedOrigins = append(c.allowedOrigins, origin)
		}
	}
}
