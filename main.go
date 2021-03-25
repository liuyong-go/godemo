package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/core"
	"github.com/liuyong-go/godemo/pkg/server/ygin"
	"github.com/liuyong-go/godemo/pkg/util/ycycle"
	"github.com/liuyong-go/godemo/pkg/util/ydefer"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
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
	fmt.Println("输出内容")

}
func testCycle() {
	state := "init"
	c := ycycle.NewCycle()
	c.Run(func() error {
		time.Sleep(time.Second * 2)
		return nil
	})
	go func() {
		select {
		case <-c.Done():
			state = "done"
		case <-time.After(time.Second * 5):
			state = "close"
		}
		c.Close()
	}()
	x := <-c.Wait()
	fmt.Println("x quit", x)
	fmt.Println("result", state)

}
func testSwap() {
	var value int32 = 10
	if atomic.CompareAndSwapInt32(&value, 12, 11) {
		fmt.Println("修改值", value)
	} else {
		fmt.Println("未修改值", value)
	}
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
