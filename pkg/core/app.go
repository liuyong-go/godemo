package core

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/server"
	"github.com/liuyong-go/godemo/pkg/util/ydefer"
	"github.com/liuyong-go/godemo/pkg/util/ygo"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
)

type Application struct {
	smu         *sync.RWMutex //读写锁
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	servers     []server.Server
	hooks       map[uint32]*ydefer.DeferStack //注册钩子函数
}

const (
	//StageAfterStop after app stop
	StageAfterStop uint32 = iota + 1
	//StageBeforeStop before app stop
	StageBeforeStop
)

//初始化hook函数
func (app *Application) initHooks(hookKeys ...uint32) {
	app.hooks = make(map[uint32]*ydefer.DeferStack, len(hookKeys))
	for _, k := range hookKeys {
		app.hooks[k] = ydefer.NewStack()
	}
}

//执行钩子函数
func (app *Application) runHooks(k uint32) {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Clean()
	}
}
func (app *Application) RegisterHooks(k uint32, fns ...func() error) error {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Push(fns...)
		return nil
	}
	return fmt.Errorf("hook stage not found" + strconv.FormatInt(int64(k), 10))
}

//初始化
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		//assign
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		//private method
		app.initHooks(StageBeforeStop, StageAfterStop)
	})
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
	app.initialize()
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
func (app *Application) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}
func (app *Application) Run(servers ...server.Server) error {
	app.runHooks(StageBeforeStop)
	app.runHooks(StageAfterStop)
	return nil
}
