package conf

var (
	ConfPath   = "godemo/demos/toml/"
	RemotePath = "http://liuyong.m.soyoung.com/"
	Conf       Config
)

type TomlConfig struct {
	TomlFile string
}
type Config struct {
	DB     *Database
	Server *Server
	YLog   *YLogConfig
}

func InitConfig() {
	Conf.DB = NewDatabase().GetRemoteConfig()
	Conf.Server = NewServer().GetConfig()
	Conf.YLog = NewLog().GetConfig()
}

//通过请求参数reload配置文件热加载方法
func reloadConfig() {

}
