package mid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Name string `header:"X-User-Name"`
}

func (t testObj) IsEmpty() bool { return t.Name == "" }

func (t testObj) BindClaims(claims jwt.Claims) {}

func TestNewGinBindUserMidWithPointer(t *testing.T) {

	tests := []struct {
		name       string
		bindObj    *testObj
		hasPanic   bool
		reqHeader  map[string][]string
		statusCode int
	}{
		{
			name:     "pointer bind object",
			bindObj:  &testObj{Name: "testaaaa"},
			hasPanic: false,
			reqHeader: map[string][]string{
				"X-User-Name": {"test"},
			},
			statusCode: http.StatusOK,
		},
		{
			name:     "nil bind object",
			bindObj:  nil,
			hasPanic: true,
		},
	}
	const userCtxKey = "User"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasPanic {
				assert.Panics(t, func() {
					NewGinBindUserMid(BindUserMidWithBindObject(tt.bindObj))
				})
				return
			}
			mid := NewGinBindUserMid(BindUserMidWithBindObject(tt.bindObj))
			assert.NotNil(t, mid)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)

			c.Request.Header = tt.reqHeader

			m := NewGinBindUserMid(
				BindUserMidWithBindObject(tt.bindObj),
				BindUserMidWithCtxKey[*testObj](userCtxKey))
			handler := m.Handler()
			handler(c)

			if w.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
			}

			if _, ok := tt.reqHeader["X-User-Name"]; ok {
				userInter, ok := c.Get(userCtxKey)
				user := userInter.(*testObj)
				assert.True(t, ok)
				assert.Equal(t, tt.reqHeader["X-User-Name"][0], user.Name)
				return
			} else {
				userInter, ok := c.Get(userCtxKey)
				assert.False(t, ok)
				assert.Nil(t, userInter)
			}
		})
	}
}

func TestBindUserMiddleHandlerWithNonePointer(t *testing.T) {

	tests := []struct {
		name       string
		bindObj    testObj
		reqHeader  map[string][]string
		statusCode int
	}{
		{
			name:    "valid request",
			bindObj: testObj{},
			reqHeader: map[string][]string{
				"X-User-Name": {"test"},
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "valid request",
			bindObj:    testObj{},
			reqHeader:  map[string][]string{},
			statusCode: http.StatusOK,
		},
	}
	const userCtxKey = "User"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)

			c.Request.Header = tt.reqHeader

			m := NewGinBindUserMid(
				BindUserMidWithBindObject(tt.bindObj),
				BindUserMidWithCtxKey[testObj](userCtxKey))
			handler := m.Handler()
			handler(c)

			if w.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
			}

			if _, ok := tt.reqHeader["X-User-Name"]; ok {
				userInter, ok := c.Get(userCtxKey)
				user := userInter.(testObj)
				assert.True(t, ok)
				assert.Equal(t, tt.reqHeader["X-User-Name"][0], user.Name)
				return
			} else {
				userInter, ok := c.Get(userCtxKey)
				assert.False(t, ok)
				assert.Nil(t, userInter)
			}
		})
	}
}

func Test_ReflectMindMap(t *testing.T) {
	type A struct {
		Name string `claims:"id"`
	}
	var emptyA A
	claims := map[string]any{
		"id": "Jack",
	}
	bindClaims(&emptyA, claims)
	fmt.Println(emptyA.Name) // print Jack
	assert.Equal(t, "Jack", emptyA.Name)
}
