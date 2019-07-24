package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	DebugLevel int = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)
const (
	printDebugLevel = "[debug  ] "
	printInfoLevel  = "[info   ] "
	printWarnLevel  = "[warn   ] "
	printErrorLevel = "[error  ] "
	printFatalLevel = "[fatal  ] "
	printPanicLevel = "[panic   ] "
)

type PrintMode int

const (
	ConsoleMode PrintMode = iota //控制台打印
	FileMode                     //文件输出
)

type Logger struct {
	level      int
	logTime    int64
	baseFolder string
	baseFile   *os.File
	baseLogger *log.Logger
}

var logger Logger

func fileMode(baseFolder string, flag int) (baseLogger *log.Logger, baseFile *os.File, err error) {
	now := time.Now()
	filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())

	file, err := os.Create(path.Join(baseFolder, filename))
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(path.Join(baseFolder, filename))
	baseLogger = log.New(file, "", flag)
	baseFile = file
	return
}

func consoleMode(baseFolder string, flag int) (baseLogger *log.Logger) {
	baseLogger = log.New(os.Stdout, "", flag)
	return
}

func New(strLevel string, mode PrintMode, baseFolder string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	case "panic":
		level = PanicLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}
	var baseLogger *log.Logger
	var baseFile *os.File
	var err error
	logger := new(Logger)
	logger.level = level
	if mode == FileMode {
		if !IsExist(baseFolder) {
			err := os.MkdirAll(baseFolder, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		baseLogger, baseFile, err = fileMode(baseFolder, flag)
		if err != nil {
			return nil, err
		}
	} else {
		baseLogger = consoleMode(baseFolder, flag)
	}
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile

	return logger, nil
}
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == FatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}
func (logger *Logger) Debug(a ...interface{}) {
	logger.doPrintf(DebugLevel, printWarnLevel, "%v", a...)
}

func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}
func (logger *Logger) Info(a ...interface{}) {
	logger.doPrintf(InfoLevel, printWarnLevel, "%v", a...)
}
func (logger *Logger) Warnf(format string, a ...interface{}) {
	logger.doPrintf(WarnLevel, printWarnLevel, format, a...)
}

func (logger *Logger) Warn(a ...interface{}) {
	logger.doPrintf(WarnLevel, printWarnLevel, "%v", a...)
}

func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}
func (logger *Logger) Error(a ...interface{}) {
	logger.doPrintf(ErrorLevel, printWarnLevel, "%v", a...)
}

func (logger *Logger) Fatalf(format string, a ...interface{}) {
	logger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}
func (logger *Logger) Fatal(a ...interface{}) {
	logger.doPrintf(FatalLevel, printWarnLevel, "%v", a...)
}

func (logger *Logger) Panicf(format string, a ...interface{}) {
	logger.doPrintf(PanicLevel, printPanicLevel, format, a...)
}
func (logger *Logger) Panic(a ...interface{}) {
	logger.doPrintf(PanicLevel, printWarnLevel, "%v", a...)
}

var gLogger, _ = New("debug", ConsoleMode, "", log.LstdFlags)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}
func Debug(a ...interface{}) {
	gLogger.Debug(a...)
}
func Info(a ...interface{}) {
	gLogger.Info(a...)
}
func Warn(a ...interface{}) {
	gLogger.Warn(a...)
}
func Error(a ...interface{}) {
	gLogger.Error(a...)
}
func Fatal(a ...interface{}) {
	gLogger.Fatal(a...)
}
func Panic(a ...interface{}) {
	gLogger.Panic(a...)
}

func Debugf(format string, a ...interface{}) {
	gLogger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}

func Infof(format string, a ...interface{}) {
	gLogger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}

func Warnf(format string, a ...interface{}) {
	gLogger.doPrintf(WarnLevel, printWarnLevel, format, a...)
}

func Errorf(format string, a ...interface{}) {
	gLogger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}

func Fatalf(format string, a ...interface{}) {
	gLogger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}
func Panicf(format string, a ...interface{}) {
	gLogger.doPrintf(PanicLevel, printPanicLevel, format, a...)
}

func Close() {
	gLogger.Close()
}
