package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

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
				serviceCloned := proto.Clone(service).(*servicespb.Service)
				serviceCloned.Checker = nil
				serviceCloned.DefaultScore = 0
				return serviceCloned
			}),
		}

		// TODO: helper for returning protojson.
		raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(resp)
		if err != nil {
			return httpext.NewErrorf(
				http.StatusInternalServerError,
				"marshaling services: %v",
				err,
			)
		}

		return c.JSONBlob(http.StatusOK, raw)
	}
}
