package test_services

import (
	"encoding/json"
	"fmt"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"io/ioutil"
	"net/http"
	"testing"

	tdata "github.com/pip-services3-gox/pip-services3-rpc-gox/test/data"
	"github.com/stretchr/testify/assert"
)

func TestHttpEndpoint(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", HttpEndpointServicePort)

	getResponse, getErr := http.Get(url + "/api/v1/dummies")
	assert.Nil(t, getErr)
	resBody, bodyErr := ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	var dummies *cdata.DataPage[tdata.Dummy]
	jsonErr := json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.False(t, dummies.HasData())
	assert.Len(t, dummies.Data, 0)
}
