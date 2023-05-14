package grpc_servers

import (
	"context"
	"errors"
	"io"

	"go.uber.org/zap"

	"google.golang.org/protobuf/types/known/emptypb"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/model/consts"
	"ydx-goadv-gophkeeper/internal/model/enum"
	pb "ydx-goadv-gophkeeper/internal/proto"
	"ydx-goadv-gophkeeper/internal/server/services"
)

type ResourceServer struct {
	log *zap.SugaredLogger
	pb.UnimplementedResourcesServer
	service       services.ResourceService
	fileProcessor services.FileProcessor
}

func NewResourcesServer(
	service services.ResourceService,
	fileProcessor services.FileProcessor,
) pb.ResourcesServer {
	return &ResourceServer{
		log:           logger.NewLogger("res-service"),
		service:       service,
		fileProcessor: fileProcessor,
	}
}

func (s *ResourceServer) Save(ctx context.Context, resource *pb.Resource) (*pb.ResourceId, error) {
	res := &model.Resource{
		UserId: s.getUserIdFromCtx(ctx),
		Data:   resource.Data,
	}
	res.Meta = resource.Meta

	res.Type = enum.ResourceType(resource.Type)

	err := s.service.Save(ctx, res)
	return &pb.ResourceId{Id: res.Id}, err
}

func (s *ResourceServer) Delete(ctx context.Context, resId *pb.ResourceId) (*emptypb.Empty, error) {
	if err := s.service.Delete(ctx, resId.Id, s.getUserIdFromCtx(ctx)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ResourceServer) GetDescriptions(query *pb.Query, stream pb.Resources_GetDescriptionsServer) error {
	t := enum.ResourceType(query.ResourceType)
	userId := s.getUserIdFromCtx(stream.Context())
	resourceDescriptions, err := s.service.GetDescriptions(stream.Context(), userId, t)
	if err != nil {
		return err
	}

	for _, resDescription := range resourceDescriptions {
		err := stream.Send(&pb.ResourceDescription{
			Id:    resDescription.Id,
			Type:  pb.TYPE(resDescription.Type),
			Alias: resDescription.Alias,
			Meta:  resDescription.Meta,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ResourceServer) Get(ctx context.Context, id *pb.ResourceId) (*pb.Resource, error) {
	result, err := s.service.Get(ctx, id.Id, s.getUserIdFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	return &pb.Resource{
		Type: pb.TYPE(result.Type),
		Data: result.Data,
		Meta: result.Meta,
	}, nil
}

// SaveFile : TODO Save file by stream chunks
func (s *ResourceServer) SaveFile(stream pb.Resources_SaveFileServer) error {
	chunk, err := stream.Recv()
	if err == io.EOF {
		return errors.New("failed to save file: empty stream")
	}
	if err != nil {
		return err
	}

	rId, err := s.service.SaveFile(
		stream.Context(),
		s.getUserIdFromCtx(stream.Context()),
		chunk.Meta,
		chunk.Data,
	)
	if err != nil {
		return err
	}

	id := &pb.ResourceId{Id: rId}

	return stream.SendAndClose(id)
}

// SaveFile : TODO Save file by stream chunks
func (s *ResourceServer) GetFile(resId *pb.ResourceId, stream pb.Resources_GetFileServer) error {

	resource, err := s.service.Get(stream.Context(), resId.GetId(), s.getUserIdFromCtx(stream.Context()))
	if err != nil {
		return err
	}
	err = stream.Send(&pb.FileChunk{
		Meta: resource.Meta,
		Data: nil,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ResourceServer) getUserIdFromCtx(ctx context.Context) int32 {
	return ctx.Value(consts.UserIDCtxKey).(int32)
}
