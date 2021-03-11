package test_rpc_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
	testrpc "github.com/pip-services3-gox/pip-services3-rpc-gox/test"
)

var (
	dummyDataPageType = reflect.TypeOf(&testrpc.DummyDataPage{})
	dummyType         = reflect.TypeOf(&testrpc.Dummy{})
)

type DummyRestClient struct {
	clients.RestClient
}

func NewDummyRestClient() *DummyRestClient {
	drc := DummyRestClient{}
	drc.RestClient = *clients.NewRestClient()
	return &drc
}

func (c *DummyRestClient) GetDummies(correlationId string, filter *cdata.FilterParams,
	paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.Call(dummyDataPageType, "get", "/dummies", correlationId, params, nil)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.DummyDataPage)
	c.Instrument(correlationId, "dummy.get_page_by_filter")
	return result, nil
}

func (c *DummyRestClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call(dummyType, "get", "/dummies/"+dummyId, correlationId, nil, nil)

	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.get_one_by_id")
	return result, nil
}

func (c *DummyRestClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call(dummyType, "post", "/dummies", correlationId, nil, dummy)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.create")
	return result, nil
}

func (c *DummyRestClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call(dummyType, "put", "/dummies", correlationId, nil, dummy)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.update")
	return result, nil
}

func (c *DummyRestClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call(dummyType, "delete", "/dummies/"+dummyId, correlationId, nil, nil)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.delete_by_id")
	return result, nil
}

func (c *DummyRestClient) CheckCorrelationId(correlationId string) (result map[string]string, err error) {

	calValue, calErr := c.Call(reflect.TypeOf(make(map[string]string)), "get", "/dummies/check/correlation_id", correlationId, nil, nil)
	if calErr != nil {
		return nil, calErr
	}

	val, _ := calValue.(*(map[string]string))
	c.Instrument(correlationId, "dummy.check_correlation_id")
	return *val, nil
}
