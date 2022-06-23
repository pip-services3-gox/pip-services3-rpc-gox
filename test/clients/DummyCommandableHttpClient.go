package test_clients

import (
	"context"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
)

type DummyCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewDummyCommandableHttpClient() *DummyCommandableHttpClient {
	dchc := DummyCommandableHttpClient{}
	dchc.CommandableHttpClient = *clients.NewCommandableHttpClient("dummies")
	return &dchc
}

func (c *DummyCommandableHttpClient) GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (result cdata.DataPage[tdata.Dummy], err error) {
	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "get_dummies", correlationId, cdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return *cdata.NewEmptyDataPage[tdata.Dummy](), err
	}

	return clients.HandleHttpResponse[cdata.DataPage[tdata.Dummy]](response, correlationId)
}

func (c *DummyCommandableHttpClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (tdata.Dummy, error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "get_dummy_by_id", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableHttpClient) CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "create_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableHttpClient) UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "update_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableHttpClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "delete_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableHttpClient) CheckCorrelationId(ctx context.Context, correlationId string) (result map[string]string, err error) {

	params := cdata.NewEmptyAnyValueMap()

	response, err := c.CallCommand(ctx, "check_correlation_id", correlationId, params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[map[string]string](response, correlationId)
}

func (c *DummyCommandableHttpClient) CheckErrorPropagation(ctx context.Context, correlationId string) error {
	params := cdata.NewEmptyAnyValueMap()
	_, calErr := c.CallCommand(ctx, "check_error_propagation", correlationId, params)
	return calErr
}
