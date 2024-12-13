package mid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Name string `header:"X-User-Name"`
}

func (t testObj) IsEmpty() bool { return t.Name == "" }

func TestNewGinBindUserMidWithPointer(t *testing.T) {

	tests := []struct {
		name     string
		bindObj  *testObj
		hasPanic bool
	}{
		{
			name:     "pointer bind object",
			bindObj:  &testObj{Name: "testaaaa"},
			hasPanic: false,
		},
		{
			name:     "nil bind object",
			bindObj:  nil,
			hasPanic: true,
		},
	}
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
