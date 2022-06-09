package test_clients

import (
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
)

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error)
	CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error)
	UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error)

	CheckCorrelationId(correlationId string) (result map[string]string, err error)

	CheckErrorPropagation(correlationId string) error
}
