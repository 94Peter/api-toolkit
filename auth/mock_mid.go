package auth

import (
	"strings"

	"github.com/94peter/api-toolkit/errors"
	"github.com/gin-gonic/gin"
)

func NewMockAuthMid() GinAuthMidInter {
	return &mockAuthMiddle{}
}

type mockAuthMiddle struct {
	errors.CommonApiErrorHandler
}

func (am *mockAuthMiddle) AddAuthPath(path string, method string, isAuth bool, group []ApiPerm) {
}

const (
	_MOCK_HEADER_KEY_UID     = "Mock_User_UID"
	_MOCK_HEADER_KEY_ACCOUNT = "Mock_User_ACC"
	_MOCK_HEADER_KEY_NAME    = "Mock_User_NAM"
	_MOCK_HEADER_KEY_ROLES   = "Mock_User_Roles"
)

func (am *mockAuthMiddle) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader(_MOCK_HEADER_KEY_UID)
		if userID == "" {
			userID = "mock-id"
		}
		userAcc := c.GetHeader(_MOCK_HEADER_KEY_ACCOUNT)
		if userAcc == "" {
			userAcc = "mock-account"
		}
		userName := c.GetHeader(_MOCK_HEADER_KEY_NAME)
		if userName == "" {
			userName = "mock-name"
		}
		roles := strings.Split(c.GetHeader(_MOCK_HEADER_KEY_ROLES), ",")
		if len(roles) == 0 {
			roles = []string{"mock"}
		}
		c.Set(
			_KEY_USER_INFO,
			NewReqUser(getHost(c.Request), userID, userAcc, userName, roles),
		)
		c.Next()
	}
}

func NewReqUser(host string, uid string, account string, name string, roles []string) ReqUser {
	return &mockUser{
		host:    host,
		UID:     uid,
		Account: account,
		Name:    name,
		Roles:   roles,
	}
}

type mockUser struct {
	host    string
	UID     string
	Account string
	Name    string
	Roles   []string
}

func (u *mockUser) Host() string {
	return u.host
}

func (u *mockUser) GetHost() string {
	return u.host
}

func (u *mockUser) GetPerm() []string {
	return u.Roles
}
