package models

import (
	"time"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type Service struct {
	ID   int
	Name string

	CheckerType           checkerpb.Type
	CheckerPath           string
	CheckerTimeoutSeconds int

	DefaultScore int

	Gets int
	Puts int

	// TODO: vulns format.
	// Places int
}

func (s *Service) String() string {
	return s.Name
}

func (s *Service) CheckerTimeout() time.Duration {
	return time.Duration(s.CheckerTimeoutSeconds) * time.Second
}
