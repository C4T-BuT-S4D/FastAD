package api

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func (s *Service) HandleGetScoreboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		sb, err := s.boardBuilder.GetScoreboard(ctx)
		if err != nil {
			return fmt.Errorf("getting scoreboard: %w", err)
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

		sb, err := s.boardBuilder.GetScoreboard(ctx)
		if err != nil {
			return fmt.Errorf("building scoreboard state: %w", err)
		}

		teams, err := s.boardBuilder.GetTeams(ctx)
		if err != nil {
			return fmt.Errorf("getting teams: %w", err)
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
