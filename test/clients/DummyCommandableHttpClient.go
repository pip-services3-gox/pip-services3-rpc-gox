package test_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
)

type DummyCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewDummyCommandableHttpClient() *DummyCommandableHttpClient {
	dchc := DummyCommandableHttpClient{}
	dchc.CommandableHttpClient = *clients.NewCommandableHttpClient("dummies")
	return &dchc
}

func (c *DummyCommandableHttpClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.CallCommand(dummyDataPageType, "get_dummies", correlationId, cdata.NewAnyValueMapFromValue(params.Value()))
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.DummyDataPage)
	return result, err
}

func (c *DummyCommandableHttpClient) GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "get_dummy_by_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "create_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "update_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "delete_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) CheckCorrelationId(correlationId string) (result map[string]string, err error) {

	params := cdata.NewEmptyAnyValueMap()

	calValue, calErr := c.CallCommand(reflect.TypeOf(make(map[string]string)), "check_correlation_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	val, _ := calValue.(*(map[string]string))
	return *val, err
}

func (c *DummyCommandableHttpClient) CheckErrorPropagation(correlationId string) error {
	params := cdata.NewEmptyAnyValueMap()
	_, calErr := c.CallCommand(nil, "check_error_propagation", correlationId, params)
	return calErr
}
