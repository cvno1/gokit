package log_test

import (
	"errors"
	"testing"

	log2 "github.com/cvno1/gokit/pkg/log"
)

func TestLog(t *testing.T) {
	logger := log2.New()
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log2.NewField("tests", 1))
	logger.Error("This is an error log", log2.NewField("err", errors.New("exampleerr")))
}

func TestWithField(t *testing.T) {

	logger := log2.New(
		log2.WithField("testkey", 1),
		log2.WithField("testkey2", 1),
	)
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log2.NewField("tests", 1))
	logger.Error("This is an error log", log2.Error(errors.New("exampleerr")))
}

func TestWithDisableCaller(t *testing.T) {
	logger := log2.New(log2.WithDisableCaller())
	defer logger.Flush()

	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log2.NewField("tests", 1))
	logger.Error("This is an error log", log2.Error(errors.New("exampleerr")))
}

func TestWithDisableStackTrace(t *testing.T) {
	logger := log2.New(log2.WithDisableStackTrace())
	defer logger.Flush()
	logger.Debug("This is a debug log")
	logger.Info("This is an info log", log2.NewField("tests", 1))
	logger.Error("This is an error log", log2.Error(errors.New("exampleerr")))
}
