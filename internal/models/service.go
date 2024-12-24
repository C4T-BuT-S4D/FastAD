package models

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/durationpb"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

const defaultRunCount = 1

type ServiceActionConfig struct {
	Timeout  time.Duration `json:"timeout"`
	RunCount int           `json:"run_count"`
}

type Service struct {
	bun.BaseModel `bun:"services,alias:s"`

	ID   int    `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull,unique" json:"name"`

	CheckerType checkerpb.Type `bun:"checker_type,notnull" json:"checker_type,omitempty"`
	CheckerPath string         `bun:"checker_path,notnull" json:"checker_path,omitempty"`

	DefaultScore float64 `bun:"default_score,notnull" json:"default_score,omitempty"`

	DefaultTimeout time.Duration                             `bun:"default_timeout,notnull" json:"default_timeout,omitempty"`
	Actions        map[checkerpb.Action]*ServiceActionConfig `bun:"actions,type:jsonb,notnull" json:"actions,omitempty"`

	Disabled bool `bun:"disabled,notnull" json:"disabled"`

	// TODO: vulns format.
	// Places int
}

func (s *Service) String() string {
	return fmt.Sprintf("Service(name=%s)", s.Name)
}

func (s *Service) CheckerTimeout(action checkerpb.Action) time.Duration {
	if cfg, ok := s.Actions[action]; ok {
		return cfg.Timeout
	}
	return s.DefaultTimeout
}

func (s *Service) GetRunCount(action checkerpb.Action) int {
	if cfg, ok := s.Actions[action]; ok {
		return cfg.RunCount
	}
	return defaultRunCount
}

func (s *Service) ToProto() *servicespb.Service {
	return &servicespb.Service{
		Id:   int64(s.ID),
		Name: s.Name,
		Checker: &servicespb.Service_Checker{
			Type:           s.CheckerType,
			Path:           s.CheckerPath,
			DefaultTimeout: durationpb.New(s.DefaultTimeout),
			Actions: lo.MapToSlice(s.Actions, func(action checkerpb.Action, actionConfig *ServiceActionConfig) *servicespb.Service_Checker_Action {
				return &servicespb.Service_Checker_Action{
					Action:   action,
					RunCount: int64(actionConfig.RunCount),
					Timeout:  durationpb.New(actionConfig.Timeout),
				}
			}),
		},

		DefaultScore: s.DefaultScore,
		Disabled:     s.Disabled,
	}
}

func NewServiceFromProto(p *servicespb.Service) *Service {
	return &Service{
		ID:   int(p.Id),
		Name: p.Name,

		CheckerType:    p.Checker.Type,
		CheckerPath:    p.Checker.Path,
		DefaultTimeout: p.Checker.DefaultTimeout.AsDuration(),
		Actions: lo.SliceToMap(
			p.Checker.Actions,
			func(t *servicespb.Service_Checker_Action) (checkerpb.Action, *ServiceActionConfig) {
				return t.Action, &ServiceActionConfig{
					Timeout:  t.Timeout.AsDuration(),
					RunCount: int(t.RunCount),
				}
			},
		),

		DefaultScore: p.DefaultScore,
		Disabled:     p.Disabled,
	}
}
