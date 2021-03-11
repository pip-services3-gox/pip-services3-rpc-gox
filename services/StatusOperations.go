package services

import (
	"net/http"
	"time"

	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
)

/*
StatusOperations helper class for status service
*/
type StatusOperations struct {
	RestOperations
	startTime   time.Time
	references2 crefer.IReferences
	contextInfo *cinfo.ContextInfo
}

// NewStatusOperations creates new instance of StatusOperations
func NewStatusOperations() *StatusOperations {
	c := StatusOperations{}
	c.startTime = time.Now()
	c.DependencyResolver.Put("context-info", crefer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	return &c
}

// SetReferences  sets references to dependent components.
//  - references  crefer.IReferences	references to locate the component dependencies.
func (c *StatusOperations) SetReferences(references crefer.IReferences) {
	c.references2 = references
	c.RestOperations.SetReferences(references)

	depRes := c.DependencyResolver.GetOneOptional("context-info")
	if depRes != nil {
		c.contextInfo = depRes.(*cinfo.ContextInfo)
	}
}

// GetStatusOperation return function for get status
func (c *StatusOperations) GetStatusOperation() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		c.Status(res, req)
	}
}

// Status method handles status requests
//   - req *http.Request  an HTTP request
//   - res  http.ResponseWriter  an HTTP response
func (c *StatusOperations) Status(res http.ResponseWriter, req *http.Request) {

	id := ""
	if c.contextInfo != nil {
		id = c.contextInfo.ContextId
	}

	name := "Unknown"
	if c.contextInfo != nil {
		name = c.contextInfo.Name
	}

	description := ""
	if c.contextInfo != nil {
		description = c.contextInfo.Description
	}

	uptime := time.Now().Sub(c.startTime)

	properties := make(map[string]string)
	if c.contextInfo != nil {
		properties = c.contextInfo.Properties
	}

	var components []string
	if c.references2 != nil {
		for _, locator := range c.references2.GetAllLocators() {
			components = append(components, cconv.StringConverter.ToString(locator))
		}
	}

	status := make(map[string]interface{})

	status["id"] = id
	status["name"] = name
	status["description"] = description
	status["start_time"] = cconv.StringConverter.ToString(c.startTime)
	status["current_time"] = cconv.StringConverter.ToString(time.Now())
	status["uptime"] = uptime
	status["properties"] = properties
	status["components"] = components

	c.SendResult(res, req, status, nil)
}
