package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type GinHttp struct {
	Host string
	Port int
	Mode string //debug, release,test
	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string `toml:"service_address"`
}

func (dma *GinHttp) getTomlFile() string {
	return "gin_server.toml"
}
func newGinHttp() *GinHttp {
	return &GinHttp{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *GinHttp) getConfig() *GinHttp {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	return c
}
