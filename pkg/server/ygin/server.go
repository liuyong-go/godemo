//返回gin服务
package ygin

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/server"
	"github.com/liuyong-go/godemo/pkg/util/ylog"
)

type Server struct {
	*gin.Engine
	Server   *http.Server
	Listener net.Listener
}

func NewServer() *Server {
	listener, err := net.Listen("tcp", Address())
	if err != nil {
		ylog.SugarLogger.Panic("new gin server err", err)
	}
	gin.SetMode(conf.YGinHttp.Mode)
	return &Server{
		Engine:   gin.New(),
		Listener: listener,
	}

}

//Serve 实现 server.Server 接口
func (s *Server) Serve() error {
	for _, route := range s.Engine.Routes() {
		ylog.SugarLogger.Info("add route", route.Method, route.Path)
	}
	s.Server = &http.Server{
		Addr:    Address(),
		Handler: s,
	}
	err := s.Server.Serve(s.Listener)
	if err == http.ErrServerClosed {
		ylog.SugarLogger.Info("close gin", Address())
		return nil
	}
	return err
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.Listener.Addr().String()
	if conf.YGinHttp.ServiceAddress != "" {
		serviceAddr = conf.YGinHttp.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme(conf.YGinHttp.Scheme),
		server.WithAddress(serviceAddr),
		server.WithName(conf.YGinHttp.Name),
		server.WithEnable(conf.YGinHttp.EnableService),
	)
	return &info
}
func Address() string {
	return conf.YGinHttp.ServiceAddress
}
