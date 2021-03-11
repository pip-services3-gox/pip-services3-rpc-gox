package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
RestOperations helper class for REST operations
*/
type RestOperations struct {
	Logger             *clog.CompositeLogger
	Counters           *ccount.CompositeCounters
	DependencyResolver *crefer.DependencyResolver
}

// NewRestOperations creates new instance of RestOperations
func NewRestOperations() *RestOperations {
	ro := RestOperations{}
	ro.Logger = clog.NewCompositeLogger()
	ro.Counters = ccount.NewCompositeCounters()
	ro.DependencyResolver = crefer.NewDependencyResolver()
	return &ro
}

// Configure method are configures this RestOperations using the given configuration parameters.
// Parameters:
//   - config *cconf.ConfigParams confif parameters
func (c *RestOperations) Configure(config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(config)
}

// SetReferences method are sets references to this RestOperations logger, counters, and connection resolver.
// Parameters:
//   - references    an IReferences object, containing references to a logger, counters,
//   and a dependency resolver.
func (c *RestOperations) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
}

// GetCorrelationId method returns CorrelationId from request
// Parameters:
//   req *http.Request  request
// Returns: string
// retrun correlation_id or empty string
func (c *RestOperations) GetCorrelationId(req *http.Request) string {
	correlationId := req.URL.Query().Get("correlation_id")
	if correlationId == "" {
		correlationId = req.Header.Get("correlation_id")
	}
	return correlationId
}

// GetFilterParams method retruns filter params object from request
// Parameters:
//   req *http.Request  request
// Returns: *cdata.FilterParams
// filter params object
func (c *RestOperations) GetFilterParams(req *http.Request) *cdata.FilterParams {

	params := req.URL.Query()
	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")
	filter := cdata.NewFilterParamsFromValue(
		params,
	)
	return filter
}

// GetPagingParams method retruns paging params object from request
// Parameters:
//   req *http.Request  request
// Returns: *cdata.PagingParams
// pagings params object
func (c *RestOperations) GetPagingParams(req *http.Request) *cdata.PagingParams {

	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	paging := cdata.NewPagingParamsFromValue(
		paginParams,
	)
	return paging
}

// GetParam methods helps get all params from query
//   - req   - incoming request
//   - name  - parameter name
// Returns value or empty string if param not exists
func (c *RestOperations) GetParam(req *http.Request, name string) string {
	param := req.URL.Query().Get(name)
	if param == "" {
		param = mux.Vars(req)[name]
	}
	return param
}

// DecodeBody methods helps decode body
//   - req   	- incoming request
//   - target  	- pointer on target variable for decode
// Returns error
func (c *RestOperations) DecodeBody(req *http.Request, target interface{}) error {

	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	err = json.Unmarshal(bytes, target)
	if err != nil {
		return err
	}
	return nil
}

func (c *RestOperations) SendResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendResult(res, req, result, err)
}

func (c *RestOperations) SendEmptyResult(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendEmptyResult(res, req, err)
}

func (c *RestOperations) SendCreatedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *RestOperations) SendDeletedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *RestOperations) SendError(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendError(res, req, err)
}

func (c *RestOperations) SendBadRequest(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewBadRequestError(correlationId, "BAD_REQUEST", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendUnauthorized(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewUnauthorizedError(correlationId, "UNAUTHORIZED", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendNotFound(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewNotFoundError(correlationId, "NOT_FOUND", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendConflict(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewConflictError(correlationId, "CONFLICT", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendSessionExpired(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	err := cerr.NewUnknownError(correlationId, "SESSION_EXPIRED", message)
	err.Status = 440
	c.SendError(res, req, err)
}

func (c *RestOperations) SendInternalError(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewUnknownError(correlationId, "INTERNAL", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendServerUnavailable(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	err := cerr.NewConflictError(correlationId, "SERVER_UNAVAILABLE", message)
	err.Status = 503
	c.SendError(res, req, err)
}
