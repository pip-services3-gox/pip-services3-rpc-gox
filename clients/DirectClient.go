package clients

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	ctrace "github.com/pip-services3-go/pip-services3-components-go/trace"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	service "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

/*
DirectClient is bstract client that calls controller directly in the same memory space.

It is used when multiple microservices are deployed in a single container (monolyth)
and communication between them can be done by direct calls rather then through
the network.

Configuration parameters:

- dependencies:
  - controller:            override controller descriptor

References:

- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
- *:counters:*:*:1.0       (optional) ICounters components to pass collected measurements
- *:controller:*:*:1.0     controller to call business methods

Example:

	type MyDirectClient struct {
	*DirectClient
	}
        func MyDirectClient()* MyDirectClient {
		  c:= MyDirectClient{}
		  c.DirectClient = NewDirectClient()
          c.DependencyResolver.Put("controller", cref.NewDescriptor(
              "mygroup", "controller", "*", "*", "*"));
		}

		func (c *MyDirectClient) SetReferences(references cref.IReferences) {
			c.DirectClient.SetReferences(references)
			specificController, ok := c.Controller.(tdata.IMyDataController)
			if !ok {
				panic("MyDirectClient: Cant't resolv dependency 'controller' to IMyDataController")
			}
			c.specificController = specificController
		}
        ...

        func (c * MyDirectClient) GetData(correlationId string, id string)(result MyData, err error) {
           timing := c.Instrument(correlationId, "myclient.get_data")
           cmRes, cmdErr := c.specificController.GetData(correlationId, id)
           timing.EndTiming();
           return  c.InstrumentError(correlationId, "myclient.get_data", cmdRes, cmdErr)
        }
        ...

    client = NewMyDirectClient();
    client.SetReferences(cref.NewReferencesFromTuples(
        cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller,
    ));

    res, err := client.GetData("123", "1")
*/
type DirectClient struct {
	//The controller reference.
	Controller interface{}
	//The open flag.
	Opened bool
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters
	Counters *ccount.CompositeCounters
	//The dependency resolver to get controller reference.
	DependencyResolver crefer.DependencyResolver
	// The tracer.
    Tracer *ctrace.CompositeTracer;
}

// NewDirectClient is creates a new instance of the client.
func NewDirectClient() *DirectClient {
	dc := DirectClient{
		Opened:             true,
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		DependencyResolver: *crefer.NewDependencyResolver(),
		Tracer: ctrace.NewCompositeTracer(nil),
	}
	dc.DependencyResolver.Put("controller", "none")
	return &dc
}

// Configure method are configures component by passing configuration parameters.
// Parameters:
//  - config  *cconf.ConfigParams  configuration parameters to be set.
func (c *DirectClient) Configure(config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(config)
}

// SetReferences method are sets references to dependent components.
// Parameters:
// - references  crefer.IReferences	references to locate the component dependencies.
func (c *DirectClient) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.Tracer.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
	res, cErr := c.DependencyResolver.GetOneRequired("controller")
	if cErr != nil {
		panic("DirectClient: Cant't resolv dependency 'controller'")
	}
	c.Controller = res
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
// Parameters:
//    - correlationId  string    (optional) transaction id to trace execution through call chain.
//    - name   string           a method name.
// Returns Timing object to end the time measurement.
func (c *DirectClient) Instrument(correlationId string, name string) *service.InstrumentTiming {
	c.Logger.Trace(correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	
	counterTiming := c.Counters.BeginTiming(name + ".call_time")
    traceTiming := c.Tracer.BeginTrace(correlationId, name, "")
    return service.NewInstrumentTiming(correlationId, name, "call",
            c.Logger, c.Counters, counterTiming, traceTiming)
}

// InstrumentError method are adds instrumentation to error handling.
// Parameters:
//    - correlationId     (optional) transaction id to trace execution through call chain.
//    - name              a method name.
//    - err               an occured error
//    - result            (optional) an execution result
// Retruns:          result interface{}, err error
// an execution result and error
// func (c *DirectClient) InstrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
// 	if inErr != nil {
// 		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
// 		c.Counters.IncrementOne(name + ".call_errors")
// 	}
// 	return inRes, inErr
// }

// IsOpen method are checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *DirectClient) IsOpen() bool {
	return c.Opened
}

// Open method are opens the component.
// 	- correlationId string	(optional) transaction id to trace execution through call chain.
// Returns: error
// error or nil no errors occured.
func (c *DirectClient) Open(correlationId string) error {
	if c.Opened {
		return nil
	}

	if c.Controller == nil {
		err := cerr.NewConnectionError(correlationId, "NO_CONTROLLER", "Controller reference is missing")
		return err
	}

	c.Opened = true

	c.Logger.Info(correlationId, "Opened direct client")
	return nil
}

// Close method are closes component and frees used resources.
// 	- correlationId string	(optional) transaction id to trace execution through call chain.
// Returns: error
// error or nil no errors occured.
func (c *DirectClient) Close(correlationId string) error {
	if c.Opened {
		c.Logger.Info(correlationId, "Closed direct client")
	}
	c.Opened = false
	return nil
}
