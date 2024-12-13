package apitool

import (
	"github.com/94peter/api-toolkit/mid"
	"github.com/gin-gonic/gin"
)

const envCtxUserKey = "API_CTX_USER_KEY"

type ConfigWithBindUser[T mid.BindUser] struct {
	bindUser   T
	CtxUserKey string
	*Config
}

func GetConfigWithBindUserFromEnv[T mid.BindUser](bind T) (*ConfigWithBindUser[T], error) {
	cfg, err := GetConfigFromEnv()
	if err != nil {
		return nil, err
	}
	ctxUserKey, err := stringFromEnv(envCtxUserKey)
	if err != nil {
		return nil, err
	}
	return &ConfigWithBindUser[T]{
		bindUser:   bind,
		CtxUserKey: ctxUserKey,
		Config:     cfg,
	}, nil
}

func (cfg *ConfigWithBindUser[T]) SetAPIs(apis ...GinAPIWithBindUser[T]) {
	cfg.Config.apis = make([]GinAPI, len(apis))
	for i, api := range apis {
		api.SetReqUserHandler(cfg.CtxUserKey)
		cfg.Config.apis[i] = api
	}
}

type GinReqUserHandler[T mid.BindUser] struct {
	ctxUserKey string
}

func (h *GinReqUserHandler[T]) SetReqUserHandler(ctxUserKey string) {
	if h == nil {
		panic("missing handler")
	}
	if h.ctxUserKey != "" {
		panic("handler already set")
	}
	h.ctxUserKey = ctxUserKey
}

func (h *GinReqUserHandler[T]) GetReqUser(c *gin.Context) *T {
	if h == nil {
		return nil
	}
	user, ok := c.Get(h.ctxUserKey)
	if !ok {
		return nil
	}
	return user.(*T)
}

type GinAPIWithBindUser[T mid.BindUser] interface {
	GinAPI
	SetReqUserHandler(ctxUserKey string)
	GetReqUser(c *gin.Context) *T
}
