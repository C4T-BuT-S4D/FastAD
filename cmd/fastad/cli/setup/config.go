package setup

import (
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Game struct {
	StartTime   time.Time  `yaml:"start_time"`
	EndTime     *time.Time `yaml:"end_time"`
	TotalRounds uint32     `yaml:"total_rounds"`

	FlagLifetimeRounds uint32        `yaml:"flag_lifetime_rounds"`
	RoundDuration      time.Duration `yaml:"round_duration"`
	Hardness           float64       `yaml:"hardness"`
	Inflation          bool          `yaml:"inflation"`

	Mode GameMode `yaml:"mode"`
}

func (g *Game) Validate() error {
	if g.StartTime.IsZero() {
		return errors.New("start_time required")
	}
	if g.EndTime != nil && g.EndTime.Before(g.StartTime) {
		return errors.New("end_time is before start_time")
	}
	if g.Hardness <= 0 {
		return errors.New("hardness must be positive")
	}
	return nil
}

func (g *Game) ToUpdateRequestProto() *gspb.UpdateRequest {
	res := &gspb.UpdateRequest{
		StartTime:   timestamppb.New(g.StartTime),
		TotalRounds: uint64(g.TotalRounds),

		FlagLifetimeRounds: uint64(g.FlagLifetimeRounds),
		RoundDuration:      durationpb.New(g.RoundDuration),
		Hardness:           g.Hardness,
		Inflation:          g.Inflation,

		Mode: gspb.GameMode(g.Mode),
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
		return errors.New("name required")
	}
	if t.Address == "" {
		return errors.New("address required")
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
	if c.Type == CheckerType(checkerpb.Type_TYPE_UNSPECIFIED) {
		c.Type = CheckerType(checkerpb.Type_TYPE_LEGACY)
	}
	if c.Path == "" {
		return errors.New("path required")
	}
	if c.DefaultTimeout == 0 {
		return errors.New("default_timeout required")
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
		return errors.New("name required")
	}
	if s.DefaultScore == 0 {
		return errors.New("default_score required")
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
			Type:           checkerpb.Type(s.Checker.Type),
			Path:           s.Checker.Path,
			DefaultTimeout: durationpb.New(s.Checker.DefaultTimeout),
			Actions: lo.MapToSlice(
				s.Checker.Actions,
				func(action CheckerAction, actionConfig CheckerActionConfig) *servicespb.Service_Checker_Action {
					return &servicespb.Service_Checker_Action{
						Action:   checkerpb.Action(action),
						RunCount: int64(actionConfig.Count),
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
