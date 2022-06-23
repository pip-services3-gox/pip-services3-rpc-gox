package test_clients

import (
	"context"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
	tlogic "github.com/pip-services3-gox/pip-services3-rpc-gox/test/logic"
)

type DummyDirectClient struct {
	clients.DirectClient
	specificController tlogic.IDummyController
}

func NewDummyDirectClient() *DummyDirectClient {
	ddc := DummyDirectClient{}
	ddc.DirectClient = *clients.NewDirectClient()
	ddc.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "*", "*", "*"))
	return &ddc
}

func (c *DummyDirectClient) SetReferences(ctx context.Context, references cref.IReferences) {
	c.DirectClient.SetReferences(ctx, references)

	specificController, ok := c.Controller.(tlogic.IDummyController)
	if !ok {
		panic("DummyDirectClient: Cant't resolv dependency 'controller' to IDummyController")
	}
	c.specificController = specificController

}

func (c *DummyDirectClient) GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (cdata.DataPage[tdata.Dummy], error) {

	timing := c.Instrument(ctx, correlationId, "dummy.get_page_by_filter")
	result, err := c.specificController.GetPageByFilter(ctx, correlationId, &filter, &paging)
	timing.EndTiming(ctx, err)
	return *result, err

}

func (c *DummyDirectClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (tdata.Dummy, error) {

	timing := c.Instrument(ctx, correlationId, "dummy.get_one_by_id")
	result, err := c.specificController.GetOneById(ctx, correlationId, dummyId)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (tdata.Dummy, error) {

	timing := c.Instrument(ctx, correlationId, "dummy.create")
	result, err := c.specificController.Create(ctx, correlationId, dummy)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (tdata.Dummy, error) {

	timing := c.Instrument(ctx, correlationId, "dummy.update")
	result, err := c.specificController.Update(ctx, correlationId, dummy)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (tdata.Dummy, error) {

	timing := c.Instrument(ctx, correlationId, "dummy.delete_by_id")
	result, err := c.specificController.DeleteById(ctx, correlationId, dummyId)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CheckCorrelationId(ctx context.Context, correlationId string) (map[string]string, error) {

	timing := c.Instrument(ctx, correlationId, "dummy.delete_by_id")
	result, err := c.specificController.CheckCorrelationId(ctx, correlationId)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CheckErrorPropagation(ctx context.Context, correlationId string) error {
	timing := c.Instrument(ctx, correlationId, "dummy.check_error_propagation")
	err := c.specificController.CheckErrorPropagation(ctx, correlationId)
	timing.EndTiming(ctx, err)
	return err
}
