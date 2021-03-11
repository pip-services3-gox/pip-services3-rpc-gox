package test_rpc_services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	testrpc "github.com/pip-services3-gox/pip-services3-rpc-gox/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyRestService(t *testing.T) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
	)

	var _dummy1 testrpc.Dummy
	var _dummy2 testrpc.Dummy
	var service *DummyRestService
	ctrl := testrpc.NewDummyController()

	service = NewDummyRestService()
	service.Configure(restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(references)
	opnErr := service.Open("")
	assert.Nil(t, opnErr)
	defer service.Close("")

	url := "http://localhost:3000"

	_dummy1 = testrpc.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 = testrpc.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	var dummy1 testrpc.Dummy

	// Create one dummy
	jsonBody, _ := json.Marshal(_dummy1)

	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/dummies", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr := ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	postResponse.Body.Close()
	var dummy testrpc.Dummy
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

	var dummies testrpc.DummyDataPage
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

	// Get OpenApi Spec From String
	// -----------------------------------------------------------------
	getResponse, getErr = http.Get(url + "/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	getResponse.Body.Close()

	var openApiContent = restConfig.GetAsString("openapi_content")
	assert.Equal(t, openApiContent, (string)(resBody))

	//Get OpenApi Spec From File
	// -----------------------------------------------------------------
	openApiContent = "swagger yaml content from file"
	filename := "dummy_" + cdata.IdGenerator.NextLong() + ".tmp"

	err := service.Close("")
	assert.Nil(t, err)
	// create temp file

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	assert.Nil(t, err)
	_, err = file.Write(([]byte)(openApiContent))
	assert.Nil(t, err)

	// recreate service with new configuration
	serviceConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3000,
		"openapi_file", filename, // for test only
		"swagger.enable", "true",
	)

	service.Configure(serviceConfig)
	service.Open("")

	getResponse, getErr = http.Get(url + "/swagger")
	assert.Nil(t, getErr)
	resBody, bodyErr = ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	assert.Equal(t, openApiContent, (string)(resBody))

	// delete temp file
	err = os.Remove(filename)
	assert.Nil(t, err)

}
