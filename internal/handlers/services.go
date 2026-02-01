package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

func (s *Service) HandleServicesList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		services, err := s.servicesClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing services")
		}

		resp := &servicespb.Service_Batch{
			Services: lo.Map(services, func(service *servicespb.Service, _ int) *servicespb.Service {
				serviceCloned := service.CloneVT()
				serviceCloned.Checker = nil
				return serviceCloned
			}),
		}

		return ProtoJSON(c, http.StatusOK, resp)
	}
}
