package app

import "github.com/liuyong-go/godemo/pkg/conf"

var (
	DB     *conf.Database
	Server *conf.Server
	YLog   *conf.YLogConfig
)

type TomlConfig struct {
	TomlFile string
}
type Config struct {
}

func InitConfig() {
	DB = conf.NewDatabase().GetRemoteConfig()
	Server = conf.NewServer().GetConfig()
	YLog = conf.NewLog().GetConfig()
}

//通过请求参数reload配置文件热加载方法
func reloadConfig() {

}
