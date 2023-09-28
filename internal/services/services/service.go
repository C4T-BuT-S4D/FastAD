package services

import (
	"context"

	"github.com/c4t-but-s4d/fastad/internal/models"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	servicespb.UnimplementedServicesServiceServer

	controller *Controller
}

func NewService(controller *Controller) *Service {
	return &Service{
		controller: controller,
	}
}

func (s *Service) List(ctx context.Context, req *servicespb.ListRequest) (*servicespb.ListResponse, error) {
	logrus.Debugf("TeamsService/List: %v", req)

	// FIXME: check admin rights.

	lastUpdate, err := s.controller.LastUpdate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting last update time: %v", err)
	}

	if req.LastUpdate != 0 {
		if req.LastUpdate > lastUpdate {
			return nil, status.Errorf(codes.InvalidArgument, "last update time is in the future")
		}
		if req.LastUpdate == lastUpdate {
			return &servicespb.ListResponse{LastUpdate: lastUpdate}, nil
		}
	}

	services, err := s.controller.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting services: %v", err)
	}

	return &servicespb.ListResponse{
		LastUpdate: lastUpdate,
		Services: lo.Map(services, func(service *models.Service, _ int) *servicespb.Service {
			return service.ToProto()
		}),
	}, nil
}

func (s *Service) CreateBatch(ctx context.Context, req *servicespb.CreateBatchRequest) (*servicespb.CreateBatchResponse, error) {
	logrus.Debugf("TeamsService/CreateBatch: %v", req)

	// FIXME: check admin rights.

	if err := s.sanitizeCreateBatch(req); err != nil {
		return nil, err
	}

	services := lo.Map(req.Services, func(service *servicespb.Service, _ int) *models.Service {
		serviceModel := models.NewServiceFromProto(service)
		serviceModel.ID = 0
		return serviceModel
	})

	if err := s.controller.CreateBatch(ctx, services); err != nil {
		return nil, status.Errorf(codes.Internal, "creating services: %v", err)
	}

	result := lo.Map(services, func(service *models.Service, _ int) *servicespb.Service {
		return service.ToProto()
	})
	return &servicespb.CreateBatchResponse{
		Services: result,
	}, nil
}
