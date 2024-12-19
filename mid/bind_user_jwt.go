package mid

import (
	"reflect"

	"github.com/94peter/api-toolkit/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type BindUserJwtMidOption[T BindUser] func(*bindUserJwtMiddle[T])

func BindUserJwtMidWithSecret[T BindUser](secret string) BindUserJwtMidOption[T] {
	return func(m *bindUserJwtMiddle[T]) {
		m.secretKey = secret
	}
}

func BindUserJwtMidWithCtxKey[T BindUser](ctxKey string) BindUserJwtMidOption[T] {
	return func(m *bindUserJwtMiddle[T]) {
		m.ctxKey = ctxKey
	}
}

func BindUserJwtMidWithBindObject[T BindUser](bindObj T) BindUserJwtMidOption[T] {
	return func(m *bindUserJwtMiddle[T]) {
		if reflect.TypeOf(bindObj).Kind() == reflect.Ptr &&
			reflect.ValueOf(bindObj).IsNil() {
			panic("bind object is nil")
		}
		m.bindType = bindObj
	}
}

func BindUserJwtMidWithMock[T BindUser]() BindUserJwtMidOption[T] {
	return func(m *bindUserJwtMiddle[T]) {
		m.isMock = true
	}
}

func NewGinBindUserJwtMid[T BindUser](opts ...BindUserJwtMidOption[T]) GinMiddle {
	middle := &bindUserJwtMiddle[T]{}
	for _, opt := range opts {
		opt(middle)
	}
	return middle
}

type bindUserJwtMiddle[T BindUser] struct {
	isMock    bool
	secretKey string
	ctxKey    string
	bindType  T
	errors.CommonApiErrorHandler
}

func (m *bindUserJwtMiddle[T]) newObj() BindUser {
	if reflect.TypeOf(m.bindType).Kind() == reflect.Ptr {
		return reflect.New(reflect.TypeOf(m.bindType).Elem()).Interface().(BindUser)
	} else {
		return reflect.New(reflect.TypeOf(m.bindType)).Interface().(BindUser)
	}
}

func (m *bindUserJwtMiddle[T]) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := getTokenString(c)
		if tokenString == "" {
			c.Next()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		})

		if err != nil {
			c.Next()
			return
		}

		newObj := m.newObj()
		bindClaims(newObj, token.Claims.(jwt.MapClaims))
		if !newObj.IsEmpty() {
			var obj T
			if reflect.TypeOf(obj).Kind() == reflect.Struct {
				newObj = reflect.ValueOf(newObj).Elem().Interface().(BindUser)
			}
			c.Set(m.ctxKey, newObj)
		}
		c.Next()
	}
}

func getTokenString(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return ""
	}
	return tokenString[7:]
}

func bindClaims(obj any, claims map[string]any) {
	rv := reflect.ValueOf(obj)
	rv = rv.Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		claimsTag := field.Tag.Get("claims")
		if claimsTag != "" {
			if value, ok := claims[claimsTag]; ok {
				rv.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
}
