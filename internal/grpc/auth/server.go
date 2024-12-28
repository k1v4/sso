package auth

import (
	"context"
	ssov1 "github.com/k1v4/protos/gen/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
}

func Register(grpc *grpc.Server) {
	ssov1.RegisterAuthServer(grpc, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverApi) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverApi) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
