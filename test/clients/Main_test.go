package test_clients

import (
	"context"
	"fmt"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	tlogic "github.com/pip-services3-gox/pip-services3-rpc-gox/test/logic"
	test_services "github.com/pip-services3-gox/pip-services3-rpc-gox/test/services"
	"os"
	"testing"
	"time"
)

const (
	DummyRestServicePort = iota + 4000
	DummyCommandableHttpServicePort
)

func TestMain(m *testing.M) {

	fmt.Println("Preparing test services for clients...")

	dummyRestService := BuildTestDummyRestService()
	err := dummyRestService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyRestService.Close(context.Background(), "")

	dummyCommandableHttpService := BuildTestDummyCommandableHttpService()
	err = dummyCommandableHttpService.Open(context.Background(), "")
	if err != nil {
		panic(err)
	}
	defer dummyCommandableHttpService.Close(context.Background(), "")
	time.Sleep(time.Second)
	fmt.Println("All test services started!")

	os.Exit(m.Run())
}

func BuildTestDummyRestService() *test_services.DummyRestService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestServicePort,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
	)

	var service *test_services.DummyRestService
	ctrl := tlogic.NewDummyController()

	service = test_services.NewDummyRestService()
	service.Configure(context.Background(), restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}

func BuildTestDummyCommandableHttpService() *test_services.DummyCommandableHttpService {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpServicePort,
		"swagger.enable", "true",
	)

	ctrl := tlogic.NewDummyController()

	service := test_services.NewDummyCommandableHttpService()

	service.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
	)
	service.SetReferences(context.Background(), references)
	return service
}
