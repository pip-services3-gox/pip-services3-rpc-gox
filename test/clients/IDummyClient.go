package test_rpc_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	testrpc "github.com/pip-services3-gox/pip-services3-rpc-gox/test"
)

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error)
	CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error)
	UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error)

	CheckCorrelationId(correlationId string) (result map[string]string, err error)
}
