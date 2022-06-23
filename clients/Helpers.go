package clients

import (
	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	"io/ioutil"
	"net/http"
)

// HandleHttpResponse method helps handle http response body
//	Parameters:
//		- ctx context.Context
//		- correlationId string (optional) transaction id to trace execution through call chain.
//	Returns: T any result, err error
func HandleHttpResponse[T any](r *http.Response, correlationId string) (T, error) {
	defer r.Body.Close()

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		var defaultValue T
		return defaultValue, cerr.ApplicationErrorFactory.
			Create(&cerr.ErrorDescription{
				Type:          "Application",
				Category:      "Application",
				Status:        r.StatusCode,
				Code:          "",
				Message:       err.Error(),
				CorrelationId: correlationId,
			}).
			WithCause(err)
	}

	return convert.NewDefaultCustomTypeJsonConvertor[T]().FromJson(string(buffer))
}
