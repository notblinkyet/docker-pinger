package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	srv *http.Server
}

func New(router *gin.Engine, port int, host string, timeout time.Duration) *App {
	return &App{
		srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      router,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}
}

func (a *App) Run() error {
	return a.srv.ListenAndServe()
}

func (a *App) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return a.srv.Shutdown(ctx)
}
