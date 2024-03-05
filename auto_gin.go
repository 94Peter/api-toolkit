package apitool

import (
	"errors"
	"net/http"
)

func AutoGinApiServer(cfg *Config) (*http.Server, error) {
	if cfg.errorHandler == nil {
		return nil, errors.New("missing error handler")
	}

	if cfg.authMid == nil {
		return nil, errors.New("missing auth middleware")
	}

	server := NewGinApiServer(cfg.GinMode, cfg.Service).
		SetServerErrorHandler(cfg.errorHandler).
		SetAuth(cfg.authMid).
		Middles(cfg.getMiddles()...).
		AddAPIs(cfg.apis...).
		SetTrustedProxies(cfg.TrustedProxies)

	if len(cfg.proms) > 0 {
		server = server.SetPromhttp(cfg.proms...)
	}

	if cfg.Logger != nil {
		authMode := "release"
		if cfg.IsMockAuth {
			authMode = "mock"
		}
		cfg.Logger.Infof("run api at port: [%d], auth mode: [%s]",
			cfg.ApiPort, authMode)
	}
	return server.GetServer(cfg.ApiPort), nil
}
