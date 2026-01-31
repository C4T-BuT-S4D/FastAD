package handlers

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (s *Service) HandleGetScoreboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		bs, err := s.boardBuilder.GetState()
		if err != nil {
			return fmt.Errorf("getting scoreboard state: %w", err)
		}

		return ProtoJSON(c, http.StatusOK, bs.Scoreboard)
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
		bs, err := s.boardBuilder.GetState()
		if err != nil {
			return fmt.Errorf("building scoreboard state: %w", err)
		}

		teamStates := make(map[int64]*ctftimeTeamState)
		for _, team := range bs.Teams {
			teamStates[team.GetId()] = &ctftimeTeamState{
				Team:  team.GetName(),
				Score: 0,
			}
		}
		for _, tss := range bs.Scoreboard.GetTeamServiceStates() {
			sla := 0.0
			if total := tss.GetChecksTotal(); total > 0 {
				sla = float64(tss.GetChecksPassed()) / float64(total)
			}
			teamStates[tss.GetTeamId()].Score += tss.GetPoints() * sla
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
