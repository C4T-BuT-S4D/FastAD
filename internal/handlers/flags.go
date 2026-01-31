package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

func (s *Service) HandleSubmitFlags() echo.HandlerFunc {
	type request struct {
		Flags []string `json:"flags"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req request
		if err := c.Bind(&req); err != nil {
			return httpext.NewErrorf(http.StatusBadRequest, "binding request: %v", err)
		}

		if len(req.Flags) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "flags are required")
		}

		resp, err := s.receiverClient.SubmitFlags(
			ctx,
			&receiverpb.SubmitFlagsRequest{
				Flags:     req.Flags,
				TeamToken: c.Request().Header.Get(TeamTokenHeader),
			},
		)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "submitting flags")
		}

		return ProtoJSON(c, http.StatusOK, resp)
	}
}
