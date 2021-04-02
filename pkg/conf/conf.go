package conf

var (
	ConfPath   = "godemo/demos/toml/"
	RemotePath = "http://liuyong.m.soyoung.com/"
	YDB        *Database
	YLog       *YLogConfig
	YGinHttp   *GinHttp
)

type Service struct {
	Name           string
	Scheme         string
	ServiceAddress string `toml:"service_address"`
	EnableService  bool   `toml:"enable_service"`
}

func InitConfig() {
	YDB = newDatabase().getRemoteConfig()
	YLog = newLog().getConfig()
	YGinHttp = newGinHttp().getConfig()
}
