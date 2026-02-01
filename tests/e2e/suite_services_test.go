//go:build e2e

package e2e

import (
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

type ServicesSuite struct {
	BaseSuite
}

func (s *ServicesSuite) TestListServices() {
	resp, err := http.Get(baseURL + "/api/services")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Assert().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var servicesResp servicespb.Service_Batch
	err = protojson.Unmarshal(body, &servicesResp)
	s.Require().NoError(err)

	services := servicesResp.GetServices()
	s.Assert().Len(services, 1, "Expected 1 service")

	service := services[0]
	s.Assert().Equal("test_service", service.GetName())
	s.Assert().Greater(service.GetId(), int64(0), "Service should have positive ID")
	s.Assert().Nil(service.GetChecker(), "Checker config should be stripped from API response")
	s.Assert().Equal(float64(2500), service.GetDefaultScore(), "Service should have default score of 2500")
}
