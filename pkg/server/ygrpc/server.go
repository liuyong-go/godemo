package ygrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/server"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	listener net.Listener
	Config   *conf.YGrpc
}

//NewServer 初始化服务
func NewServer() *Server {
	config := conf.NewYGrpc().GetConfig()
	fmt.Println(config)
	newServer := grpc.NewServer(config.ServerOptions...)
	listener, err := net.Listen(config.Network, config.ServiceAddress)
	if err != nil {
		fmt.Println("grpc LISTENER error", err)
	}

	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
	}
}

//Serve 实现接口，启动服务
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	return err
}

//Stop 实现接口 ，停止服务
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

//GracefulStop 实现接口 优雅停止服务
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddr = s.Config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme(s.Config.Scheme),
		server.WithAddress(serviceAddr),
		server.WithName(s.Config.Name),
		server.WithEnable(s.Config.EnableService),
	)
	return &info
}
