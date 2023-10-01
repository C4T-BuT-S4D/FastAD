package setup

import (
	"fmt"
	"time"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Game struct {
	StartTime   time.Time  `yaml:"start_time"`
	EndTime     *time.Time `yaml:"end_time"`
	TotalRounds uint       `yaml:"total_rounds"`

	FlagLifetimeRounds uint          `yaml:"flag_lifetime_rounds"`
	RoundDuration      time.Duration `yaml:"round_duration"`

	Mode GameMode `yaml:"mode"`
}

func (g *Game) Validate() error {
	if g.StartTime.IsZero() {
		return fmt.Errorf("start_time required")
	}
	if g.EndTime != nil && g.EndTime.Before(g.StartTime) {
		return fmt.Errorf("end_time is before start_time")
	}

	if g.Mode == "" {
		g.Mode = GameModeClassic
	}
	if err := g.Mode.Validate(); err != nil {
		return fmt.Errorf("mode: %w", err)
	}

	return nil
}

func (g *Game) ToUpdateRequestProto() *gspb.UpdateRequest {
	res := &gspb.UpdateRequest{
		StartTime:   timestamppb.New(g.StartTime),
		TotalRounds: uint32(g.TotalRounds),

		FlagLifetimeRounds: uint32(g.FlagLifetimeRounds),
		RoundDuration:      durationpb.New(g.RoundDuration),

		Mode: g.Mode.ToProto(),
	}

	if g.EndTime != nil {
		res.EndTime = timestamppb.New(*g.EndTime)
	}

	return res
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

type CheckerActionConfig struct {
	Count   int           `yaml:"count"`
	Timeout time.Duration `yaml:"timeout"`
}

type Checker struct {
	Type CheckerType `yaml:"type"`
	Path string      `yaml:"path"`

	DefaultTimeout time.Duration                         `yaml:"default_timeout"`
	Actions        map[CheckerAction]CheckerActionConfig `yaml:"actions"`
}

func (c *Checker) Validate() error {
	if c.Type == "" {
		c.Type = CheckerTypeLegacy
	}

	if err := c.Type.Validate(); err != nil {
		return fmt.Errorf("type: %w", err)
	}

	if c.Path == "" {
		return fmt.Errorf("path required")
	}

	if c.DefaultTimeout == 0 {
		return fmt.Errorf("default_timeout required")
	}

	for action := range c.Actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("action %s: %w", action, err)
		}
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
	if s.DefaultScore == 0 {
		return fmt.Errorf("default_score required")
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
			Type:           s.Checker.Type.ToProto(),
			Path:           s.Checker.Path,
			DefaultTimeout: durationpb.New(s.Checker.DefaultTimeout),
			Actions: lo.MapToSlice(
				s.Checker.Actions,
				func(action CheckerAction, actionConfig CheckerActionConfig) *servicespb.Service_Checker_Action {
					return &servicespb.Service_Checker_Action{
						Action:   action.ToProto(),
						RunCount: int32(actionConfig.Count),
						Timeout:  durationpb.New(actionConfig.Timeout),
					}
				},
			),
		},
	}
}

type GameConfig struct {
	Game *Game `yaml:"game"`

	Teams    []*Team    `yaml:"teams"`
	Services []*Service `yaml:"services"`
}

func (c *GameConfig) Validate() error {
	if err := c.Game.Validate(); err != nil {
		return fmt.Errorf("game: %w", err)
	}

	for i, team := range c.Teams {
		if team == nil {
			return fmt.Errorf("team %d: nil", i)
		}
		if err := team.Validate(); err != nil {
			return fmt.Errorf("team %d: %w", i, err)
		}
	}

	for i, service := range c.Services {
		if service == nil {
			return fmt.Errorf("service %d: nil", i)
		}
		if err := service.Validate(); err != nil {
			return fmt.Errorf("service %d: %w", i, err)
		}
	}
	return nil
}
