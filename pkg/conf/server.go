package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Http httpCll `toml:"http"`
	Rpc  rpcCall `toml:"rpc"`
}
type httpCll struct {
	Addr string
	Port int
}
type rpcCall struct {
	Addr string
	Port int
}

func (dma *Server) getTomlFile() string {
	return "server.toml"
}
func NewServer() *Server {
	return &Server{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *Server) GetConfig() *Server {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	return c
}
