package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model/enum"
	"ydx-goadv-gophkeeper/internal/model/resources"
	pb "ydx-goadv-gophkeeper/internal/proto"
	"ydx-goadv-gophkeeper/internal/services"
)

type ResourceService interface {
	Save(ctx context.Context, resType enum.ResourceType, data []byte, meta []byte) (int32, error)
	Delete(ctx context.Context, resId int32) error
	GetDescriptions(ctx context.Context, resType enum.ResourceType) ([]*resources.ResourceDescription, error)
	Get(ctx context.Context, resId int32) (*resources.ResourceInfo, error)
	SaveFile(ctx context.Context, path string, meta []byte) (int32, error)
	GetFile(ctx context.Context, resId int32) (string, error)
}

type resourceService struct {
	log            *zap.SugaredLogger
	resourceClient pb.ResourcesClient
	fileService    services.FileService
}

func NewResourceService(
	client pb.ResourcesClient,
	fileService services.FileService,
) ResourceService {
	return &resourceService{
		log:            logger.NewLogger("res-service"),
		resourceClient: client,
		fileService:    fileService,
	}
}

func (s *resourceService) Save(
	ctx context.Context,
	resType enum.ResourceType,
	data []byte,
	meta []byte,
) (int32, error) {
	resId, err := s.resourceClient.Save(ctx, &pb.Resource{
		Type: pb.TYPE(resType),
		Data: data,
		Meta: meta,
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

func (s *resourceService) GetDescriptions(ctx context.Context, resType enum.ResourceType) ([]*resources.ResourceDescription, error) {
	stream, err := s.resourceClient.GetDescriptions(ctx, &pb.Query{ResourceType: pb.TYPE(resType)})
	if err != nil {
		return nil, err
	}
	results := make([]*resources.ResourceDescription, 0)
	for {
		descr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, &resources.ResourceDescription{
			Id:   descr.Id,
			Meta: descr.Meta,
			Type: enum.ResourceType(descr.Type),
		})
	}
	return results, nil
}

func (s *resourceService) Get(ctx context.Context, resId int32) (*resources.ResourceInfo, error) {
	resource, err := s.resourceClient.Get(ctx, &pb.ResourceId{Id: resId})
	if err != nil {
		return nil, err
	}
	return s.parseResource(resource)
}

func (s *resourceService) parseResource(resource *pb.Resource) (*resources.ResourceInfo, error) {
	switch enum.ResourceType(resource.Type) {
	case enum.LoginPassword:
		var loginPassword resources.LoginPassword
		if err := json.Unmarshal(resource.Data, &loginPassword); err != nil {
			return nil, err
		}

		return &resources.ResourceInfo{Resource: &loginPassword, Meta: resource.Meta}, nil

	case enum.BankCard:
		var bankCard resources.BankCard
		if err := json.Unmarshal(resource.Data, &bankCard); err != nil {
			return nil, err
		}
		return &resources.ResourceInfo{Resource: &bankCard, Meta: resource.Meta}, nil
	}
	return nil, fmt.Errorf("undefined type %v", resource.Type)
}

func (s *resourceService) SaveFile(ctx context.Context, path string, meta []byte) (int32, error) {
	stream, err := s.resourceClient.SaveFile(ctx)
	if err != nil {
		return 0, err
	}
	errCh := make(chan error)
	chunks, stat, err := s.fileService.ReadFile(path, errCh)
	if err != nil {
		return 0, err
	}
	fileDescriptionJson, err := json.Marshal(resources.File{
		Name:      stat.Name(),
		Extension: filepath.Ext(path),
		Size:      stat.Size(),
	})
	err = stream.Send(&pb.FileChunk{
		Meta: meta,
		Data: fileDescriptionJson,
	})
	if err != nil {
		return 0, err
	}

	for {
		chunk, ok := <-chunks
		if !ok {
			break
		}
		err := stream.Send(&pb.FileChunk{
			Meta: nil,
			Data: chunk,
		})
		if err != nil {
			errCh <- err
			return 0, err
		}
	}
	resId, err := stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}
	return resId.Id, nil
}

func (s *resourceService) GetFile(ctx context.Context, resId int32) (string, error) {
	stream, err := s.resourceClient.GetFile(ctx, &pb.ResourceId{Id: resId})
	if err != nil {
		return "", err
	}
	chunk, err := stream.Recv()
	if err != nil {
		return "", err
	}
	var fileDescription resources.File
	err = json.Unmarshal(chunk.Data, &fileDescription)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("./%s", fileDescription.Name)
	chunks := make(chan []byte)
	errCh, err := s.fileService.SaveFile(path, chunks)
	if err != nil {
		return "", err
	}
Loop:
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			close(chunks)
			break Loop
		}
		if err != nil {
			close(chunks)
			s.log.Errorf("failed to recieve file stream chunk: %v", err)
			return "", err
		}

		select {
		case chunks <- chunk.Data:
		case _ = <-errCh:
			close(chunks)
			break Loop
		}
	}
	return path, err
}
