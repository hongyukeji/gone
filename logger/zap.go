package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

// Logger 业务逻辑日志记录器
var recoder *zap.Logger

// Logger 服务日志记录器
var logzap *zap.Logger

type ZapConfig struct {
	Level                  string
	Outputs                string
	Encode                 string
	ColorLevel             bool
	EnableTrace            bool
	EnableCaller           bool
	EnableNotFound         bool
	EnableMethodNotAllowed bool
	Debug                  bool
}

var zapConfig ZapConfig

//func init() {
//	zapConfig = ZapConfig{Level: "debug", Outputs: "outputs", Encode: "console", ColorLevel: false, EnableTrace: true, EnableCaller: true, EnableNotFound: false, EnableMethodNotAllowed: false}
//	err := NewZapLogger(zapConfig)
//	if err != nil {
//		fmt.Errorf(err.Error())
//	}
//}

// SetLogger 设置logger
func NewZapLogger(zc ZapConfig) error {
	var err error

	// 设置日志记录级别
	var zapLevel zapcore.Level
	switch zc.Level {
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.DebugLevel
	}

	// 配置级别编码器
	var encodeLevel zapcore.LevelEncoder
	if zc.ColorLevel == true && zc.Encode == "console" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeLevel = zapcore.CapitalLevelEncoder
	}

	outputs := strings.Split(zc.Outputs, "|")

	// 配置编码器的参数
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message", // 消息字段名
		LevelKey:   "level",   // 级别字段名
		TimeKey:    "time",    // 时间字段名
		CallerKey:  "file",    // 记录源码文件的字段名
		// 编码时间字符串的格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel:  encodeLevel,                // 日志级别的编码器
		EncodeCaller: zapcore.ShortCallerEncoder, // Caller的编码器
	}

	var disableStacktrace bool
	var disableCaller bool
	if zc.EnableCaller == false {
		disableStacktrace = true
	}
	if zc.EnableCaller == false {
		disableCaller = true
	}

	// 设置Logger
	recoder, err = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志记录级别
		Development:       zc.Debug,                       // 开发模式
		Encoding:          zc.Encode,                      // 日志格式json/console
		EncoderConfig:     encoderConfig,                  // 编码器配置
		OutputPaths:       outputs,                        // 输出路径
		DisableStacktrace: disableStacktrace,              // 屏蔽堆栈跟踪
		DisableCaller:     disableCaller,                  // 屏蔽调用信息
	}.Build()
	if err != nil {
		return err
	}

	// 设置logzap
	logzap, err = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志记录级别
		Development:       zc.Debug,                       // 开发模式
		Encoding:          zc.Encode,                      // 日志格式json/console
		EncoderConfig:     encoderConfig,                  // 编码器配置
		OutputPaths:       outputs,                        // 输出路径
		DisableStacktrace: true,                           // 屏蔽堆栈跟踪
		DisableCaller:     true,                           // 屏蔽跟踪
	}.Build()
	if err != nil {
		return err
	}
	return err
}
