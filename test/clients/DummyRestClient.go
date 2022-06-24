package test_clients

import (
	"context"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
	"net/http"
)

type DummyRestClient struct {
	clients.RestClient
}

func NewDummyRestClient() *DummyRestClient {
	drc := DummyRestClient{}
	drc.RestClient = *clients.NewRestClient()
	return &drc
}

func (c *DummyRestClient) GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams,
	paging cdata.PagingParams) (result cdata.DataPage[tdata.Dummy], err error) {

	defer c.Instrument(ctx, correlationId, "dummy.get_page_by_filter")

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.Call(ctx, http.MethodGet, "/dummies", correlationId, params, nil)
	if err != nil {
		return *cdata.NewEmptyDataPage[tdata.Dummy](), err
	}

	return clients.HandleHttpResponse[cdata.DataPage[tdata.Dummy]](response, correlationId)
}

func (c *DummyRestClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {

	defer c.Instrument(ctx, correlationId, "dummy.get_one_by_id")

	response, err := c.Call(ctx, http.MethodGet, "/dummies/"+dummyId, correlationId, nil, nil)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyRestClient) CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {

	defer c.Instrument(ctx, correlationId, "dummy.create")

	response, err := c.Call(ctx, http.MethodPost, "/dummies", correlationId, nil, dummy)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyRestClient) UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {

	defer c.Instrument(ctx, correlationId, "dummy.update")

	response, err := c.Call(ctx, http.MethodPut, "/dummies", correlationId, nil, dummy)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyRestClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {

	defer c.Instrument(ctx, correlationId, "dummy.delete_by_id")

	response, err := c.Call(ctx, http.MethodDelete, "/dummies/"+dummyId, correlationId, nil, nil)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return clients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyRestClient) CheckCorrelationId(ctx context.Context, correlationId string) (result map[string]string, err error) {

	defer c.Instrument(ctx, correlationId, "dummy.check_correlation_id")

	response, err := c.Call(ctx, http.MethodGet, "/dummies/check/correlation_id", correlationId, nil, nil)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[map[string]string](response, correlationId)
}

func (c *DummyRestClient) CheckErrorPropagation(ctx context.Context, correlationId string) error {

	c.Instrument(ctx, correlationId, "dummy.check_error_propagation")

	_, err := c.Call(ctx, http.MethodGet, "/dummies/check/error_propagation", correlationId, nil, nil)
	return err
}
