package clients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"ydx-goadv-gophkeeper/internal/client/interceptors"
	"ydx-goadv-gophkeeper/internal/client/model"
)

func CreateGrpcConnection(port string, tokenHolder *model.TokenHolder) (*grpc.ClientConn, error) {
	tokenProcessor := interceptors.NewRequestTokenProcessor(tokenHolder)
	return grpc.Dial(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(tokenProcessor.TokenInterceptor()),
		grpc.WithStreamInterceptor(tokenProcessor.TokenStreamInterceptor()),
	)
}
