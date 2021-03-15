package conf

var (
	ConfPath = "godemo/demos/toml/"
	Conf     Config
)

type TomlConfig struct {
	TomlFile string
}
type Config struct {
	DB     *Database
	Server *Server
}

func InitConfig() {
	Conf.DB = NewDatabase().GetConfig()
	Conf.Server = NewServer().GetConfig()
}
