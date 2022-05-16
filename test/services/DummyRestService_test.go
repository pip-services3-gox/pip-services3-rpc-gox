package test_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	tdata "github.com/pip-services3-go/pip-services3-rpc-go/test/data"
	"github.com/stretchr/testify/assert"
)

func TestDummyRestService(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", DummyRestServicePort)

	_dummy1 := tdata.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 := tdata.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	var dummy1 tdata.Dummy

	// Create one dummy
	jsonBody, _ := json.Marshal(_dummy1)

	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/dummies", "application/json", bodyReader)
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

	dummy1 = dummy

	// Create another dummy
	jsonBody, _ = json.Marshal(_dummy2)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()

	jsonErr = json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)
	//dummy2 = dummy

	// Get all dummies
	getResponse, getErr := http.Get(url + "/dummies")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()

	var dummies tdata.DummyDataPage
	jsonErr = json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy

	dummy1.Content = "Updated Content 1"
	jsonBody, _ = json.Marshal(dummy1)

	client := &http.Client{}
	data := bytes.NewReader(jsonBody)
	putReq, putErr := http.NewRequest(http.MethodPut, url+"/dummies", data)
	assert.Nil(t, putErr)
	putRes, putErr := client.Do(putReq)
	assert.Nil(t, putErr)
	resBody, bodyErr = ioutil.ReadAll(putRes.Body)
	putRes.Body.Close()
	jsonErr = json.Unmarshal(resBody, &dummy)
	assert.Nil(t, putErr)
	assert.NotNil(t, dummy)

	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)
	dummy1 = dummy

	// Delete dummy
	delReq, delErr := http.NewRequest(http.MethodDelete, url+"/dummies/"+dummy1.Id, nil)
	assert.Nil(t, delErr)
	_, delErr = client.Do(delReq)
	assert.Nil(t, delErr)

	// Try to get delete dummy
	dummies.Data = dummies.Data[:0]
	*dummies.Total = 0
	getResponse, getErr = http.Get(url + "/dummies/" + dummy1.Id)
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()
	jsonErr = json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 0)

	// Testing transmit correlationId
	getResponse, getErr = http.Get(url + "/dummies/check/correlation_id?correlation_id=test_cor_id")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()
	values := make(map[string]string, 0)
	jsonErr = json.Unmarshal(resBody, &values)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, values)
	assert.Equal(t, values["correlationId"], "test_cor_id")

	req, reqErr := http.NewRequest("GET", url+"/dummies/check/correlation_id", bytes.NewBuffer(make([]byte, 0, 0)))
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
	getResponse, getErr = http.Get(url + "/dummies/check/error_propagation?correlation_id=test_error_propagation")
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
	getResponse, getErr = http.Get(url + "/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()

	assert.Equal(t, "swagger yaml or json content", (string)(resBody))

	//Get OpenApi Spec From File
	// -----------------------------------------------------------------
	url = fmt.Sprintf("http://localhost:%d", DummyOpenAPIFileRestServicePort)
	getResponse, getErr = http.Get(url + "/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	assert.Equal(t, "swagger yaml content from file", (string)(resBody))
}
