package conf

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

// Config ...
type YLogConfig struct {
	Development  bool         //开发，线上环境
	InfoPath     string       `toml:"info_path"`
	ErrorPath    string       `toml:"error_path"`
	EncodeConfig EncodeConfig `toml:"encode_config"`
	RotationLogs RotateLogs
}
type RotateLogs struct {
	MaxAge       time.Duration `toml:"max_age"`       // 保存小时数
	RotationTime time.Duration `toml:"rotation_time"` //切割频率 小时记录
}
type EncodeConfig struct {
	MessageKey string `toml:"message_key"`
	LevelKey   string `toml:"level_key"`
	TimeKey    string `toml:"time_key"`
	CallerKey  string `toml:"caller_key"`
}

func (c *YLogConfig) getTomlFile() string {
	return "log.toml"
}

func NewLog() *YLogConfig {
	return &YLogConfig{}
}

//GetConfig dbConf := conf.NewDatabase().GetConfig()
func (c *YLogConfig) GetConfig() *YLogConfig {
	if _, err := toml.DecodeFile(ConfPath+c.getTomlFile(), &c); err != nil {
		fmt.Println(err)
		return c
	}
	c.RotationLogs.MaxAge = c.RotationLogs.MaxAge * time.Hour
	return c
}
