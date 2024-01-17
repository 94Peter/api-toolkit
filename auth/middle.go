package auth

import (
	"github.com/94peter/api-toolkit/mid"
)

type GinAuthMidInter interface {
	mid.GinMiddle
	AddAuthPath(path string, method string, isAuth bool, group []ApiPerm)
	IsAuth(path string, method string) bool
	HasPerm(path, method string, perm []string) bool
}
