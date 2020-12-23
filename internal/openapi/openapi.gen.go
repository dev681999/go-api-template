// Package openapi provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Error defines model for Error.
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Status defines model for Status.
type Status struct {
	Message string `json:"message"`
}

// UserLoginRequest defines model for UserLoginRequest.
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token"`
}

// UserProfile defines model for UserProfile.
type UserProfile struct {
	Address   string `json:"address"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	ImageUrl  string `json:"image_url"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
}

// UserProfileUpdateRequest defines model for UserProfileUpdateRequest.
type UserProfileUpdateRequest struct {
	Address   *string `json:"address,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	ImageUrl  *string `json:"image_url,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Mobile    *string `json:"mobile,omitempty"`
}

// UserRegistrationRequest defines model for UserRegistrationRequest.
type UserRegistrationRequest struct {
	Address   string `json:"address"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
}

// ActivateUserParams defines parameters for ActivateUser.
type ActivateUserParams struct {

	// The activation token
	Token *string `json:"token,omitempty"`
}

// LoginUserJSONBody defines parameters for LoginUser.
type LoginUserJSONBody UserLoginRequest

// UpdateUserProfileJSONBody defines parameters for UpdateUserProfile.
type UpdateUserProfileJSONBody UserProfileUpdateRequest

// RegisterUserJSONBody defines parameters for RegisterUser.
type RegisterUserJSONBody UserRegistrationRequest

// LoginUserRequestBody defines body for LoginUser for application/json ContentType.
type LoginUserJSONRequestBody LoginUserJSONBody

// UpdateUserProfileRequestBody defines body for UpdateUserProfile for application/json ContentType.
type UpdateUserProfileJSONRequestBody UpdateUserProfileJSONBody

// RegisterUserRequestBody defines body for RegisterUser for application/json ContentType.
type RegisterUserJSONRequestBody RegisterUserJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /user/activate)
	ActivateUser(ctx echo.Context, params ActivateUserParams) error

	// (POST /user/login)
	LoginUser(ctx echo.Context) error

	// (GET /user/profile)
	UserProfile(ctx echo.Context) error

	// (PATCH /user/profile)
	UpdateUserProfile(ctx echo.Context) error

	// (POST /user/register)
	RegisterUser(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ActivateUser converts echo context to params.
func (w *ServerInterfaceWrapper) ActivateUser(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ActivateUserParams
	// ------------- Optional query parameter "token" -------------

	err = runtime.BindQueryParameter("form", true, false, "token", ctx.QueryParams(), &params.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter token: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ActivateUser(ctx, params)
	return err
}

// LoginUser converts echo context to params.
func (w *ServerInterfaceWrapper) LoginUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.LoginUser(ctx)
	return err
}

// UserProfile converts echo context to params.
func (w *ServerInterfaceWrapper) UserProfile(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserProfile(ctx)
	return err
}

// UpdateUserProfile converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUserProfile(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateUserProfile(ctx)
	return err
}

// RegisterUser converts echo context to params.
func (w *ServerInterfaceWrapper) RegisterUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RegisterUser(ctx)
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

	router.POST(baseURL+"/user/activate", wrapper.ActivateUser)
	router.POST(baseURL+"/user/login", wrapper.LoginUser)
	router.GET(baseURL+"/user/profile", wrapper.UserProfile)
	router.PATCH(baseURL+"/user/profile", wrapper.UpdateUserProfile)
	router.POST(baseURL+"/user/register", wrapper.RegisterUser)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXS2/bOBD+KwJ3j4KlJHtY6JYFtkXaAg3yQA+GEdDS2GYqkcxw5NYI9N8LkpKtxFTc",
	"BnEeRU9RNNTMfB+/efiW5arSSoIkw7JbZvIFVNw9/o+o0D5oVBqQBLjX0L2mlQaWMUMo5Jw1MavAGD6H",
	"gK2JGcJNLRAKlo3XB+PW2STuPlDTa8jJOjsnTrXZDv/LQULOLw3gJzUX8gxuajAUQFlxUQZRam7MN4XF",
	"7gy8j94XO1IxWkkD27nMBBq6kryCYEIlf8hK6ivI3bn6Y3E/Vt/zUOanqGaiDOTMiwLBmGBKw9zuQCoq",
	"PoerGstH8FCpaZvoz13aABFrR/EaYj+vHTxd6oITDIruTZEWxHkGc2EIOQklnxfmI4E8pp4HpbF2FVDJ",
	"tjKamBnIaxS0Ordd15MzBY6AxzUtNv+9U1hxYhn78OWCxb5HW0/eytaeF0SaNdaxkDPlAAmysNl7FR2f",
	"nkQXUOmSk81rCWiEkixjB6N0lFomlAbJtWAZOxqlo0MHiBYuq6Q2gAnPSSzt5/ZKlb/aAkyOQpP3ddye",
	"iKwYmHPpxXBS9KytUXPkFRCgYdn4vqeLBURtPKFk1PUnYW03NeCKxcxf97p3+dkVuseJvUjfXh2cwzS1",
	"f3IlCaSDwbUuRe5iJddGyc0stE9/I8xYxv5KNsMyaSdl0g4qR/tdCJ8/Wlb/ecJYfiQHQk15EWFbby7m",
	"0bPHLGDG65L2H7eW8F1DTlBE0J7ZVBLLxrbS+NxqyrUkNrF2L+DSDtph9bo5HJauM7WWFvR/qlg9Gdqt",
	"jSQA3KfXHej3JsIamj1qfHtJGZT7W5CB3mwtcwjowJ6PutXmvhL6a8+eKe/CvH6y786t8aTZZt9Nk3zx",
	"MNuR35C2SXev71O/nzIM7moBTgJ5v1hx/lZKWdcpum0ScLhjn7Unwk27s+65b4eW3gAz/WMvIpQ/m8rr",
	"31R2F4sBXHYrs/tV59b+LElKlfNyoQxl/6ZpmnAtkuUBaybNjwAAAP//FFHGAuQRAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
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

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
