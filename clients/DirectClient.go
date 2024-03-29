package clients

import (
	"context"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
	service "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// DirectClient is abstract client that calls controller directly in the same memory space.
//
// It is used when multiple microservices are deployed in a single container (monolyth)
// and communication between them can be done by direct calls rather then through
// the network.
//
//	Configuration parameters:
//		- dependencies:
//			- controller: override controller descriptor
//
//	References:
//		- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0       (optional) ICounters components to pass collected measurements
//		- *:controller:*:*:1.0     controller to call business methods
//
//	Example:
//		type MyDirectClient struct {
//			*DirectClient
//		}
//		func NewMyDirectClient()* MyDirectClient {
//			c:= MyDirectClient{}
//			c.DirectClient = NewDirectClient()
//			c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor(
//              "mygroup", "controller", "*", "*", "*"));
//			return &c
//		}
//
//		func (c *MyDirectClient) SetReferences(ctx context.Context, references cref.IReferences) {
//			c.DirectClient.SetReferences(ctx, references)
//			specificController, ok := c.Controller.(tdata.IMyDataController)
//			if !ok {
//				panic("MyDirectClient: Cant't resolv dependency 'controller' to IMyDataController")
//			}
//			c.specificController = specificController
//		}
//		...
//		func (c * MyDirectClient) GetData(ctx context.Context, correlationId string, id string)(result MyData, err error) {
//			timing := c.Instrument(ctx, correlationId, "myclient.get_data")
//			defer timing.EndTiming(ctx);
//
//			return c.specificController.GetData(ctx, correlationId, id)
//		}
//		...
//
//		client = NewMyDirectClient();
//		client.SetReferences(context.Background(), cref.NewReferencesFromTuples(
//			cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller,
//		));
//		res, err := client.GetData(context.Background(), "123", "1")
type DirectClient struct {
	//The controller reference.
	Controller any
	//The open flag.
	Opened bool
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters
	Counters *ccount.CompositeCounters
	//The dependency resolver to get controller reference.
	DependencyResolver *crefer.DependencyResolver
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

// NewDirectClient is creates a new instance of the client.
func NewDirectClient() *DirectClient {
	dc := DirectClient{
		Opened:             true,
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		DependencyResolver: crefer.NewDependencyResolver(),
		Tracer:             ctrace.NewCompositeTracer(),
	}
	dc.DependencyResolver.Put(context.Background(), "controller", "none")
	return &dc
}

// Configure method are configures component by passing configuration parameters.
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams  configuration parameters to be set.
func (c *DirectClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(ctx, config)
}

// SetReferences method are sets references to dependent components.
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences references to locate the component dependencies.
func (c *DirectClient) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
	res, cErr := c.DependencyResolver.GetOneRequired("controller")
	if cErr != nil {
		panic("DirectClient: Cant't resolv dependency 'controller'")
	}
	c.Controller = res
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//	Parameters:
//		- ctx context.Context
//		- correlationId  string (optional) transaction id to trace execution through call chain.
//		- name string a method name.
//	Returns: Timing object to end the time measurement.
func (c *DirectClient) Instrument(ctx context.Context, correlationId string, name string) *service.InstrumentTiming {
	c.Logger.Trace(ctx, correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, correlationId, name, "")
	return service.NewInstrumentTiming(correlationId, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// InstrumentError method are adds instrumentation to error handling.
// Parameters:
//    - correlationId     (optional) transaction id to trace execution through call chain.
//    - name              a method name.
//    - err               an occured error
//    - result            (optional) an execution result
// Retruns:          result any, err error
// an execution result and error
// func (c *DirectClient) InstrumentError(correlationId string, name string, inErr error, inRes any) (result any, err error) {
// 	if inErr != nil {
// 		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
// 		c.Counters.IncrementOne(name + ".call_errors")
// 	}
// 	return inRes, inErr
// }

// IsOpen method are checks if the component is opened.
//	Returns: true if the component has been opened and false otherwise.
func (c *DirectClient) IsOpen() bool {
	return c.Opened
}

// Open method are opens the component.
//	Parameters:
//		- ctx context.Context
//		- correlationId string	(optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *DirectClient) Open(ctx context.Context, correlationId string) error {
	if c.Opened {
		return nil
	}

	if c.Controller == nil {
		err := cerr.NewConnectionError(correlationId, "NO_CONTROLLER", "Controller reference is missing")
		return err
	}

	c.Opened = true

	c.Logger.Info(ctx, correlationId, "Opened direct client")
	return nil
}

// Close method are closes component and frees used resources.
//	Parameters:
//		- ctx context.Context
//		- correlationId string	(optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *DirectClient) Close(ctx context.Context, correlationId string) error {
	if c.Opened {
		c.Logger.Info(ctx, correlationId, "Closed direct client")
	}
	c.Opened = false
	return nil
}
