package mid

import (
	"github.com/94peter/api-toolkit/errors"
	"github.com/gin-gonic/gin"
)

type GinMiddle interface {
	errors.ApiErrorHandler
	Handler() gin.HandlerFunc
}
