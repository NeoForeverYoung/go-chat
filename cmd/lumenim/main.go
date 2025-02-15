package main

import (
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/comet"
	"go-chat/internal/mission"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/logger"
	_ "go-chat/internal/pkg/server"

	"github.com/urfave/cli/v2"
)

// Version 服务版本号（默认）
// 构建时传入版本号
// go build -o lumenim -ldflags "-X main.Version=${IMAGE_TAG}" ./cmd/lumenim
var Version = "1.0.0"

func main() {
	app := core.NewApp(Version)
	app.Register(NewHttpCommand)
	app.Register(NewCometCommand)
	app.Register(NewCrontabCommand)
	app.Register(NewQueueCommand)
	app.Register(NewTempCommand)
	app.Register(NewMigrateCommand)
	app.Run()
}

func NewHttpCommand() core.Command {
	return core.Command{
		Name:  "http",
		Usage: "Http Command - Http API 接口服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			// 初始化日志系统
			// - 日志文件路径为 app.log
			// - 日志级别为 Info
			// - 日志标识为 "http"
			logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "http")
			// 启动 HTTP 服务
			// NewHttpInjector(conf) 创建依赖注入器
			// apis.Run 启动 HTTP API 服务
			return apis.Run(ctx, NewHttpInjector(conf))
		},
	}
}

func NewCometCommand() core.Command {
	return core.Command{
		Name:  "comet",
		Usage: "Comet Command - Websocket、TCP 服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "comet")
			return comet.Run(ctx, NewCometInjector(conf))
		},
	}
}

func NewCrontabCommand() core.Command {
	return core.Command{
		Name:  "crontab",
		Usage: "Crontab Command - 定时任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "crontab")
			return mission.Cron(ctx, NewCronInjector(conf))
		},
	}
}

func NewQueueCommand() core.Command {
	return core.Command{
		Name:  "queue",
		Usage: "Queue Command - 队列任务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "group",
				Usage: "分组",
				Value: "default",
			},
		},
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "queue")
			return mission.Queue(ctx, NewQueueInjector(conf))
		},
	}
}

func NewMigrateCommand() core.Command {
	return core.Command{
		Name:  "migrate",
		Usage: "Migrate Command - 数据库初始化",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "migrate")
			return mission.Migrate(ctx, NewMigrateInjector(conf))
		},
	}
}

func NewTempCommand() core.Command {
	return core.Command{
		Name:  "temp",
		Usage: "Temp Command - 临时命令",
		Subcommands: []core.Command{
			{
				Name:  "test",
				Usage: "Test Command",
				Action: func(ctx *cli.Context, conf *config.Config) error {
					logger.Init(conf.Log.LogFilePath("app.log"), logger.LevelInfo, "temp")
					return NewTempInjector(conf).TestCommand.Do(ctx)
				},
			},
		},
	}
}
