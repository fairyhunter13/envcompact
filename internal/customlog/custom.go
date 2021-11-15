package customlog

import (
	"go.uber.org/zap"
)

var instance *zap.Logger

func Set(logger *zap.Logger) {
	instance = logger
}

func Get() *zap.Logger {
	return instance
}
