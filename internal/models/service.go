package models

import (
	"time"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/durationpb"
)

const defaultRunCount = 1

type Service struct {
	bun.BaseModel `bun:"services,alias:s"`

	ID   int    `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull,unique"`

	CheckerType checkerpb.Type `bun:"checker_type,notnull"`
	CheckerPath string         `bun:"checker_path,notnull"`

	DefaultScore float64 `bun:"default_score,notnull"`

	DefaultTimeout time.Duration                      `bun:"default_timeout,notnull"`
	ActionTimeouts map[checkerpb.Action]time.Duration `bun:"action_timeouts"`

	ActionRunCounts map[checkerpb.Action]int `bun:"action_run_counts"`

	// TODO: vulns format.
	// Places int
}

func (s *Service) String() string {
	return s.Name
}

func (s *Service) CheckerTimeout(action checkerpb.Action) time.Duration {
	if timeout, ok := s.ActionTimeouts[action]; ok {
		return timeout
	}
	return s.DefaultTimeout
}

func (s *Service) RunCount(action checkerpb.Action) int {
	if count, ok := s.ActionRunCounts[action]; ok {
		return count
	}
	return defaultRunCount
}

func (s *Service) ToProto() *servicespb.Service {
	return &servicespb.Service{
		Id:   int32(s.ID),
		Name: s.Name,
		Checker: &servicespb.Service_Checker{
			Type:           s.CheckerType,
			Path:           s.CheckerPath,
			DefaultTimeout: durationpb.New(s.DefaultTimeout),
			ActionTimeouts: lo.MapToSlice(s.ActionTimeouts, func(action checkerpb.Action, timeout time.Duration) *servicespb.Service_Checker_ActionTimeout {
				return &servicespb.Service_Checker_ActionTimeout{
					Action:  action,
					Timeout: durationpb.New(timeout),
				}
			}),
			ActionRunCounts: lo.MapToSlice(s.ActionRunCounts, func(action checkerpb.Action, count int) *servicespb.Service_Checker_ActionRunCount {
				return &servicespb.Service_Checker_ActionRunCount{
					Action:   action,
					RunCount: int32(count),
				}
			}),
		},

		DefaultScore: s.DefaultScore,
	}
}

func NewServiceFromProto(p *servicespb.Service) *Service {
	return &Service{
		ID:   int(p.Id),
		Name: p.Name,

		CheckerType:    p.Checker.Type,
		CheckerPath:    p.Checker.Path,
		DefaultTimeout: p.Checker.DefaultTimeout.AsDuration(),
		ActionTimeouts: lo.SliceToMap(
			p.Checker.ActionTimeouts,
			func(t *servicespb.Service_Checker_ActionTimeout) (checkerpb.Action, time.Duration) {
				return t.Action, t.Timeout.AsDuration()
			},
		),
		ActionRunCounts: lo.SliceToMap(
			p.Checker.ActionRunCounts,
			func(t *servicespb.Service_Checker_ActionRunCount) (checkerpb.Action, int) {
				return t.Action, int(t.RunCount)
			},
		),

		DefaultScore: p.DefaultScore,
	}
}
