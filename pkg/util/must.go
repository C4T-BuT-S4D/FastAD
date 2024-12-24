package util

import "go.uber.org/zap"

func Must[T any](message string) func(t T, err error) T {
	return func(t T, err error) T {
		if err != nil {
			zap.L().Fatal(message, zap.Error(err))
		}
		return t
	}
}
