package test_services

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cinfo "github.com/pip-services3-gox/pip-services3-components-gox/info"
	"github.com/pip-services3-gox/pip-services3-rpc-gox/services"
	tlogic "github.com/pip-services3-gox/pip-services3-rpc-gox/test/logic"
)

const (
	StatusRestServicePort = iota + 3000
	HeartbeatRestServicePort
	HttpEndpointServicePort
	DummyRestServicePort
	DummyOpenAPIFileRestServicePort
	DummyCommandableHttpServicePort
	DummyCommandableSwaggerHttpServicePort
)

func TestMain(m *testing.M) {

	fmt.Println("Preparing test services...")

	statusRestService := BuildTestStatusRestService()
	err := statusRestService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer statusRestService.Close(context.Background(), "")

	heartbeatRestService := BuildTestHeartbeatRestService()
	err = heartbeatRestService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer heartbeatRestService.Close(context.Background(), "")

	httpEndpointService, endpoint := BuildTestHttpEndpointService()
	err = endpoint.Open(context.Background(), "")
	if err != nil {
		panic(err)
	} else {
		err = httpEndpointService.Open(context.Background(), "")
		if err != nil {
			panic(err)
		} else {
			defer endpoint.Close(context.Background(), "")
			defer httpEndpointService.Close(context.Background(), "")
		}
	}

	dummyRestService := BuildTestDummyRestService()
	err = dummyRestService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyRestService.Close(context.Background(), "")

	dummyOpenAPIFileRestService, filename := BuildTestDummyOpenAPIFileRestService()
	err = dummyOpenAPIFileRestService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyOpenAPIFileRestService.Close(context.Background(), "")
	//defer os.Remove(filename)
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()

	dummyCommandableHttpService := BuildTestDummyCommandableHttpService()
	err = dummyCommandableHttpService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyCommandableHttpService.Close(context.Background(), "")

	dummyCommandableSwaggerHttpService := BuildTestDummyCommandableSwaggerHttpService()
	err = dummyCommandableSwaggerHttpService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyCommandableSwaggerHttpService.Close(context.Background(), "")
	time.Sleep(time.Second)
	fmt.Println("All test services started!")

	code := m.Run()

	noc := dummyRestService.GetNumberOfCalls()
	fmt.Println("Number of calls:", noc, "from 4")
	if noc != 4 {
		panic("Number of calls test failed!")
	}
	os.Exit(code)
}

func BuildTestStatusRestService() *services.StatusRestService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", StatusRestServicePort,
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	service := services.NewStatusRestService()
	service.Configure(context.Background(), restConfig)

	contextInfo := cinfo.NewContextInfo()
	contextInfo.Name = "Test"
	contextInfo.Description = "This is a test container"

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "status-service", "http", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}

func BuildTestHttpEndpointService() (*DummyRestService, *services.HttpEndpoint) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", HttpEndpointServicePort,
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	ctrl := tlogic.NewDummyController()
	service := NewDummyRestService()
	service.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"base_route",
		"/api/v1",
	))

	endpoint := services.NewHttpEndpoint()
	endpoint.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
		cref.NewDescriptor("pip-services", "endpoint", "http", "default", "1.0"), endpoint,
	)
	service.SetReferences(context.Background(), references)
	return service, endpoint
}

func BuildTestDummyRestService() *DummyRestService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestServicePort,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	var service *DummyRestService
	ctrl := tlogic.NewDummyController()

	service = NewDummyRestService()
	service.Configure(context.Background(), restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}

func BuildTestDummyOpenAPIFileRestService() (*DummyRestService, string) {

	openApiContent := "swagger yaml content from file"
	filename := path.Join(".", "dummy_"+cdata.IdGenerator.NextLong()+".tmp")

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(([]byte)(openApiContent))
	if err != nil {
		panic(err)
	}
	//err = file.Close()
	//if err != nil {
	//	panic(err)
	//}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyOpenAPIFileRestServicePort,
		"openapi_file", filename, // for test only
		"swagger.enable", "true",
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	var service *DummyRestService
	ctrl := tlogic.NewDummyController()

	service = NewDummyRestService()
	service.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service, filename
}

func BuildTestDummyCommandableHttpService() *DummyCommandableHttpService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpServicePort,
		"swagger.enable", "true",
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	ctrl := tlogic.NewDummyController()

	service := NewDummyCommandableHttpService()

	service.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}

func BuildTestDummyCommandableSwaggerHttpService() *DummyCommandableHttpService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableSwaggerHttpServicePort,
		"swagger.enable", "true",
		"swagger.auto", false,
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	ctrl := tlogic.NewDummyController()

	service := NewDummyCommandableHttpService()

	service.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}

func BuildTestHeartbeatRestService() *services.HeartbeatRestService {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", HeartbeatRestServicePort,
		"cors_headers", "correlation_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	service := services.NewHeartbeatRestService()
	service.Configure(context.Background(), restConfig)
	return service
}
