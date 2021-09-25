// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/service"
	"go-chat/config"
	"go-chat/connect"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, conf *config.Config) *Service {
	auth := &v1.Auth{
		Conf: conf,
	}
	user := &v1.User{}
	download := &v1.Download{}
	index := &open.Index{}
	redis := connect.RedisConnect(ctx, conf)
	wsClient := &cache.WsClient{
		Redis: redis,
	}
	clientService := &service.ClientService{
		WsClient: wsClient,
	}
	wsWs := &ws.Ws{
		ClientService: clientService,
	}
	handlerHandler := &handler.Handler{
		Auth:     auth,
		User:     user,
		Download: download,
		Index:    index,
		Ws:       wsWs,
	}
	engine := router.NewRouter(conf, handlerHandler)
	server := connect.NewHttp(conf, engine)
	serverRunID := cache.NewServerRun(redis)
	socketService := &service.SocketService{
		Conf:        conf,
		ServerRunID: serverRunID,
	}
	mainService := &Service{
		HttpServer:   server,
		SocketServer: socketService,
	}
	return mainService
}

// wire.go:

var providerSet = wire.NewSet(connect.RedisConnect, connect.NewHttp, router.NewRouter, cache.NewServerRun, wire.Struct(new(cache.WsClient), "*"), wire.Struct(new(v1.Auth), "*"), wire.Struct(new(v1.User), "*"), wire.Struct(new(v1.Download), "*"), wire.Struct(new(open.Index), "*"), wire.Struct(new(ws.Ws), "*"), wire.Struct(new(handler.Handler), "*"), wire.Struct(new(service.ClientService), "*"), wire.Struct(new(service.UserService), "*"), wire.Struct(new(service.SocketService), "*"), wire.Struct(new(Service), "*"))
