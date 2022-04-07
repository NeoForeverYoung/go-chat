package provider

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
)

type HttpServer *http.Server

type WebsocketServer *http.Server

func NewHttpServer(conf *config.Config, handler *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.App.Port),
		Handler: handler,
	}
}

func NewWebsocketServer(handler *gin.Engine) WebsocketServer {
	return &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9504),
		Handler: handler,
	}
}
