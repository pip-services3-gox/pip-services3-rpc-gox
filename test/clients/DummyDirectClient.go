package test_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	tlogic "github.com/pip-services3-go/pip-services3-rpc-go/test/logic"
)

type DummyDirectClient struct {
	clients.DirectClient
	specificController tlogic.IDummyController
}

func NewDummyDirectClient() *DummyDirectClient {
	ddc := DummyDirectClient{}
	ddc.DirectClient = *clients.NewDirectClient()
	ddc.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "*", "*", "*"))
	return &ddc
}

func (c *DummyDirectClient) SetReferences(references cref.IReferences) {
	c.DirectClient.SetReferences(references)

	specificController, ok := c.Controller.(tlogic.IDummyController)
	if !ok {
		panic("DummyDirectClient: Cant't resolv dependency 'controller' to IDummyController")
	}
	c.specificController = specificController

}

func (c *DummyDirectClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error) {

	timing := c.Instrument(correlationId, "dummy.get_page_by_filter")
	result, err = c.specificController.GetPageByFilter(correlationId, filter, paging)
	timing.EndTiming(err)
	return result, err

}

func (c *DummyDirectClient) GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.get_one_by_id")
	result, err = c.specificController.GetOneById(correlationId, dummyId)
	timing.EndTiming(err)
	return result, err
}

func (c *DummyDirectClient) CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.create")
	result, err = c.specificController.Create(correlationId, dummy)
	timing.EndTiming(err)
	return result, err
}

func (c *DummyDirectClient) UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.update")
	result, err = c.specificController.Update(correlationId, dummy)
	timing.EndTiming(err)
	return result, err
}

func (c *DummyDirectClient) DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.delete_by_id")
	result, err = c.specificController.DeleteById(correlationId, dummyId)
	timing.EndTiming(err)
	return result, err
}

func (c *DummyDirectClient) CheckCorrelationId(correlationId string) (result map[string]string, err error) {

	timing := c.Instrument(correlationId, "dummy.delete_by_id")
	result, err = c.specificController.CheckCorrelationId(correlationId)
	timing.EndTiming(err)
	return result, err
}
