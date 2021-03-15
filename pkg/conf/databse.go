package conf

import (
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
)

type Database struct {
	tomlFile string
	DB       map[string]db `toml:"mysql"`
}
type db struct {
	Host     string
	Port     int
	Username string
	Password string
	schema   string
}

func (dma *Database) getTomlFile() string {
	return "database.toml"
}

func NewDatabase() *Database {
	return &Database{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (dma *Database) GetConfig() *Database {
	if _, err := toml.DecodeFile(ConfPath+dma.getTomlFile(), &dma); err != nil {
		fmt.Println(err)
		return dma
	}
	return dma
}

//GetRemoteConfig 读取远程文件 dbConf := conf.NewDatabase().GetRemoteConfig()
func (dma *Database) GetRemoteConfig() *Database {
	remotePath := RemotePath + dma.getTomlFile()
	resp, err := http.Get(remotePath)
	if err != nil {
		fmt.Println(err)
		return dma
	}
	defer resp.Body.Close()
	if _, err := toml.DecodeReader(resp.Body, &dma); err != nil {
		fmt.Println(err)
		return dma
	}
	return dma
}
