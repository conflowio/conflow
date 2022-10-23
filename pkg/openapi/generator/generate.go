// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"fmt"
	"regexp"
	"strings"
	texttemplate "text/template"

	"github.com/go-openapi/inflect"

	conflowgenerator "github.com/conflowio/conflow/pkg/conflow/generator"
	"github.com/conflowio/conflow/pkg/conflow/generator/template"
	"github.com/conflowio/conflow/pkg/openapi"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

type generator struct {
	o           *openapi.OpenAPI
	router      string
	packageName string
	imports     map[string]string
	schemas     []schema.Schema
}

type Operation struct {
	Path        string
	Method      string
	OperationID string
	Parameters  map[string]*openapi.Parameter
}

func Generate(o *openapi.OpenAPI, router string, packageName, outputDir string) error {
	g := generator{
		o:           o,
		router:      router,
		packageName: packageName,
		imports: map[string]string{
			packageName: "",
		},
	}
	if err := g.generateTypes(outputDir); err != nil {
		return err
	}

	g.imports = map[string]string{
		packageName: "",
	}
	if err := g.generateSchemas(outputDir); err != nil {
		return err
	}

	g.imports = map[string]string{
		packageName: "",
	}
	if err := g.generateServer(outputDir); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateTypes(outputDir string) error {
	body := &strings.Builder{}

	for _, name := range util.SortedMapKeys(g.o.Schemas) {
		g.writeType(body, name, g.o.Schemas[name])
	}

	for _, name := range util.SortedMapKeys(g.o.RequestBodies) {
		g.writeRequestBody(body, name, g.o.RequestBodies[name])
	}

	var operations []Operation
	for _, path := range util.SortedMapKeys(g.o.Paths) {
		p := g.o.Paths[path]
		_ = p.IterateOperations(func(method string, op *openapi.Operation) error {
			operations = append(operations, Operation{
				Path:        path,
				Method:      method,
				OperationID: op.OperationID,
			})
			return nil
		})

		if err := g.writePath(body, p); err != nil {
			return fmt.Errorf("failed to generate path %q: %w", path, err)
		}
	}

	b, err := util.GenerateTemplate(
		serverTemplate,
		ServerTemplateParams{
			Operations: operations,
			Imports:    g.imports,
		},
		nil,
	)
	if err != nil {
		return err
	}
	body.Write(b)
	body.WriteRune('\n')

	packageParts := strings.Split(g.packageName, "/")

	header, err := template.GenerateHeader(template.HeaderParams{
		Package: packageParts[len(packageParts)-1],
		Imports: g.imports,
	})
	if err != nil {
		return err
	}

	return conflowgenerator.WriteGeneratedFile(outputDir, "openapi_types", append(header, []byte(body.String())...))
}

func (g *generator) generateSchemas(outputDir string) error {
	for _, s := range g.schemas {
		g.replaceReferences(s)

		// TODO: examples might contain types like time.Time
		// We have to encode that field without any types
		_ = schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
			meta.SetExamples(nil)
			return nil
		})
	}

	body, err := util.GenerateTemplate(
		schemaTemplate,
		schemasTemplateParams{
			Schemas: g.schemas,
			Imports: g.imports,
		},
		nil,
	)
	if err != nil {
		return err
	}

	packageParts := strings.Split(g.packageName, "/")

	header, err := template.GenerateHeader(template.HeaderParams{
		Package: packageParts[len(packageParts)-1],
		Imports: g.imports,
	})
	if err != nil {
		return err
	}

	return conflowgenerator.WriteGeneratedFile(outputDir, "openapi_schemas", append(header, body...))
}

