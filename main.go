package main

import (
	"fmt"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/core"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"go.uber.org/zap"
)

type App struct {
	core.Application
}

func main() {
	conf.ConfPath = "/Users/liuyong/go/src/godemo/toml/toml_dev/"
	var app = core.NewApp()
	app.Start()

	defer ylog.Logger.Sync()
	fmt.Println(conf.YDB)
	ylog.SugarLogger.Infow("测试日志", zap.String("name", "测试"))
	fmt.Println("输出内容")

}
func (app *App) serverHttp() {

}
