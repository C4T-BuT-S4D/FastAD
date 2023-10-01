package setup

import (
	"fmt"
	"time"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Game struct {
	StartTime   time.Time  `yaml:"start_time"`
	EndTime     *time.Time `yaml:"end_time"`
	TotalRounds int        `yaml:"total_rounds"`

	FlagLifetimeRounds int           `yaml:"flag_lifetime_rounds"`
	RoundDuration      time.Duration `yaml:"round_duration"`
}

func (g *Game) Validate() error {
	if g.StartTime.IsZero() {
		return fmt.Errorf("start time required")
	}
	if g.EndTime != nil && g.EndTime.Before(g.StartTime) {
		return fmt.Errorf("end time is before start time")
	}
	if g.TotalRounds < 0 {
		return fmt.Errorf("total rounds must be non-negative")
	}
	return nil
}

type Team struct {
	Name    string            `yaml:"name"`
	Address string            `yaml:"address"`
	Labels  map[string]string `yaml:"labels"`
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name required")
	}
	if t.Address == "" {
		return fmt.Errorf("address required")
	}
	return nil
}

func (t *Team) ToProto() *teamspb.Team {
	return &teamspb.Team{
		Name:    t.Name,
		Address: t.Address,
		Labels:  t.Labels,
	}
}

type CheckerAction struct {
	Count   int           `yaml:"count"`
	Timeout time.Duration `yaml:"timeout"`
}

type Checker struct {
	Type checkerpb.Type `yaml:"type"`
	Path string         `yaml:"path"`

	DefaultTimeout time.Duration                      `yaml:"default_timeout"`
	Actions        map[checkerpb.Action]CheckerAction `yaml:"actions"`
}

func (c *Checker) Validate() error {
	if c.Type == checkerpb.Type_TYPE_UNSPECIFIED {
		c.Type = checkerpb.Type_TYPE_LEGACY
	}
	if c.Path == "" {
		return fmt.Errorf("path required")
	}
	return nil
}

type Service struct {
	Name         string  `yaml:"name"`
	DefaultScore float64 `yaml:"default_score"`

	Checker *Checker `yaml:"checker"`
}

func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name required")
	}
	if s.Checker == nil {
		s.Checker = &Checker{}
	}
	if err := s.Checker.Validate(); err != nil {
		return fmt.Errorf("checker: %w", err)
	}
	return nil
}

func (s *Service) ToProto() *servicespb.Service {
	return &servicespb.Service{
		Name:         s.Name,
		DefaultScore: s.DefaultScore,
		Checker: &servicespb.Service_Checker{
			Type:           s.Checker.Type,
			Path:           s.Checker.Path,
			DefaultTimeout: durationpb.New(s.Checker.DefaultTimeout),
			Actions: lo.MapToSlice(s.Checker.Actions, func(action checkerpb.Action, actionConfig CheckerAction) *servicespb.Service_Checker_Action {
				return &servicespb.Service_Checker_Action{
					Action:   action,
					RunCount: int32(actionConfig.Count),
					Timeout:  durationpb.New(actionConfig.Timeout),
				}
			}),
		},
	}
}

type GameConfig struct {
	Game Game `yaml:"game"`

	Teams    []Team    `yaml:"teams"`
	Services []Service `yaml:"services"`
}

func (c *GameConfig) Validate() error {
	if err := c.Game.Validate(); err != nil {
		return fmt.Errorf("game: %w", err)
	}
	for i, team := range c.Teams {
		if err := team.Validate(); err != nil {
			return fmt.Errorf("team %d: %w", i, err)
		}
	}
	for i, service := range c.Services {
		if err := service.Validate(); err != nil {
			return fmt.Errorf("service %d: %w", i, err)
		}
	}
	return nil
}
