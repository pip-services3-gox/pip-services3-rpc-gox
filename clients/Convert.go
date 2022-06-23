package clients

import (
	"encoding/json"
	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	"io/ioutil"
	"net/http"
	"reflect"
)

// ConvertCommandResult method helps get correct result from JSON by prototype
//	Parameters:
//		- comRes any  input JSON string
//		- prototype reflect.Type output object prototype
//	Returns: convRes any, err error
func ConvertCommandResult(comRes any, prototype reflect.Type) (convRes any, err error) {

	str, ok := comRes.([]byte)
	if !ok || string(str) == "null" {
		return nil, nil
	}

	if prototype.Kind() == reflect.Ptr {
		prototype = prototype.Elem()
	}
	convRes = reflect.New(prototype).Interface()

	convErr := json.Unmarshal(comRes.([]byte), &convRes)
	if convErr != nil {
		return nil, convErr
	}

	return convRes, nil
}

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
