package mid

import (
	"github.com/94peter/api-toolkit/errors"
	"github.com/gin-gonic/gin"
)

type GinMiddle interface {
	errors.ApiErrorHandler
	Handler() gin.HandlerFunc
}

func NewGinMiddle(handler gin.HandlerFunc) GinMiddle {
	return &baseMiddle{
		handler: handler,
	}
}

type baseMiddle struct {
	handler gin.HandlerFunc
	errors.CommonApiErrorHandler
}

func (m *baseMiddle) Handler() gin.HandlerFunc {
	return m.handler
}
