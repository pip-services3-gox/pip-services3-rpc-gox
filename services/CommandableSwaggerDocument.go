package services

import (
	"strings"

	ccomands "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type CommandableSwaggerDocument struct {
	content string

	Commands []ccomands.ICommand

	Version   string
	BaseRoute string

	InfoTitle          string
	InfoDescription    string
	InfoVersion        string
	InfoTermsOfService string

	InfoContactName  string
	InfoContactUrl   string
	InfoContactEmail string

	InfoLicenseName string
	InfoLicenseUrl  string
}

func NewCommandableSwaggerDocument(baseRoute string, config *cconf.ConfigParams, commands []ccomands.ICommand) *CommandableSwaggerDocument {
	c := &CommandableSwaggerDocument{
		content:     "",
		Version:     "3.0.2",
		InfoVersion: "1",
		BaseRoute:   baseRoute,
		Commands:    make([]ccomands.ICommand, 0),
	}

	if commands != nil {
		c.Commands = commands
	}

	if config == nil {
		config = cconf.NewEmptyConfigParams()
	}

	c.InfoTitle = config.GetAsStringWithDefault("name", "CommandableHttpService")
	c.InfoDescription = config.GetAsStringWithDefault("description", "Commandable microservice")
	return c
}

func (c *CommandableSwaggerDocument) ToString() string {
	var data = map[string]any{
		"openapi": c.Version,
		"info": map[string]any{
			"title":          c.InfoTitle,
			"description":    c.InfoDescription,
			"version":        c.InfoVersion,
			"termsOfService": c.InfoTermsOfService,
			"contact": map[string]any{
				"name":  c.InfoContactName,
				"url":   c.InfoContactUrl,
				"email": c.InfoContactEmail,
			},
			"license": map[string]any{
				"name": c.InfoLicenseName,
				"url":  c.InfoLicenseUrl,
			},
		},
		"paths": c.createPathsData(),
	}

	c.WriteData(0, data)

	//console.log(c.content);
	return c.content
}

func (c *CommandableSwaggerDocument) createPathsData() map[string]any {
	var data = make(map[string]any, 0)

	for index := 0; index < len(c.Commands); index++ {
		command := c.Commands[index]

		var path = c.BaseRoute + "/" + command.Name()
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		data[path] = map[string]any{

			"post": map[string]any{
				"tags":        []any{c.BaseRoute},
				"operationId": command.Name(),
				"requestBody": c.createRequestBodyData(command),
				"responses":   c.createResponsesData(),
			},
		}
	}

	return data
}

func (c *CommandableSwaggerDocument) createRequestBodyData(command ccomands.ICommand) map[string]any {
	var schemaData = c.createSchemaData(command)
	if schemaData == nil {
		return nil
	}

	return map[string]any{
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": schemaData,
			},
		},
	}
}

func (c *CommandableSwaggerDocument) createSchemaData(command ccomands.ICommand) map[string]any {
	var schema = command.(*ccomands.Command).GetSchema().(*cvalid.ObjectSchema)

	if schema == nil || schema.Properties() == nil {
		return nil
	}

	var properties = make(map[string]any, 0)
	var required = make([]any, 0)

	for _, property := range schema.Properties() {
		tp, _ := property.Type().(cconv.TypeCode)

		properties[property.Name()] = map[string]any{

			"type": c.typeToString(tp),
		}
		if property.Required() {
			required = append(required, property.Name())
		}
	}

	var data = map[string]any{
		"properties": properties,
	}

	if len(required) > 0 {
		data["required"] = required
	}

	return data
}

func (c *CommandableSwaggerDocument) createResponsesData() map[string]any {
	return map[string]any{

		"200": map[string]any{
			"description": "Successful response",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{
						"type": "object",
					},
				},
			},
		},
	}
}

func (c *CommandableSwaggerDocument) WriteData(indent int, data map[string]any) {

	for key, value := range data {
		if val, ok := value.(string); ok {
			c.writeAsString(indent, key, val)
		} else {
			if arr, ok := value.([]any); ok {
				if len(arr) > 0 {
					c.WriteName(indent, key)
					for index := 0; index < len(arr); index++ {
						item := arr[index].(string)
						c.writeArrayItem(indent+1, item, false)
					}
				}
			} else {
				if m, ok := value.(map[string]any); ok {
					notEmpty := false
					for _, v := range m {
						if v != nil {
							notEmpty = true
							break
						}
					}
					if notEmpty {
						c.WriteName(indent, key)
						c.WriteData(indent+1, m)
					}
				} else {
					c.writeAsObject(indent, key, value)
				}
			}
		}
	}
}

func (c *CommandableSwaggerDocument) WriteName(indent int, name string) {
	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ":\n"
}

func (c *CommandableSwaggerDocument) writeArrayItem(indent int, name string, isObjectItem bool) {
	var spaces = c.GetSpaces(indent)
	c.content += spaces + "- "

	if isObjectItem {
		c.content += name + ":\n"
	} else {
		c.content += name + "\n"
	}
}

func (c *CommandableSwaggerDocument) writeAsObject(indent int, name string, value any) {
	if value == nil {
		return
	}

	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ": " + cconv.StringConverter.ToString(value) + "\n"
}

func (c *CommandableSwaggerDocument) writeAsString(indent int, name string, value string) {
	if value == "" {
		return
	}

	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ": '" + value + "'\n"
}

func (c *CommandableSwaggerDocument) GetSpaces(length int) string {
	return strings.Repeat(" ", length*2)
}

func (c *CommandableSwaggerDocument) typeToString(tp cconv.TypeCode) string {
	// allowed types: array, boolean, integer, number, object, string
	if tp == cconv.Integer || tp == cconv.Long {
		return "integer"
	}
	if tp == cconv.Double || tp == cconv.Float {
		return "number"
	}
	if tp == cconv.String {
		return "string"
	}
	if tp == cconv.Boolean {
		return "boolean"
	}
	if tp == cconv.Array {
		return "array"
	}

	return "object"
}
