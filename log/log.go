package log

import "go.uber.org/zap"

var logger *zap.Logger

func Init() error {
	l, err := zap.NewDevelopment()
	logger = l

	return err
}

func Close() error {
	return logger.Sync()
}

func Get() *zap.Logger {
	return logger
}
