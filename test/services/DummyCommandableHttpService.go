package test_services

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type DummyCommandableHttpService struct {
	services.CommandableHttpService
}

func NewDummyCommandableHttpService() *DummyCommandableHttpService {
	c := &DummyCommandableHttpService{}
	c.CommandableHttpService = *services.InheritCommandableHttpService(c, "dummies")
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyCommandableHttpService) Register() {
	if !c.SwaggerAuto && c.SwaggerEnabled {
		c.RegisterOpenApiSpec("swagger yaml content")
	}
	c.CommandableHttpService.Register()
}
