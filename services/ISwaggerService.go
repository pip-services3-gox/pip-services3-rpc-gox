package services

// ISwaggerService Interface to perform Swagger registrations.
type ISwaggerService interface {

	//  Perform required Swagger registration steps.
	RegisterOpenApiSpec(baseRoute string, swaggerRoute string)
}
