package logging

import (
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/log"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

func NewTemporalAdapter(logger *logrus.Entry) log.Logger {
	return logur.LoggerToKV(logrusadapter.NewFromEntry(logger))
}
