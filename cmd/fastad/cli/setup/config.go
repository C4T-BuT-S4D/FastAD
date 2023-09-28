package setup

import (
	"fmt"
	"time"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Game struct {
	StartTime time.Time `yaml:"start_time"`
	EndTime   time.Time `yaml:"end_time"`
}

func (g *Game) Validate() error {
	if g.StartTime.IsZero() {
		return fmt.Errorf("start time is not set")
	}
	if g.EndTime.Before(g.StartTime) {
		return fmt.Errorf("end time is before start time")
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
		return fmt.Errorf("name is not set")
	}
	if t.Address == "" {
		return fmt.Errorf("address is not set")
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

type Service struct {
	Name string `yaml:"name"`
}

func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is not set")
	}
	return nil
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
