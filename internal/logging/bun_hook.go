package logging

import (
	"time"

	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

func AddBunQueryHook(db *bun.DB) {
	db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
		Logger:     logrus.StandardLogger(),
		QueryLevel: logrus.DebugLevel,
		ErrorLevel: logrus.ErrorLevel,
		SlowLevel:  logrus.WarnLevel,
		LogSlow:    time.Millisecond * 100,
	}))
}
