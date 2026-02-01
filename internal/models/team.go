package models

import (
	"fmt"

	"github.com/uptrace/bun"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Team struct {
	bun.BaseModel `bun:"teams,alias:t"`

	ID        int               `bun:"id,pk,autoincrement"`
	Name      string            `bun:"name,notnull,unique"`
	Address   string            `bun:"address,notnull"`
	Token     string            `bun:"token,notnull"`
	Labels    map[string]string `bun:"labels,type:jsonb,notnull"`
	AvatarURL string            `bun:"avatar_url"`
}

func (t *Team) String() string {
	return fmt.Sprintf("Team(name=%s, ip=%s)", t.Name, t.Address)
}

func (t *Team) ToProto() *teamspb.Team {
	return &teamspb.Team{
		Id:        int64(t.ID),
		Name:      t.Name,
		Address:   t.Address,
		Token:     t.Token,
		Labels:    t.Labels,
		AvatarUrl: t.AvatarURL,
	}
}

func NewTeamFromProto(team *teamspb.Team) *Team {
	return &Team{
		ID:        int(team.GetId()),
		Name:      team.GetName(),
		Address:   team.GetAddress(),
		Token:     team.GetToken(),
		Labels:    team.GetLabels(),
		AvatarURL: team.GetAvatarUrl(),
	}
}
