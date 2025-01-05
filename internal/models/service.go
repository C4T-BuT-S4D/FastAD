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

type ServiceActionConfig struct {
	Timeout  time.Duration `json:"timeout"`
	RunCount int           `json:"run_count"`
}

type Service struct {
	bun.BaseModel `bun:"services,alias:s"`

	ID   int    `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull,unique"`

	CheckerType checkerpb.Type `bun:"checker_type,notnull"`
	CheckerPath string         `bun:"checker_path,notnull"`

	DefaultScore float64 `bun:"default_score,notnull"`

	DefaultTimeout time.Duration                             `bun:"default_timeout,notnull"`
	Actions        map[checkerpb.Action]*ServiceActionConfig `bun:"actions,type:jsonb,notnull"`

	Disabled bool `bun:"disabled,notnull"`

	// TODO: vulns format.
	// Places int
}

func (s *Service) String() string {
	return fmt.Sprintf("Service(name=%s)", s.Name)
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
		ID:   int(p.GetId()),
		Name: p.GetName(),

		CheckerType:    p.GetChecker().GetType(),
		CheckerPath:    p.GetChecker().GetPath(),
		DefaultTimeout: p.GetChecker().GetDefaultTimeout().AsDuration(),
		Actions: lo.SliceToMap(
			p.GetChecker().GetActions(),
			func(t *servicespb.Service_Checker_Action) (checkerpb.Action, *ServiceActionConfig) {
				return t.GetAction(), &ServiceActionConfig{
					Timeout:  t.GetTimeout().AsDuration(),
					RunCount: int(t.GetRunCount()),
				}
			},
		),

		DefaultScore: p.GetDefaultScore(),
		Disabled:     p.GetDisabled(),
	}
}
