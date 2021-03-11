package test_rpc_services

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type DummyCommandableHttpService struct {
	*services.CommandableHttpService
}

func NewDummyCommandableHttpService() *DummyCommandableHttpService {
	c := DummyCommandableHttpService{
		CommandableHttpService: services.NewCommandableHttpService("dummies"),
	}
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.CommandableHttpService.IRegisterable = &c
	return &c
}

func (c *DummyCommandableHttpService) Register() {
	if !c.SwaggerAuto && c.SwaggerEnable {
		c.RegisterOpenApiSpec("swagger yaml content")
	}
	c.CommandableHttpService.Register()
}
