package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/core"
	"github.com/liuyong-go/godemo/pkg/server/ygin"
	"github.com/liuyong-go/godemo/pkg/util/ydefer"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"go.uber.org/zap"
)

type App struct {
	core.Application
}

func main() {
	conf.ConfPath = "/Users/liuyong/go/src/godemo/toml/toml_dev/"
	var app = &App{}
	app.Start(app.serverHttp)
	app.RegisterHooks(core.StageAfterStop, app.afterStop)
	app.RegisterHooks(core.StageBeforeStop, app.beforeStop)
	app.Run()
	defer ydefer.Clean()
	defer ylog.Logger.Sync()
	ylog.SugarLogger.Infow("测试日志", zap.String("name", "测试"))
	fmt.Println("输出内容")

}
func (app *App) serverHttp() error {
	server := ygin.NewServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, "HELLO YONG")
	})
	return app.Serve(server)
}
func (app *App) beforeStop() error {
	ylog.SugarLogger.Infow("before stop")
	return nil
}
func (app *App) afterStop() error {
	ylog.SugarLogger.Infow("after stop")
	return nil
}
