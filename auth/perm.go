package auth

type ApiPerm string

type ReqUser interface {
	Host() string
	GetPerm() []string
}
