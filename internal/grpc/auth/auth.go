package auth

import (
	"context"
	ssov1 "github.com/Woland-prj/microtasks_protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	r *ssov1.LoginRequest,
) (*ssov1.LoginRespones, error) {
	panic("implement me")
}

func (s *serverAPI) Register(
	ctx context.Context,
	r *ssov1.RegisterRequest,
) (*ssov1.RegisterRespones, error) {
	return &ssov1.RegisterRespones{
		Uid: 12345,
	}, nil
}