func (g *generator) generateServer(outputDir string) error {
	body := &strings.Builder{}

	var operations []Operation
	for _, path := range util.SortedMapKeys(g.o.Paths) {
		p := g.o.Paths[path]
		_ = p.IterateOperations(func(method string, op *openapi.Operation) error {
			operations = append(operations, Operation{
				Path:        path,
				Method:      method,
				OperationID: op.OperationID,
				Parameters:  g.mergeParams(p.Parameters, op.Parameters),
			})
			return nil
		})
	}

	switch g.router {
	case "echo":
		paramRegex := regexp.MustCompile(`{([^{}]+)}`)
		b, err := util.GenerateTemplate(
			echoServerTemplate,
			EchoServerTemplateParams{
				Operations: operations,
				Imports:    g.imports,
			},
			texttemplate.FuncMap{
				"convertPath": func(path string) string {
					return paramRegex.ReplaceAllString(path, ":$1")
				},
				"bindParameterFunc": func(p *openapi.Parameter, imports map[string]string) string {
					if a, ok := p.Schema.(*schema.Array); ok {
						return fmt.Sprintf("BindParameterArray[%s]", a.Items.GoType(imports))
					} else {
						if n, ok := p.Schema.(schema.Nullable); ok && n.GetNullable() {
							return fmt.Sprintf("BindParameterPtr[%s]", strings.TrimPrefix(p.Schema.GoType(imports), "*"))
						} else {
							return fmt.Sprintf("BindParameter[%s]", p.Schema.GoType(imports))
						}
					}
				},
			},
		)
		if err != nil {
			return err
		}

		body.Write(b)
		body.WriteRune('\n')
	default:
		return fmt.Errorf("unsupported router: %s", g.router)
	}

	packageParts := strings.Split(g.packageName, "/")

	header, err := template.GenerateHeader(template.HeaderParams{
		Package: packageParts[len(packageParts)-1],
		Imports: g.imports,
	})
	if err != nil {
		return err
	}

	return conflowgenerator.WriteGeneratedFile(outputDir, "openapi_server", append(header, []byte(body.String())...))
}

