package test_services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	"github.com/pip-services3-go/pip-services3-rpc-go/services"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	tlogic "github.com/pip-services3-go/pip-services3-rpc-go/test/logic"
)

type DummyRestService struct {
	*services.RestService
	controller     tlogic.IDummyController
	numberOfCalls  int
	openApiContent string
	openApiFile    string
}

func NewDummyRestService() *DummyRestService {
	c := &DummyRestService{}
	c.RestService = services.InheritRestService(c)
	c.numberOfCalls = 0
	c.DependencyResolver.Put("controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyRestService) Configure(config *cconf.ConfigParams) {
	c.openApiContent = *config.GetAsNullableString("openapi_content")
	c.openApiFile = *config.GetAsNullableString("openapi_file")
	c.RestService.Configure(config)
}

func (c *DummyRestService) SetReferences(references crefer.IReferences) {
	c.RestService.SetReferences(references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(tlogic.IDummyController)
	}
}

func (c *DummyRestService) GetNumberOfCalls() int {
	return c.numberOfCalls
}

func (c *DummyRestService) incrementNumberOfCalls(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	c.numberOfCalls++
	next.ServeHTTP(res, req)
}

func (c *DummyRestService) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	result, err := c.controller.GetPageByFilter(
		c.GetCorrelationId(req),
		cdata.NewFilterParamsFromValue(params), // W! need test
		cdata.NewPagingParamsFromTuples(paginParams),
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) getOneById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}
	result, err := c.controller.GetOneById(
		c.GetCorrelationId(req),
		dummyId)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) create(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)
	var dummy tdata.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}

	result, err := c.controller.Create(
		correlationId,
		dummy,
	)
	c.SendCreatedResult(res, req, result, err)
}

func (c *DummyRestService) update(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)

	var dummy tdata.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}
	result, err := c.controller.Update(
		correlationId,
		dummy,
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) deleteById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}

	result, err := c.controller.DeleteById(
		c.GetCorrelationId(req),
		dummyId,
	)
	c.SendDeletedResult(res, req, result, err)
}

func (c *DummyRestService) checkCorrelationId(res http.ResponseWriter, req *http.Request) {
	result, err := c.controller.CheckCorrelationId(c.GetCorrelationId(req))
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) checkErrorPropagation(res http.ResponseWriter, req *http.Request) {
	err := c.controller.CheckErrorPropagation(c.GetCorrelationId(req))
	c.SendError(res, req, err)
}

func (c *DummyRestService) Register() {
	c.RegisterInterceptor("/dummies$", c.incrementNumberOfCalls)

	c.RegisterRoute(
		"get", "/dummies",
		&cvalid.NewObjectSchema().WithOptionalProperty("skip", cconv.String).
			WithOptionalProperty("take", cconv.String).
			WithOptionalProperty("total", cconv.String).
			WithOptionalProperty("body", cvalid.NewFilterParamsSchema()).Schema,
		c.getPageByFilter,
	)

	c.RegisterRoute(
		"get", "/dummies/check/correlation_id",
		&cvalid.NewObjectSchema().Schema,
		c.checkCorrelationId,
	)

	c.RegisterRoute(
		"get", "/dummies/check/error_propagation",
		&cvalid.NewObjectSchema().Schema,
		c.checkErrorPropagation,
	)

	c.RegisterRoute(
		"get", "/dummies/{dummy_id}",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.getOneById,
	)

	c.RegisterRoute(
		"post", "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.create,
	)

	c.RegisterRoute(
		"put", "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.update,
	)

	c.RegisterRoute(
		"delete", "/dummies/{dummy_id}",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.deleteById,
	)

	if c.openApiContent != "" {
		c.RegisterOpenApiSpec(c.openApiContent)
	}

	if c.openApiFile != "" {
		c.RegisterOpenApiSpecFromFile(c.openApiFile)
	}
}
