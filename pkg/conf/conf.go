package conf

var (
	ConfPath   = "godemo/demos/toml/"
	RemotePath = "http://liuyong.m.soyoung.com/"
	YDB        *Database
	YServer    *Server
	YLog       *YLogConfig
	YGinHttp   *GinHttp
)

func InitConfig() {
	YDB = newDatabase().getRemoteConfig()
	YServer = newServer().getConfig()
	YLog = newLog().getConfig()
	YGinHttp = newGinHttp().getConfig()
}
