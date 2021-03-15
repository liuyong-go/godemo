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
}

func InitConfig() {
	Conf.DB = NewDatabase().GetRemoteConfig()
	Conf.Server = NewServer().GetConfig()
}
