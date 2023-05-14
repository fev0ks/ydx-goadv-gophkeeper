package services

import (
	"context"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/model/enum"
	"ydx-goadv-gophkeeper/internal/server/repositories"
)

type ResourceService interface {
	Save(ctx context.Context, res *model.Resource) error
	Delete(ctx context.Context, resId, userId int32) error
	GetDescriptions(ctx context.Context, userId int32, resType enum.ResourceType) ([]*model.ResourceDescription, error)
	Get(ctx context.Context, resId int32, userId int32) (*model.Resource, error)
	SaveFile(ctx context.Context, userId int32, meta []byte, data []byte) (int32, error)
	GetFile(ctx context.Context, resource *model.Resource) ([]byte, error)
}

type resourceService struct {
	log  *zap.SugaredLogger
	repo repositories.ResourceRepository
}

func NewResourceService(repo repositories.ResourceRepository) ResourceService {
	return &resourceService{log: logger.NewLogger("res-service"), repo: repo}
}

func (s *resourceService) Save(ctx context.Context, data *model.Resource) error {
	return s.repo.Save(ctx, data)
}

func (s *resourceService) Delete(ctx context.Context, resId int32, userId int32) error {
	return s.repo.Delete(ctx, resId, userId)
}

func (s *resourceService) GetDescriptions(ctx context.Context, userId int32, resType enum.ResourceType) ([]*model.ResourceDescription, error) {
	return s.repo.GetResDescriptionsByType(ctx, userId, resType)
}

func (s *resourceService) Get(ctx context.Context, resId int32, userId int32) (*model.Resource, error) {
	return s.repo.Get(ctx, resId, userId)
}

func (s *resourceService) SaveFile(ctx context.Context, userId int32, meta []byte, data []byte) (int32, error) {
	resource := &model.Resource{
		UserId: userId,
		Data:   data,
	}
	resource.Type = enum.File
	resource.Meta = meta

	err := s.repo.Save(ctx, resource)
	if err != nil {
		return 0, err
	}

	return resource.Id, nil
}

func (s *resourceService) GetFile(ctx context.Context, resource *model.Resource) ([]byte, error) {
	res, err := s.repo.Get(ctx, resource.Id, resource.UserId)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}