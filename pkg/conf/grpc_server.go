package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"google.golang.org/grpc"
)

type YGrpc struct {
	*Service
	Network            string
	ServerOptions      []grpc.ServerOption
	StreamInterceptors []grpc.StreamServerInterceptor //流式拦截器
	SnaryInterceptors  []grpc.UnaryServerInterceptor  //一次性拦截器
}

func (dma *YGrpc) getTomlFile() string {
	return "grpc_server.toml"
}
func NewYGrpc() *YGrpc {
	return &YGrpc{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *YGrpc) GetConfig() *YGrpc {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	return c
}
