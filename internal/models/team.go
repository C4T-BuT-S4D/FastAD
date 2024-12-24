package models

import (
	"fmt"

	"github.com/uptrace/bun"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Team struct {
	bun.BaseModel `bun:"teams,alias:t"`

	ID      int               `bun:"id,pk,autoincrement" json:"id"`
	Name    string            `bun:"name,notnull,unique" json:"name"`
	Address string            `bun:"address,notnull" json:"address"`
	Token   string            `bun:"token,notnull" json:"token,omitempty"`
	Labels  map[string]string `bun:"labels,type:jsonb,notnull" json:"labels"`
}

func (t *Team) String() string {
	return fmt.Sprintf("Team(name=%s, ip=%s)", t.Name, t.Address)
}

func (t *Team) ToProto() *teamspb.Team {
	return &teamspb.Team{
		Id:      int64(t.ID),
		Name:    t.Name,
		Address: t.Address,
		Token:   t.Token,
		Labels:  t.Labels,
	}
}

func NewTeamFromProto(team *teamspb.Team) *Team {
	return &Team{
		ID:      int(team.Id),
		Name:    team.Name,
		Address: team.Address,
		Token:   team.Token,
		Labels:  team.Labels,
	}
}
