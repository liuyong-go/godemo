package core

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/server"
	"github.com/liuyong-go/godemo/pkg/util/signals"
	"github.com/liuyong-go/godemo/pkg/util/ycycle"
	"github.com/liuyong-go/godemo/pkg/util/ydefer"
	"github.com/liuyong-go/godemo/pkg/util/ygo"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	smu         *sync.RWMutex //读写锁
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	servers     []server.Server
	hooks       map[uint32]*ydefer.DeferStack //注册钩子函数
	cycle       *ycycle.Cycle
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
		app.cycle = ycycle.NewCycle()
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
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()
	app.waitSignals() //开启协程监听信号
	defer app.clean()
	app.cycle.Run(app.startServers)
	if err := <-app.cycle.Wait(); err != nil {
		ylog.SugarLogger.Info("shutdown with error", err)
		return err
	}
	ylog.SugarLogger.Info("shut down,bye")
	return nil
}
func (app *Application) startServers() error {
	fmt.Println("start server")
	var eg errgroup.Group
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			ylog.SugarLogger.Info("start server", s.Info().Name)
			defer ylog.SugarLogger.Info("end server", s.Info().Name)
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

// waitSignals wait signal
func (app *Application) waitSignals() {
	ylog.SugarLogger.Info("init listen signal")
	signals.Shutdown(func(grace bool) {
		if grace {
			app.GracefulStop(context.TODO())
		} else {
			app.Stop()
		}
	})
}
func (app *Application) clean() {
	ylog.Logger.Sync()
}
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.Unlock()
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return
}
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.smu.Unlock()
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return
}
