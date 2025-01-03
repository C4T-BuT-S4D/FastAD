package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func (s *Service) HandleGetGameState() echo.HandlerFunc {
	type request struct {
		Format string `query:"format"`
	}

	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		var req request
		if err := c.Bind(&req); err != nil {
			return httpext.NewErrorf(http.StatusBadRequest, "parsing request: %v", err)
		}

		gs, err := s.gameStateClient.Get(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "getting game state")
		}

		switch req.Format {
		case "proto":
			return ProtoRaw(c, http.StatusOK, gs)
		default:
			return ProtoJSON(c, http.StatusOK, gs)
		}
	}
}

func (s *Service) HandleGetAttackData() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		var lastSnapshot models.AttackDataSnapshot
		if err := s.db.NewSelect().
			Model(&lastSnapshot).
			OrderExpr("created_at DESC").
			Scan(ctx); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return httpext.NewErrorf(http.StatusNotFound, "no attack data snapshots")
			}
			return httpext.NewErrorFromStatus(err, "getting last attack data snapshot")
		}

		return c.JSON(http.StatusOK, lastSnapshot.Payload)
	}
}
