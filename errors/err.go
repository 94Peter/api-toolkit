package errors

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type GinApiErrorHandler func(c *gin.Context, err error)
type ApiErrorHandler interface {
	SetApiErrorHandler(GinApiErrorHandler)
}

type GinServerErrorHandler func(c *gin.Context, service string, err error)

type ApiError interface {
	GetStatus() int
	error
}

type myApiError struct {
	statusCode int
	error
}

func (e myApiError) GetStatus() int {
	return e.statusCode
}

func (e myApiError) String() string {
	return fmt.Sprintf("%v: %v", e.statusCode, e.error)
}

func New(status int, msg string) ApiError {
	return myApiError{statusCode: status, error: errors.New(msg)}
}

func PkgError(status int, err error) ApiError {
	return myApiError{statusCode: status, error: err}
}

type CommonApiErrorHandler struct {
	GinApiErrorHandler
}

func (api *CommonApiErrorHandler) SetApiErrorHandler(handler GinApiErrorHandler) {
	api.GinApiErrorHandler = handler
}

var (
	Error_Auth_Path_NotFound  = errors.New("auth path not found")
	Error_Auth_Miss_Token     = errors.New("miss token")
	Error_Auth_Invalid_Token  = errors.New("invalid token")
	Error_Auth_Host_Not_Match = errors.New("host not match")
	Error_Auth_No_Perm        = errors.New("no permission")
)
