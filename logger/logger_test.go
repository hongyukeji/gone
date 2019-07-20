package logger_test

import (
	"fmt"
	"github.com/wx11055/gone/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

/**
官方go包使用
*/

func TestLog(t *testing.T) {
	filename := "test.log"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	debuglog := log.New(file, "[INFO] ", log.Llongfile)
	debuglog.Println("a debug message here")
	debuglog.SetPrefix("[DEBUG] ")
	debuglog.Println("a debug message here")

	panicln()
	fatal()
}
func panicln() {
	recov()
	log.Panicln("test for defer panic")
	fmt.Println("--------end---------")
	defer func() {
		fmt.Println("second")
	}()
}
func fatal() {
	defer func() {
		fmt.Println("first")
	}()
	log.Fatalln("test for defer Fatal")
}
func recov() {
	defer func() {
		fmt.Println("defer panic first")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}
func TestZap(t *testing.T) {
	log, _ := zap.NewProduction()
	url := "Hello"
	recov()
	defer log.Sync()
	log.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backpff", time.Second))
	log.Warn("debug log", zap.String("level", url))
	log.Error("error message", zap.String("error", url))
	log.Panic("Panic log", zap.String("level", url))
}
func TestC(t *testing.T) {

	alevel := zap.NewAtomicLevel()
	// curl  http://localhost:9090/handle/level
	//curl -XPUT --data '{"level":"debug"}' http://localhost:9090/handle/level

	http.HandleFunc("/handle/level", alevel.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			panic(err)
		}
	}()
	cf := zap.NewProductionConfig()
	cf.Level = alevel
	log, err := cf.Build()
	if err != nil {
		fmt.Println("err", err)
	}
	defer log.Sync()
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Second)
		log.Debug("debug log ", zap.String("level", alevel.String()))
		log.Info("info log  ", zap.String("level", alevel.String()))
	}
}
func TestLumberjack(t *testing.T) {
	log := initLogger("all.log", "info")
	log.Info(" test log", zap.Int("line", 47))
	log.Warn(" testlog", zap.Int("line", 47))
}
func initLogger(logpath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    1024,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}
	w := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	cf := zap.NewProductionEncoderConfig()

	cf.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cf), w, level)
	log := zap.New(core)
	log.Info("test.log", zap.Int("line", 47))
	log.Warn("testlog", zap.Int("line", 47))
	return log
}

func TestLogger(t *testing.T) {
	name := "MATOSIKI_LOGGER"

	logger.Debug("My name is %v", name)
	logger.Info("My name is %v", name)
	logger.Error("My name is %v", name)

	l, err := logger.New("info", logger.FileMode, "test", log.LstdFlags)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer l.Close()

	l.Debug("will not print")
	l.Info("My name is %v", name)

	logger.Export(l)

	l.Debug("will not print")
	l.Info("My name is %v", name)
}
