package apitool

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/94peter/api-toolkit/mid"
)

func autoGinApiServer(cfg *Config) (*http.Server, error) {
	if cfg.errorHandler == nil {
		return nil, errors.New("missing error handler")
	}

	server := NewGinApiServer(cfg.GinMode, cfg.Service).
		SetServerErrorHandler(cfg.errorHandler)

	if cfg.store != nil {
		if cfg.SessionHeaderName == "" {
			return nil, errors.New("missing env SESSION_HEADER_NAME")
		}
		if cfg.SessionExpired < 0 {
			return nil, errors.New("missing env SESSION_EXPIRED or set SessionExpired must > 0s")
		}
		server = server.SetSession(cfg.SessionHeaderName, cfg.store, cfg.SessionExpired)
	}
	if cfg.authMid != nil {
		server = server.SetAuth(cfg.authMid)
	}
	server = server.Middles(cfg.getMiddles()...).
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

func AutoGinApiRun(ctx context.Context, cfg *Config) error {
	var apiWait sync.WaitGroup
	server, err := autoGinApiServer(cfg)
	if err != nil {
		return err
	}
	const fiveSecods = 5 * time.Second
	apiWait.Add(1)
	go func(srv *http.Server) {

		for {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				cfg.Logger.Fatalf("listen: %s", err)
				time.Sleep(fiveSecods)
			} else if err == http.ErrServerClosed {
				apiWait.Done()
				return
			}
		}
	}(server)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecods)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger.Fatalf("Server forced to shutdown: %v", err)
	}
	apiWait.Wait()
	return nil
}

func autoGinApiServerWithBindUser[T mid.BindUser](cfg *ConfigWithBindUser[T]) (*http.Server, error) {
	if cfg.errorHandler == nil {
		return nil, errors.New("missing error handler")
	}

	server := NewGinApiServer(cfg.GinMode, cfg.Service).
		SetServerErrorHandler(cfg.errorHandler)

	if cfg.store != nil {
		if cfg.SessionHeaderName == "" {
			return nil, errors.New("missing env SESSION_HEADER_NAME")
		}
		if cfg.SessionExpired < 0 {
			return nil, errors.New("missing env SESSION_EXPIRED or set SessionExpired must > 0s")
		}
		server = server.SetSession(cfg.SessionHeaderName, cfg.store, cfg.SessionExpired)
	}

	middles := append([]mid.GinMiddle{
		mid.NewGinBindUserMid(
			mid.BindUserMidWithCtxKey[T](cfg.CtxUserKey),
			mid.BindUserMidWithBindObject(cfg.bindUser),
		)}, cfg.getMiddles()...)

	server = server.Middles(
		middles...).
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

func AutoGinApiRunWithBindUser[T mid.BindUser](ctx context.Context, cfg *ConfigWithBindUser[T]) error {
	var apiWait sync.WaitGroup
	server, err := autoGinApiServerWithBindUser(cfg)
	if err != nil {
		return err
	}
	const fiveSecods = 5 * time.Second
	apiWait.Add(1)
	go func(srv *http.Server) {

		for {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				cfg.Logger.Fatalf("listen: %s", err)
				time.Sleep(fiveSecods)
			} else if err == http.ErrServerClosed {
				apiWait.Done()
				return
			}
		}
	}(server)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecods)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger.Fatalf("Server forced to shutdown: %v", err)
	}
	apiWait.Wait()
	return nil
}
