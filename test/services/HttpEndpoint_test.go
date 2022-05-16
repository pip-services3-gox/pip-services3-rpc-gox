package test_services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	"github.com/stretchr/testify/assert"
)

func TestHttpEndpoint(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", HttpEndpointServicePort)

	getResponse, getErr := http.Get(url + "/api/v1/dummies")
	assert.Nil(t, getErr)
	resBody, bodyErr := ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	var dummies tdata.DummyDataPage
	jsonErr := json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 0)
}
