package util

import (
	"github.com/sirupsen/logrus"
)

func Must[T any](message string) func(t T, err error) T {
	return func(t T, err error) T {
		if err != nil {
			logrus.Fatalf("%s: %v", message, err)
		}
		return t
	}
}
