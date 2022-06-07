package test

import (
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
)

type TestRestClient struct {
	clients.RestClient
}

func NewTestRestClient(baseRoute string) *TestRestClient {
	c := &TestRestClient{}
	c.RestClient = *clients.NewRestClient()
	c.BaseRoute = baseRoute
	return c
}
