package logger_test

import (
	"github.com/wx11055/gone/logger"
	"log"
	"testing"
)

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
