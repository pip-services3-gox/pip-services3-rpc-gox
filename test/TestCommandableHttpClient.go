package test

import (
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
)

type TestCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewTestCommandableHttpClient(baseRoute string) *TestCommandableHttpClient {
	c := &TestCommandableHttpClient{}
	c.CommandableHttpClient = *clients.NewCommandableHttpClient(baseRoute)
	return c
}
