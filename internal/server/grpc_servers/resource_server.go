package grpc_servers

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "ydx-goadv-gophkeeper/api/proto"
	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model/consts"
	"ydx-goadv-gophkeeper/internal/model/enum"
	"ydx-goadv-gophkeeper/internal/model/resources"
	"ydx-goadv-gophkeeper/internal/server/services"
	intsrv "ydx-goadv-gophkeeper/internal/services"
)

type ResourceServer struct {
	log *zap.SugaredLogger
	pb.UnimplementedResourcesServer
	service     services.ResourceService
	fileService intsrv.FileService
}

func NewResourcesServer(
	service services.ResourceService,
	fileService intsrv.FileService,
) pb.ResourcesServer {
	return &ResourceServer{
		log:         logger.NewLogger("res-service"),
		service:     service,
		fileService: fileService,
	}
}

func (s *ResourceServer) Save(ctx context.Context, resource *pb.Resource) (*pb.ResourceId, error) {
	res := &resources.Resource{
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
			Id:   resDescription.Id,
			Type: pb.TYPE(resDescription.Type),
			Meta: resDescription.Meta,
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
	chunks := make(chan []byte)

	userId := s.getUserIdFromCtx(stream.Context())
	resId, err := s.service.SaveFileDescription(
		stream.Context(),
		userId,
		chunk.Meta,
		chunk.Data,
	)
	if err != nil {
		return err
	}
	errCh, err := s.fileService.SaveFile(fmt.Sprintf("./cmd/server/%d", resId), chunks)
	if err != nil {
		return err
	}
Loop:
	for {
		chunk, err = stream.Recv()
		if err == io.EOF {
			close(chunks)
			break Loop
		}
		if err != nil {
			close(chunks)
			return fmt.Errorf("failed to save file: %v", err)
		}
		select {
		case chunks <- chunk.Data:
		case <-errCh:
			close(chunks)
			break Loop
		}
	}

	id := &pb.ResourceId{Id: resId}

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
		Data: resource.Data,
	})
	if err != nil {
		return err
	}
	errCh := make(chan error)
	chunks, _, err := s.fileService.ReadFile(fmt.Sprintf("./cmd/server/%d", resource.Id), errCh)
	if err != nil {
		return err
	}

Loop:
	for {
		chunk, ok := <-chunks
		if !ok {
			break Loop
		}
		err := stream.Send(&pb.FileChunk{
			Meta: nil,
			Data: chunk,
		})
		if err != nil {
			errCh <- err
			return err
		}
	}
	return nil
}

func (s *ResourceServer) getUserIdFromCtx(ctx context.Context) int32 {
	return ctx.Value(consts.UserIDCtxKey).(int32)
}
