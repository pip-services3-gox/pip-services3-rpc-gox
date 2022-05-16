package test_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommandableHttpService(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", DummyCommandableHttpServicePort)

	_dummy1 := tdata.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 := tdata.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Create one dummy

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = _dummy1

	jsonBody, _ := json.Marshal(bodyMap)

	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/dummies/create_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr := ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()

	var dummy tdata.Dummy
	jsonErr := json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)

	dummy1 := dummy

	// Create another dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy"] = _dummy2

	jsonBody, _ = json.Marshal(bodyMap)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/create_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()

	jsonErr = json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)

	// Get all dummies

	postResponse, postErr = http.Post(url+"/dummies/get_dummies", "application/json", nil)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()
	var dummies tdata.DummyDataPage
	jsonErr = json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	bodyMap = make(map[string]interface{})
	bodyMap["dummy"] = dummy1

	jsonBody, _ = json.Marshal(bodyMap)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/update_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()
	jsonErr = json.Unmarshal(resBody, &dummy)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)

	// Delete dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy_id"] = dummy1.Id
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/delete_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	postResponse.Body.Close()
	assert.Nil(t, bodyErr)

	// Try to get delete dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy_id"] = dummy1.Id
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/get_dummy_by_id", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()
	dummy = tdata.Dummy{}
	jsonErr = json.Unmarshal(resBody, &dummy)
	assert.Nil(t, jsonErr)
	assert.Empty(t, dummy)

	// Testing transmit correlationId
	bodyMap = make(map[string]interface{})
	bodyMap["dummy_id"] = dummy1.Id
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	getResponse, getErr := http.Post(url+"/dummies/check_correlation_id?correlation_id=test_cor_id", "application/json", bodyReader)
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()
	values := make(map[string]string, 0)
	jsonErr = json.Unmarshal(resBody, &values)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, values)
	assert.Equal(t, values["correlationId"], "test_cor_id")

	req, reqErr := http.NewRequest("POST", url+"/dummies/check_correlation_id", bytes.NewBuffer(make([]byte, 0, 0)))
	assert.Nil(t, reqErr)
	req.Header.Set("correlation_id", "test_cor_id")
	localClient := http.Client{}
	getResponse, getErr = localClient.Do(req)
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()
	values = make(map[string]string, 0)
	jsonErr = json.Unmarshal(resBody, &values)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, values)
	assert.Equal(t, values["correlationId"], "test_cor_id")

	// Testing error propagation
	getResponse, getErr = http.Post(url+"/dummies/check_error_propagation?correlation_id=test_error_propagation", "application/json", nil)
	assert.Nil(t, getErr)
	assert.NotNil(t, getResponse)

	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()

	appErr := cerr.ApplicationError{}
	jsonErr = json.Unmarshal(resBody, &appErr)
	assert.Nil(t, jsonErr)

	assert.Equal(t, appErr.CorrelationId, "test_error_propagation")
	assert.Equal(t, appErr.Status, 404)
	assert.Equal(t, appErr.Code, "NOT_FOUND_TEST")
	assert.Equal(t, appErr.Message, "Not found error")

	// Get OpenApi Spec From String
	// -----------------------------------------------------------------
	getResponse, getErr = http.Get(url + "/dummies/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	fmt.Println((string)(resBody))
	assert.True(t, strings.Index((string)(resBody), "openapi:") >= 0)

	// Get OpenApi Spec From File
	// -----------------------------------------------------------------
	url = fmt.Sprintf("http://localhost:%d", DummyCommandableSwaggerHttpServicePort)
	getResponse, getErr = http.Get(url + "/dummies/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	assert.Equal(t, "swagger yaml content", (string)(resBody))
}
