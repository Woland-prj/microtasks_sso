package auth

import (
	"context"
	"errors"

	ssov1 "github.com/Woland-prj/microtasks_protos/gen/go/sso"
	"github.com/Woland-prj/microtasks_sso/internal/domain/cerrors"
	"github.com/Woland-prj/microtasks_sso/internal/domain/dtos"
	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(
		ctx context.Context,
		dto dtos.LoginDto,
	) (*entities.JwtTokenPair, error)
	Register(
		ctx context.Context,
		dto dtos.RegisterDto,
	) (int64, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	validate    *validator.Validate
	authService AuthService
}

func Register(gRPC *grpc.Server, service AuthService, validate *validator.Validate) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		authService: service,
		validate:    validate,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	r *ssov1.LoginRequest,
) (*ssov1.LoginRespones, error) {
	dto := dtos.LoginDto{
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
		AppId:    r.GetAppId(),
	}

	if err := s.validate.Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := s.authService.Login(ctx, dto)

	if err != nil {
		if errors.Is(err, &cerrors.InvalidCredentialsError{}) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginRespones{
		AuthToken:    tokens.RefreshToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	r *ssov1.RegisterRequest,
) (*ssov1.RegisterRespones, error) {
	dto := dtos.RegisterDto{
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
	}

	if err := s.validate.Struct(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	uid, err := s.authService.Register(ctx, dto)

	if err != nil {
		if errors.Is(err, &cerrors.AlreadyExistsError{}) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterRespones{
		Uid: uid,
	}, nil
}
