package apitool

import (
	"strconv"

	"github.com/94peter/api-toolkit/auth"
	"github.com/94peter/api-toolkit/errors"
	"github.com/94peter/api-toolkit/mid"
	"github.com/gin-gonic/gin"
)

type GinApiHandler struct {
	Handler func(c *gin.Context)
	Method  string
	Path    string
	Auth    bool
	Group   []auth.ApiPerm
}

type GinAPI interface {
	errors.ApiErrorHandler
	GetAPIs() []*GinApiHandler
}

type GinApiServer interface {
	AddAPIs(handlers ...GinAPI) GinApiServer
	Middles(mids ...mid.GinMiddle) GinApiServer
	SetServerErrorHandler(errors.GinServerErrorHandler) GinApiServer
	SetAuth(authmid auth.GinAuthMidInter) GinApiServer
	SetTrustedProxies([]string) GinApiServer
	Static(relativePath, root string) GinApiServer
	Run(port int) error
	errorHandler(c *gin.Context, err error)
}

type ginApiServ struct {
	*gin.Engine
	service      string
	authMid      auth.GinAuthMidInter
	myErrHandler errors.GinServerErrorHandler
}

func (serv *ginApiServ) SetServerErrorHandler(handler errors.GinServerErrorHandler) GinApiServer {
	serv.myErrHandler = handler
	return serv
}

func (serv *ginApiServ) errorHandler(c *gin.Context, err error) {
	serv.myErrHandler(c, serv.service, err)
}

func (serv *ginApiServ) Static(relativePath, root string) GinApiServer {
	serv.Engine.Static(relativePath, root)
	return serv
}

func (serv *ginApiServ) SetAuth(authMid auth.GinAuthMidInter) GinApiServer {
	serv.authMid = authMid
	return serv
}

func (serv *ginApiServ) Middles(mids ...mid.GinMiddle) GinApiServer {
	for _, m := range mids {
		m.SetApiErrorHandler(serv.errorHandler)
		serv.Engine.Use(m.Handler())
	}
	return serv
}

func (serv *ginApiServ) AddAPIs(apis ...GinAPI) GinApiServer {
	for _, api := range apis {
		api.SetApiErrorHandler(serv.errorHandler)
		for _, h := range api.GetAPIs() {
			if serv.authMid != nil {
				serv.authMid.AddAuthPath(h.Path, h.Method, h.Auth, h.Group)
			}
			switch h.Method {
			case "GET":
				serv.Engine.GET(h.Path, h.Handler)
			case "POST":
				serv.Engine.POST(h.Path, h.Handler)
			case "PUT":
				serv.Engine.PUT(h.Path, h.Handler)
			case "DELETE":
				serv.Engine.DELETE(h.Path, h.Handler)
			}
		}
	}
	return serv
}

func (serv *ginApiServ) SetTrustedProxies(proxies []string) GinApiServer {
	if len(proxies) == 0 {
		return serv
	}
	serv.Engine.ForwardedByClientIP = true
	serv.Engine.SetTrustedProxies(proxies)
	return serv
}

func (serv *ginApiServ) Run(port int) error {
	return serv.Engine.Run(":" + strconv.Itoa(port))
}

func NewGinApiServer(mode string, service string) GinApiServer {
	gin.SetMode(mode)
	return &ginApiServ{
		Engine:  gin.New(),
		service: service,
	}
}