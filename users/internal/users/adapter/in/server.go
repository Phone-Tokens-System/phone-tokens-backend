package in

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"

	httpadapter "users/internal/users/adapter/in/http"
	"users/internal/users/model"
	"users/internal/users/service/users"
	userspb "users/pkg/proto"
)

type Server struct {
	userspb.UnimplementedUserServiceServer
	svc users.Service
}

func NewServer(svc users.Service) *Server {
	return &Server{svc: svc}
}

// NewHTTPServer создаёт HTTP‑сервер поверх доменного сервиса.
func NewHTTPServer(httpPort string, jwtSecret string, svc users.Service) (*http.Server, error) {
	if httpPort == "" {
		return nil, errors.New("HTTP_PORT is required")
	}

	handler := httpadapter.NewHandler(svc)

	mux := http.NewServeMux()
	httpadapter.RegisterRoutes(mux, handler, httpadapter.AuthConfig{
		JWTSecret: jwtSecret,
	})

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	return server, nil
}

// Start launches a gRPC server on the given address (e.g. ":9090").
// It uses JSON over gRPC for simplicity and to avoid generated code.
func Start(addr string, svc users.Service) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("listen: %w", err)
	}

	encoding.RegisterCodec(userspb.JSONCodec{})
	server := grpc.NewServer(grpc.ForceServerCodec(userspb.JSONCodec{}))

	userspb.RegisterUserServiceServer(server, NewServer(svc))

	return server, lis, nil
}

func (s *Server) Register(ctx context.Context, req *userspb.RegisterRequest) (*userspb.RegisterResponse, error) {
	role := model.RoleUser
	switch req.Role {
	case userspb.Role_ROLE_ADMIN:
		role = model.RoleAdmin
	case userspb.Role_ROLE_USER:
		role = model.RoleUser
	}

	user, err := s.svc.Register(ctx, req.Phone, req.Password, role)
	if err != nil {
		switch err {
		case users.ErrPhoneAlreadyUsed:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userspb.RegisterResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *Server) Login(ctx context.Context, req *userspb.LoginRequest) (*userspb.LoginResponse, error) {
	token, user, err := s.svc.Authenticate(ctx, req.Phone, req.Password)
	if err != nil {
		switch err {
		case users.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &userspb.LoginResponse{
		Token: token,
		User:  toProtoUser(user),
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.GetUserResponse, error) {
	user, err := s.svc.GetByID(ctx, req.Id)
	if err != nil {
		if err == users.ErrNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userspb.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

func toProtoUser(u *model.User) *userspb.User {
	if u == nil {
		return nil
	}

	return &userspb.User{
		Id:        u.ID,
		Phone:     u.Phone,
		Role:      fromModelRole(u.Role),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func fromModelRole(r model.Role) userspb.Role {
	switch r {
	case model.RoleAdmin:
		return userspb.Role_ROLE_ADMIN
	case model.RoleUser:
		return userspb.Role_ROLE_USER
	default:
		return userspb.Role_ROLE_UNSPECIFIED
	}
}
