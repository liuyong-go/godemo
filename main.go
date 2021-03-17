package main

import (
	"fmt"

	"github.com/liuyong-go/godemo/pkg/app"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"go.uber.org/zap"
)

func main() {
	conf.ConfPath = "/Users/liuyong/go/src/godemo/toml/toml_dev/"
	app.InitConfig()
	fmt.Println(app.YLog)
	ylog.InitLog()
	ylog.SugarLogger.Infow("测试日志", zap.String("name", "测试"))
	ylog.SugarLogger.Error("错误", "错误信息")
}
