package services

import (
	"context"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
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
	zap.L().Debug("ServicesService/List", zap.Any("request", req))

	// FIXME: check admin rights.

	gotVersion, err := s.controller.Versions.Get(ctx, VersionKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting version: %v", err)
	}

	requestedVersion := int(req.GetVersion().GetVersion())
	if requestedVersion > gotVersion {
		return nil, status.Errorf(codes.FailedPrecondition, "requested version is greater than current")
	}
	if requestedVersion == gotVersion {
		return &servicespb.ListResponse{Version: version.NewVersionProto(gotVersion)}, nil
	}

	services, err := s.controller.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting services: %v", err)
	}

	return &servicespb.ListResponse{
		Services: lo.Map(services, func(service *models.Service, _ int) *servicespb.Service {
			return service.ToProto()
		}),
		Version: version.NewVersionProto(gotVersion),
	}, nil
}

func (s *Service) CreateBatch(ctx context.Context, req *servicespb.CreateBatchRequest) (*servicespb.CreateBatchResponse, error) {
	zap.L().Debug("ServicesService/CreateBatch", zap.Any("request", req))

	// FIXME: check admin rights.

	if err := s.validateCreateBatch(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validating request: %v", err)
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
