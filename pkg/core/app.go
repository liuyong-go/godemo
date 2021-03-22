package core

import (
	"sync"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/server"
	"github.com/liuyong-go/godemo/pkg/util/ygo"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
)

type Application struct {
	smu         *sync.RWMutex //读写锁
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	servers     []server.Server
}

func NewApp() *Application {
	return &Application{smu: &sync.RWMutex{}, servers: make([]server.Server, 0)}
}

func (app *Application) startup() (err error) {
	app.startupOnce.Do(func() {
		err = ygo.SerialUntilError(
			//初始化方法
			app.loadConfig,
			app.loadLog,
		)()
	})
	return
}
func (app *Application) Start(fns ...func() error) error {
	if err := app.startup(); err != nil {
		return err
	}
	return ygo.SerialUntilError(fns...)()
}

//初始解析flag
func (app *Application) loadFlag() error {
	return nil
}

//初始化配置文件
func (app *Application) loadConfig() error {
	conf.InitConfig()
	return nil
}

//初始化日志
func (app *Application) loadLog() error {
	ylog.InitLog()
	return nil
}
func (app *Application) Serve(s ...server.Server) {

}
