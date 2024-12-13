package mid

import (
	"reflect"

	"github.com/94peter/api-toolkit/errors"
	"github.com/gin-gonic/gin"
)

type BindUser interface {
	IsEmpty() bool
}

type BindUserMidOption[T BindUser] func(*bindUserMiddle[T])

func BindUserMidWithCtxKey[T BindUser](ctxKey string) BindUserMidOption[T] {
	return func(m *bindUserMiddle[T]) {
		m.ctxKey = ctxKey
	}
}

func BindUserMidWithBindObject[T BindUser](bindObj T) BindUserMidOption[T] {
	return func(m *bindUserMiddle[T]) {
		if reflect.TypeOf(bindObj).Kind() == reflect.Ptr &&
			reflect.ValueOf(bindObj).IsNil() {
			panic("bind object is nil")
		}
		m.bindType = bindObj
	}
}

func NewGinBindUserMid[T BindUser](opts ...BindUserMidOption[T]) GinMiddle {
	middle := &bindUserMiddle[T]{}
	for _, opt := range opts {
		opt(middle)
	}
	return middle
}

type bindUserMiddle[T BindUser] struct {
	ctxKey   string
	bindType T
	errors.CommonApiErrorHandler
}

func (m *bindUserMiddle[T]) newObj() BindUser {
	return reflect.New(reflect.TypeOf(m.bindType)).Interface().(BindUser)
}

func (m *bindUserMiddle[T]) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		newObj := m.newObj()
		if err := c.BindHeader(newObj); err != nil {
			c.Abort()
			return
		}
		if !newObj.IsEmpty() {
			c.Set(m.ctxKey, newObj)
			return
		}

		c.Next()
	}
}
