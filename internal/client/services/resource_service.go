package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/model/enum"
	pb "ydx-goadv-gophkeeper/internal/proto"
)

type ResourceService interface {
	Save(ctx context.Context, resType enum.ResourceType, alias string, data []byte, meta []byte) (int32, error)
	Delete(ctx context.Context, resId int32) error
	GetDescriptions(ctx context.Context, resType enum.ResourceType) ([]*model.ResourceDescription, error)
	Get(ctx context.Context, resId int32) (*model.ResourceInfo, error)
	SaveFile(ctx context.Context, description, path string) (int32, error)
	GetFile(ctx context.Context, resId int32) (string, error)
}

type resourceService struct {
	log            *zap.SugaredLogger
	resourceClient pb.ResourcesClient
}

func NewResourceService(
	client pb.ResourcesClient,
) ResourceService {
	return &resourceService{
		log:            logger.NewLogger("res-service"),
		resourceClient: client,
	}
}

func (s *resourceService) Save(
	ctx context.Context,
	resType enum.ResourceType,
	alias string,
	data []byte,
	meta []byte,
) (int32, error) {
	resId, err := s.resourceClient.Save(ctx, &pb.Resource{
		Type:  pb.TYPE(resType),
		Alias: alias,
		Data:  data,
		Meta:  meta,
	})
	if err != nil {
		return 0, err
	}
	return resId.GetId(), nil
}

func (s *resourceService) Delete(ctx context.Context, resId int32) error {
	_, err := s.resourceClient.Delete(ctx, &pb.ResourceId{Id: resId})
	return err
}

func (s *resourceService) GetDescriptions(ctx context.Context, resType enum.ResourceType) ([]*model.ResourceDescription, error) {
	stream, err := s.resourceClient.GetDescriptions(ctx, &pb.Query{ResourceType: pb.TYPE(resType)})
	if err != nil {
		return nil, err
	}
	results := make([]*model.ResourceDescription, 0)
	for {
		descr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, &model.ResourceDescription{
			Id:    descr.Id,
			Alias: descr.Alias,
			Meta:  descr.Meta,
			Type:  enum.ResourceType(descr.Type),
		})
	}
	return results, nil
}

func (s *resourceService) Get(ctx context.Context, resId int32) (*model.ResourceInfo, error) {
	resource, err := s.resourceClient.Get(ctx, &pb.ResourceId{Id: resId})
	if err != nil {
		return nil, err
	}
	return s.parseResource(resource)
}

func (s *resourceService) parseResource(resource *pb.Resource) (*model.ResourceInfo, error) {
	switch enum.ResourceType(resource.Type) {
	case enum.LoginPassword:
		var loginPassword model.LoginPassword
		if err := json.Unmarshal(resource.Data, &loginPassword); err != nil {
			return nil, err
		}

		return &model.ResourceInfo{Resource: &loginPassword, Meta: resource.Meta}, nil

	case enum.BankCard:
		var bankCard model.BankCard
		if err := json.Unmarshal(resource.Data, &bankCard); err != nil {
			return nil, err
		}
		return &model.ResourceInfo{Resource: &bankCard, Meta: resource.Meta}, nil
	}
	return nil, fmt.Errorf("undefined type %v", resource.Type)
}

func (s *resourceService) SaveFile(ctx context.Context, description, path string) (int32, error) {
	return 0, nil
}

func (s *resourceService) GetFile(ctx context.Context, resId int32) (string, error) {
	return "", nil
}
