package build

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// DefaultRpcFactory are creates RPC components by their descriptors.

// See Factory
// See HttpEndpoint
// See HeartbeatRestService
// See StatusRestService
type DefaultRpcFactory struct {
	cbuild.Factory
	Descriptor                 *cref.Descriptor
	HttpEndpointDescriptor     *cref.Descriptor
	StatusServiceDescriptor    *cref.Descriptor
	HeartbeatServiceDescriptor *cref.Descriptor
}

// NewDefaultRpcFactorymethod create a new instance of the factory.
func NewDefaultRpcFactory() *DefaultRpcFactory {
	c := DefaultRpcFactory{}
	c.Factory = *cbuild.NewFactory()
	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "rpc", "default", "1.0")
	c.HttpEndpointDescriptor = cref.NewDescriptor("pip-services", "endpoint", "http", "*", "1.0")
	c.StatusServiceDescriptor = cref.NewDescriptor("pip-services", "status-service", "http", "*", "1.0")
	c.HeartbeatServiceDescriptor = cref.NewDescriptor("pip-services", "heartbeat-service", "http", "*", "1.0")

	c.RegisterType(c.HttpEndpointDescriptor, services.NewHttpEndpoint)
	c.RegisterType(c.HeartbeatServiceDescriptor, services.NewHeartbeatRestService)
	c.RegisterType(c.StatusServiceDescriptor, services.NewStatusRestService)
	return &c
}
