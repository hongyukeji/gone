package logger

import (
	"errors"
	"fmt"
	"iki-go/file"
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
		file.CreateDirs(baseFolder)
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

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.doPrintf(WarnLevel, printWarnLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}

func (logger *Logger) Panic(format string, a ...interface{}) {
	logger.doPrintf(PanicLevel, printPanicLevel, format, a...)
}

var gLogger, _ = New("debug", ConsoleMode, "", log.LstdFlags)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.doPrintf(WarnLevel, printWarnLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}
func Panic(format string, a ...interface{}) {
	gLogger.doPrintf(PanicLevel, printPanicLevel, format, a...)
}

func Close() {
	gLogger.Close()
}
