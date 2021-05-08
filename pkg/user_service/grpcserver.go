package user_service

import (
	"context"
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	api.UnimplementedUserServiceServer
}

func (s *GRPCServer) CreateUser(context context.Context, request *api.CreateUserRequest) (*api.User, error) {
	fmt.Println(request)
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (s *GRPCServer) UpdateUser(context context.Context, request *api.UpdateUserRequest) (*api.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (s *GRPCServer) DeleteUser(context context.Context, request *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (s *GRPCServer) ListUser(context context.Context, request *api.ListUserRequest) (*api.ListUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUser not implemented")
}
func (s *GRPCServer) GetUser(context context.Context, request *api.GetUserRequest) (*api.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
