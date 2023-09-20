package logging

import (
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/log"
	"logur.dev/logur"

	logrusadapter "logur.dev/adapter/logrus"
)

func NewTemporalAdapter(logger *logrus.Entry) log.Logger {
	return logur.LoggerToKV(logrusadapter.NewFromEntry(logger))
}
