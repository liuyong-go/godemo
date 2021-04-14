package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type YTrace struct {
	ServiceName       string `toml:"service_name"`
	LocalAgentAddress string `toml:"local_agent_address"`
	LogSpans          bool   `toml:"log_spans"`
}

func (dma *YTrace) getTomlFile() string {
	return "trace.toml"
}
func NewYTrace() *YTrace {
	return &YTrace{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *YTrace) GetConfig() *YTrace {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	return c
}
