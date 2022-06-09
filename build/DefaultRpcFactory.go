package build

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// DefaultRpcFactory are creates RPC components by their descriptors
//	see Factory
//	see HttpEndpoint
//	see HeartbeatRestService
//	see StatusRestService
type DefaultRpcFactory struct {
	cbuild.Factory
}

// NewDefaultRpcFactory creates a new instance of the factory.
func NewDefaultRpcFactory() *DefaultRpcFactory {
	c := DefaultRpcFactory{}
	c.Factory = *cbuild.NewFactory()

	httpEndpointDescriptor := cref.NewDescriptor("pip-services", "endpoint", "http", "*", "1.0")
	statusServiceDescriptor := cref.NewDescriptor("pip-services", "status-service", "http", "*", "1.0")
	heartbeatServiceDescriptor := cref.NewDescriptor("pip-services", "heartbeat-service", "http", "*", "1.0")

	c.RegisterType(httpEndpointDescriptor, services.NewHttpEndpoint)
	c.RegisterType(heartbeatServiceDescriptor, services.NewHeartbeatRestService)
	c.RegisterType(statusServiceDescriptor, services.NewStatusRestService)
	return &c
}
