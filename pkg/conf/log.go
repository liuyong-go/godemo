package conf

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type YLogConfig struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int `toml:"max_size"` //M
	MaxAge    int `toml:"max_age"`  //day
	MaxBackup int `toml:"max_backup"`
	// 日志磁盘刷盘间隔
	Interval      time.Duration //*time.Hour
	CallerSkip    int           `toml:"caller_skip"`
	Async         bool
	Queue         bool
	QueueSleep    time.Duration `toml:"queue_sleep"`
	Core          zapcore.Core
	Debug         bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
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
	c.Interval = c.Interval * time.Hour
	c.QueueSleep = c.QueueSleep * time.Microsecond
	return c
}
