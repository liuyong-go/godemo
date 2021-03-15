package conf

import (
	"fmt"

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
