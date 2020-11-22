package logging

import (
	"bytes"
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

type OutputSplitter struct{}

func (splitter *OutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.HasPrefix(p, []byte("[ERROR]")) || bytes.HasPrefix(p, []byte("[FATAL]")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

func New(cfg *calcapp.Config) (*logrus.Logger, error) {
	logger := &logrus.Logger{
		Out:   &OutputSplitter{},
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}

	if cfg.LogsToFile != nil {
		f, err := os.OpenFile(*cfg.LogsToFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			return nil, fmt.Errorf("opening log file %s error %w", *cfg.LogsToFile, err)
		}
		logger.SetOutput(f)
	}

	if cfg.LogsToJSON {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	switch cfg.LogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	}

	return logger, nil
}
