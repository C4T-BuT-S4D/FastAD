package api

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

type teamServiceKey struct {
	TeamID    int
	ServiceID int
}

type teamServiceState struct {
	TeamID       int     `json:"team_id"`
	ServiceID    int     `json:"service_id"`
	ChecksTotal  int     `json:"checks_total"`
	ChecksPassed int     `json:"checks_passed"`
	StolenFlags  int     `json:"stolen_flags"`
	LostFlags    int     `json:"lost_flags"`
	Points       float64 `json:"points"`
}

type scoreboardState struct {
	TeamServiceStates []*teamServiceState `json:"team_service_states"`
}

func (s *Service) HandleGetScoreboard() echo.HandlerFunc {
	// TODO: cache this for as long as possible.
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		teams, err := s.teamsClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing teams")
		}

		sb, err := s.buildScoreboardState(ctx, teams)
		if err != nil {
			return fmt.Errorf("building scoreboard state: %w", err)
		}

		return c.JSON(http.StatusOK, sb)
	}
}

func (s *Service) HandleGetCTFTimeScoreboard() echo.HandlerFunc {
	type ctftimeTeamState struct {
		Pos   int     `json:"pos"`
		Team  string  `json:"team"`
		Score float64 `json:"score"`
	}

	type response struct {
		Standings []*ctftimeTeamState `json:"standings"`
	}

	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		teams, err := s.teamsClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing teams")
		}

		sb, err := s.buildScoreboardState(ctx, teams)
		if err != nil {
			return fmt.Errorf("building scoreboard state: %w", err)
		}

		teamStates := make(map[int]*ctftimeTeamState)
		for _, team := range teams {
			teamStates[team.ID] = &ctftimeTeamState{
				Team:  team.Name,
				Score: 0,
			}
		}
		for _, tss := range sb.TeamServiceStates {
			sla := 0.0
			if tss.ChecksTotal > 0 {
				sla = float64(tss.ChecksPassed) / float64(tss.ChecksTotal)
			}
			teamStates[tss.TeamID].Score += tss.Points * sla
		}

		teamStatesList := lo.Filter(lo.Values(teamStates), func(item *ctftimeTeamState, _ int) bool {
			return item.Score > 0
		})
		slices.SortFunc(teamStatesList, func(t1, t2 *ctftimeTeamState) int {
			return cmp.Compare(t2.Score, t1.Score)
		})

		for i, teamState := range teamStatesList {
			teamState.Pos = i + 1
		}

		return c.JSON(http.StatusOK, response{Standings: teamStatesList})
	}
}

func (s *Service) buildScoreboardState(ctx context.Context, teams []*models.Team) (*scoreboardState, error) {
	services, err := s.servicesClient.List(ctx)
	if err != nil {
		return nil, httpext.NewErrorFromStatus(err, "listing services")
	}

	sbMap := make(map[teamServiceKey]*teamServiceState)
	for _, team := range teams {
		for _, service := range services {
			sbMap[teamServiceKey{TeamID: team.ID, ServiceID: service.ID}] = &teamServiceState{
				TeamID:    team.ID,
				ServiceID: service.ID,
				Points:    service.DefaultScore,
			}
		}
	}

	slaState, err := s.scoreboardClient.GetState(ctx, &scoreboard.GetStateRequest{})
	if err != nil {
		return nil, httpext.NewErrorFromStatus(err, "getting scoreboard state")
	}

	receiverState, err := s.receiverClient.GetState(ctx, &receiverpb.GetStateRequest{})
	if err != nil {
		return nil, httpext.NewErrorFromStatus(err, "getting receiver state")
	}

	for _, tss := range slaState.GetScoreboard().GetTeamServiceStates() {
		key := teamServiceKey{TeamID: int(tss.GetTeamId()), ServiceID: int(tss.GetServiceId())}
		if sbs, ok := sbMap[key]; ok {
			sbs.ChecksTotal = int(tss.GetChecksTotal())
			sbs.ChecksPassed = int(tss.GetChecksPassed())
		}
	}

	for _, tss := range receiverState.GetState().GetTeamServices() {
		key := teamServiceKey{TeamID: int(tss.GetTeamId()), ServiceID: int(tss.GetServiceId())}
		if sbs, ok := sbMap[key]; ok {
			sbs.StolenFlags = int(tss.GetStolenFlags())
			sbs.LostFlags = int(tss.GetLostFlags())
			sbs.Points = tss.GetPoints()
		}
	}

	return &scoreboardState{TeamServiceStates: lo.Values(sbMap)}, nil
}
