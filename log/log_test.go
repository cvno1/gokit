package log_test

import (
	"errors"
	"testing"

	"github.com/cvno1/gokit/log"
)

func TestLog(t *testing.T) {
	logger := log.New()
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log.NewField("tests", 1))
	logger.Error("This is an error log", log.NewField("err", errors.New("exampleerr")))
}

func TestWithField(t *testing.T) {

	logger := log.New(
		log.WithField("testkey", 1),
		log.WithField("testkey2", 1),
	)
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log.NewField("tests", 1))
	logger.Error("This is an error log", log.Error(errors.New("exampleerr")))
}

func TestWithDisableCaller(t *testing.T) {
	logger := log.New(log.WithDisableCaller())
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log.NewField("tests", 1))
	logger.Error("This is an error log", log.Error(errors.New("exampleerr")))
}

func TestWithDisableStackTrace(t *testing.T) {
	logger := log.New(log.WithDisableStackTrace())
	defer logger.Flush()
	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log.NewField("tests", 1))
	logger.Error("This is an error log", log.Error(errors.New("exampleerr")))
}
