package test_services

import (
	"context"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type DummyCommandableHttpService struct {
	services.CommandableHttpService
}

func NewDummyCommandableHttpService() *DummyCommandableHttpService {
	c := &DummyCommandableHttpService{}
	c.CommandableHttpService = *services.InheritCommandableHttpService(c, "dummies")
	c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyCommandableHttpService) Register() {
	if !c.SwaggerAuto && c.SwaggerEnabled {
		c.RegisterOpenApiSpec("swagger yaml content")
	}
	c.CommandableHttpService.Register(context.Background())
}
