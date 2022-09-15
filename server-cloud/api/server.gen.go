// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Organization defines model for Organization.
type Organization struct {
	CreatedOn   *openapi_types.Date `json:"createdOn,omitempty"`
	Description *string             `json:"description,omitempty"`
	Name        string              `json:"name"`
	UpdatedOn   *openapi_types.Date `json:"updatedOn,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// CreateOrganizationParams defines parameters for CreateOrganization.
type CreateOrganizationParams struct {
	Name        string  `form:"name" json:"name"`
	Description *string `form:"description,omitempty" json:"description,omitempty"`
}

// UpdateOrganizationParams defines parameters for UpdateOrganization.
type UpdateOrganizationParams struct {
	Description *string `form:"description,omitempty" json:"description,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /v1/organizations)
	CreateOrganization(ctx echo.Context, params CreateOrganizationParams) error

	// (DELETE /v1/organizations/{name})
	DeleteOrganization(ctx echo.Context, name string) error

	// (GET /v1/organizations/{name})
	FindOrganization(ctx echo.Context, name string) error

	// (PUT /v1/organizations/{name})
	UpdateOrganization(ctx echo.Context, name string, params UpdateOrganizationParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateOrganization converts echo context to params.
func (w *ServerInterfaceWrapper) CreateOrganization(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateOrganizationParams
	// ------------- Required query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, true, "name", ctx.QueryParams(), &params.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// ------------- Optional query parameter "description" -------------

	err = runtime.BindQueryParameter("form", true, false, "description", ctx.QueryParams(), &params.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter description: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateOrganization(ctx, params)
	return err
}

// DeleteOrganization converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteOrganization(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteOrganization(ctx, name)
	return err
}

// FindOrganization converts echo context to params.
func (w *ServerInterfaceWrapper) FindOrganization(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindOrganization(ctx, name)
	return err
}

// UpdateOrganization converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateOrganization(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params UpdateOrganizationParams
	// ------------- Optional query parameter "description" -------------

	err = runtime.BindQueryParameter("form", true, false, "description", ctx.QueryParams(), &params.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter description: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateOrganization(ctx, name, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/v1/organizations", wrapper.CreateOrganization)
	router.DELETE(baseURL+"/v1/organizations/:name", wrapper.DeleteOrganization)
	router.GET(baseURL+"/v1/organizations/:name", wrapper.FindOrganization)
	router.PUT(baseURL+"/v1/organizations/:name", wrapper.UpdateOrganization)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RUwW7UMBD9FTRwtJpte8sNFSrBpVIBcah6cOPZrKvEdsfjSssq/47G3m02JKBKBYHE",
	"Xtaesd/Me8+THTS+D96h4wj1Dghj8C5i3rwn8nS9j0ig8Y7RsSx1CJ1tNFvvqvvoncRis8Fey+oN4Rpq",
	"eF2N6FXJxiqjwiA/tb8yVpNFIB+Q2JYmGm9y7bWnXjPUYB2fn4EC3gYsW2yRYFDQY4y6zaf3ychkXQtS",
	"ifAhWUID9U3BHM/fPoH5u3tsWLCuqNXOfssEF5oi1Izmyk06M5pxbOxQW4HB2JANB6hZ3ukeFxMpmGeX",
	"+YFixpwTE8mxSWR5+0mkL3TuUBPS28SbcXd5KPfx62fYGyVIJTvW3zAHyF5at/aZhuVOMte2615ddD4Z",
	"UPCIFDN/OD1ZnayEnQ/odLBQw3kOKQiaN7mh6vG08kcO5GDwMT89MSJHPxio4SJbMbFLgEj3yEgR6psd",
	"WKn7kJC2cFC7/B1rxpRQHb3hmb7LOMfm/ur6rZoO19lq9dtGasK+DJbBtU4d/+zmUyvVdMjL5Zn81U7Y",
	"DoJmsEPGuQ/vcvwZPojJL7BhWccXMVbQ4sLLurTO/C0+/+C7UBDSgkpf8jfqz+v034zfMHwPAAD//+78",
	"3T6VBwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}