func (g *generator) writePath(b *strings.Builder, p *openapi.PathItem) error {
	if err := p.IterateOperations(func(_ string, op *openapi.Operation) error {
		if err := g.writeOperation(b, op, p.Parameters); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (g *generator) writeRequestBody(b *strings.Builder, name string, requestBody *openapi.RequestBody) {
	for _, contentType := range util.SortedMapKeys(requestBody.Content) {
		g.writeContentType(b, contentType, requestBody.Content[contentType].Schema, name, "RequestBody")
	}
}

func (g *generator) writeOperation(b *strings.Builder, operation *openapi.Operation, pathParams []*openapi.Parameter) error {
	if operation == nil {
		return nil
	}

	if operation.RequestBody != nil {
		g.writeRequestBody(b, operation.OperationID, operation.RequestBody)
	}

	g.writeRequest(b, operation.OperationID, pathParams, operation.Parameters, operation.RequestBody)

	t, err := util.GenerateTemplate(
		responseTemplate,
		ResponseTemplateParams{
			OperationID: operation.OperationID,
			Imports:     g.imports,
		},
		nil,
	)
	if err != nil {
		return err
	}

	b.Write(t)
	b.WriteRune('\n')

	for _, responseCode := range util.SortedMapKeys(operation.Responses) {
		if err := g.writeResponse(b, operation.OperationID, responseCode, operation.Responses[responseCode]); err != nil {
			return fmt.Errorf("failed to generate response %s for operation %s: %w", responseCode, operation.OperationID, err)
		}
	}

	return nil
}

func (g *generator) writeRequest(
	b *strings.Builder,
	name string,
	pathParams []*openapi.Parameter,
	operationParams []*openapi.Parameter,
	requestBody *openapi.RequestBody,
) {
	params := g.mergeParams(pathParams, operationParams)

	request := &schema.Object{
		Properties: map[string]schema.Schema{},
	}

	if requestBody != nil {
		for _, contentType := range util.SortedMapKeys(requestBody.Content) {
			request.Properties[fmt.Sprintf("Body%s", g.contentTypeName(contentType))] = requestBody.Content[contentType].Schema
		}
	}

	for _, name := range util.SortedMapKeys(params) {
		p := params[name]
		if p.Schema == nil {
			p.Schema = schema.StringValue()
		}
		request.Properties[name] = p.Schema
		if ptr.Value(p.Required) {
			request.Required = append(request.Required, name)
		}
	}

	g.writeType(b, fmt.Sprintf("%sRequest", inflect.Camelize(name)), request)
}

func (g *generator) writeResponse(b *strings.Builder, operationID string, responseCode string, response *openapi.Response) error {
	responseName := fmt.Sprintf("%sResponse%s", inflect.Camelize(operationID), inflect.Camelize(responseCode))

	responseObj := &schema.Object{
		Properties: map[string]schema.Schema{},
	}
	if responseCode == "default" {
		responseObj.Properties["ResponseCode"] = schema.IntegerValue()
		responseObj.Required = append(responseObj.Required, "ResponseCode")
	}

	for _, contentType := range util.SortedMapKeys(response.Content) {
		content := response.Content[contentType]
		g.writeContentType(b, contentType, content.Schema, fmt.Sprintf("%sResponse%s", inflect.Camelize(operationID), inflect.Camelize(responseCode)), "Body")
		content.Schema.(schema.Nullable).SetNullable(true)

		responseObj.Properties[fmt.Sprintf("Body%s", g.contentTypeName(contentType))] = content.Schema
	}

	g.writeType(b, responseName, responseObj)

	t, err := util.GenerateTemplate(
		writeResponseFuncTemplate,
		WriteResponseFuncTemplateParams{
			OperationID:  operationID,
			ResponseCode: responseCode,
			ContentTypes: util.SortedMapKeys(response.Content),
			Imports:      g.imports,
		},
		texttemplate.FuncMap{
			"contentTypeName": g.contentTypeName,
		},
	)
	if err != nil {
		return err
	}

	b.Write(t)
	b.WriteRune('\n')

	return nil
}

func (g *generator) writeContentType(b *strings.Builder, contentType string, s schema.Schema, prefix, suffix string) {
	name := fmt.Sprintf("%s%s%s", inflect.Camelize(prefix), g.contentTypeName(contentType), suffix)
	g.writeType(b, name, s)
}

func (g *generator) contentTypeName(contentType string) string {
	return inflect.Camelize(strings.ReplaceAll(contentType, "/", "_"))
}

func (g *generator) writeType(b *strings.Builder, name string, s schema.Schema) {
	if s.GetID() != "" {
		return
	}

	if so, ok := s.(*schema.Object); ok {
		for jsonPropertyName := range so.Properties {
			fieldName, ok := so.FieldNames[jsonPropertyName]
			if !ok {
				fieldName = inflect.Camelize(jsonPropertyName)
			}
			fieldName = inflect.Capitalize(fieldName)
			if fieldName != jsonPropertyName {
				if so.FieldNames == nil {
					so.FieldNames = map[string]string{}
				}
				so.FieldNames[jsonPropertyName] = fieldName
			}
		}
	}

	_, _ = b.Write([]byte(fmt.Sprintf("type %s %s\n\n", name, s.GoType(g.imports))))

	_ = schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetID(fmt.Sprintf("%s.%s", g.packageName, name))
		return nil
	})

	g.schemas = append(g.schemas, s)
}

func (g *generator) replaceReferences(s schema.Schema) {
	switch st := s.(type) {
	case *schema.Array:
		g.replaceReferences(st.Items)
	case *schema.Map:
		g.replaceReferences(st.AdditionalProperties)
	case *schema.Object:
		for _, p := range st.Properties {
			g.replaceReferences(p)
		}
	case *schema.Reference:
		if strings.HasPrefix(st.Ref, "#/components/schemas/") {
			name := strings.TrimPrefix(st.Ref, "#/components/schemas/")
			if s, ok := g.o.Schemas[name]; ok {
				st.Ref = s.GetID()
			}
		}
	}
}

func (g *generator) mergeParams(pathParams, opParams []*openapi.Parameter) map[string]*openapi.Parameter {
	params := map[string]*openapi.Parameter{}
	for _, p := range pathParams {
		name := fmt.Sprintf("%s%s", inflect.Camelize(p.In), inflect.Camelize(p.Name))
		params[name] = p
	}
	for _, p := range opParams {
		name := fmt.Sprintf("%s%s", inflect.Camelize(p.In), inflect.Camelize(p.Name))
		params[name] = p
	}
	return params
}
