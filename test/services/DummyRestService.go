package test_services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
	tlogic "github.com/pip-services3-gox/pip-services3-rpc-gox/test/logic"
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
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyRestService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	if _val, ok := config.GetAsNullableString("openapi_content"); ok {
		c.openApiContent = _val
	}
	if _val, ok := config.GetAsNullableString("openapi_file"); ok {
		c.openApiFile = _val
	}
	c.RestService.Configure(ctx, config)
}

func (c *DummyRestService) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.RestService.SetReferences(ctx, references)
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
		req.Context(),
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
		req.Context(),
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
		req.Context(),
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
		req.Context(),
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
		req.Context(),
		c.GetCorrelationId(req),
		dummyId,
	)
	c.SendDeletedResult(res, req, result, err)
}

func (c *DummyRestService) checkCorrelationId(res http.ResponseWriter, req *http.Request) {
	result, err := c.controller.CheckCorrelationId(req.Context(), c.GetCorrelationId(req))
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) checkErrorPropagation(res http.ResponseWriter, req *http.Request) {
	err := c.controller.CheckErrorPropagation(req.Context(), c.GetCorrelationId(req))
	c.SendError(res, req, err)
}

func (c *DummyRestService) checkGracefulShutdownContext(res http.ResponseWriter, req *http.Request) {
	err := c.controller.CheckGracefulShutdownContext(req.Context(), c.GetCorrelationId(req))
	c.SendError(res, req, err)
}

func (c *DummyRestService) Register() {
	c.RegisterInterceptor("/dummies$", c.incrementNumberOfCalls)

	c.RegisterRoute(
		http.MethodGet, "/dummies",
		&cvalid.NewObjectSchema().WithOptionalProperty("skip", cconv.String).
			WithOptionalProperty("take", cconv.String).
			WithOptionalProperty("total", cconv.String).
			WithOptionalProperty("body", cvalid.NewFilterParamsSchema()).Schema,
		c.getPageByFilter,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/correlation_id",
		&cvalid.NewObjectSchema().Schema,
		c.checkCorrelationId,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/error_propagation",
		&cvalid.NewObjectSchema().Schema,
		c.checkErrorPropagation,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/check/graceful_shutdown",
		&cvalid.NewObjectSchema().Schema,
		c.checkGracefulShutdownContext,
	)

	c.RegisterRoute(
		http.MethodGet, "/dummies/{dummy_id}",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.getOneById,
	)

	c.RegisterRoute(
		http.MethodPost, "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.create,
	)

	c.RegisterRoute(
		http.MethodPut, "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", tdata.NewDummySchema()).Schema,
		c.update,
	)

	c.RegisterRoute(
		http.MethodDelete, "/dummies/{dummy_id}",
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
