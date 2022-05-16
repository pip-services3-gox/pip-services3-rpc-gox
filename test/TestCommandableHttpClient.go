package test

import (
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
)

type TestCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewTestCommandableHttpClient(baseRoute string) *TestCommandableHttpClient {
	c := &TestCommandableHttpClient{}
	c.CommandableHttpClient = *clients.NewCommandableHttpClient(baseRoute)
	return c
}
