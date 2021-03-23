package ylog

import (
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/liuyong-go/godemo/pkg/util/ydefer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//只能输出结构化日志，但是性能要高于 SugaredLogger
	Logger *zap.Logger
	//可以输出 结构化日志、非结构化日志。性能低于zap.Logger，
	SugarLogger *zap.SugaredLogger
)

func InitLog() {
	// config := zapcore.EncoderConfig{}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   conf.YLog.EncodeConfig.MessageKey,
		LevelKey:     conf.YLog.EncodeConfig.LevelKey,  //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      conf.YLog.EncodeConfig.TimeKey,   //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    conf.YLog.EncodeConfig.CallerKey, //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder,      //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,       //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= zapcore.InfoLevel
	})
	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	var wsInfo zapcore.WriteSyncer
	var wsWarn zapcore.WriteSyncer
	var core zapcore.Core
	if conf.YLog.Development == true {
		wsInfo = os.Stdout
	} else {
		wsInfo = zapcore.AddSync(getWriter(conf.YLog.InfoPath))
		wsWarn = zapcore.AddSync(getWriter(conf.YLog.ErrorPath))
	}
	if conf.YLog.Async == true {
		var close CloseFunc
		wsInfo, close = Buffer(wsInfo, defaultBufferSize, defaultFlushInterval)
		wsWarn, close = Buffer(wsWarn, defaultBufferSize, defaultFlushInterval)
		ydefer.Register(close)
	}
	if conf.YLog.Development == true {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), wsInfo, infoLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), wsInfo, infoLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), wsWarn, warnLevel),
		)
	}
	//实现多个输出

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	SugarLogger = Logger.Sugar()
}
func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*conf.YLog.RotationLogs.MaxAge),             // 保存30天
		rotatelogs.WithRotationTime(time.Hour*conf.YLog.RotationLogs.RotationTime), //切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}
