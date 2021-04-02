package conf

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type EtcdV3 struct {
	Endpoints []string `json:"endpoints"`
	CertFile  string   `json:"certFile" toml:"cert_file"`
	KeyFile   string   `json:"keyFile" toml:"key_file"`
	CaCert    string   `json:"caCert" toml:"ca_cert"`
	BasicAuth bool     `json:"basicAuth" toml:"basic_auth"`
	UserName  string   `json:"userName" toml:"user_name"`
	Password  string   `json:"-" toml:"password"`
	// 连接超时时间
	ConnectTimeout time.Duration `json:"connectTimeout" toml:"connect_timeout"`
	Secure         bool          `json:"secure"`
	// 自动同步member list的间隔
	AutoSyncInterval time.Duration `json:"autoAsyncInterval" toml:"auto_sync_interval"`
	ServicePrefix    string        `toml:"service_prefix"`
}

func (dma *EtcdV3) getTomlFile() string {
	return "etcdv3.toml"
}
func NewEtcdV3() *EtcdV3 {
	return &EtcdV3{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *EtcdV3) GetConfig() *EtcdV3 {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	c.AutoSyncInterval = c.AutoSyncInterval * time.Second
	c.ConnectTimeout = c.ConnectTimeout * time.Second
	return c
}
