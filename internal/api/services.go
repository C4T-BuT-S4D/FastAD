package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

func (s *Service) HandleServicesList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		services, err := s.servicesClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing services")
		}

		for _, service := range services {
			service.Actions = nil
			service.CheckerPath = ""
			service.DefaultScore = 0
			service.CheckerType = checkerpb.Type_TYPE_UNSPECIFIED
			service.DefaultTimeout = 0
		}

		return c.JSON(http.StatusOK, services)
	}
}
