package authgrpc

import (
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/k1v4/protos/gen/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/internal/services/auth"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email, password string, appId int32) (string, error)
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(grpc *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpc, &serverApi{
		auth: auth,
	})
}

func (s *serverApi) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	fmt.Println(userId, err)
	if err != nil {
		if errors.Is(err, auth.ErrUserExist) {
			return nil, status.Error(codes.InvalidArgument, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{UserId: userId}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
