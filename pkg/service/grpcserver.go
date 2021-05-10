package service

import (
	"context"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/service/errorhandler"
	"github.com/fev0ks/UserServiceSC/pkg/service/postgres"
	"github.com/fev0ks/UserServiceSC/pkg/service/validation"
)

type GRPCServer struct {
	api.UnimplementedUserServiceServer
}

func (s *GRPCServer) CreateUser(context context.Context, request *api.CreateUserRequest) (*api.User, error) {
	if err := validation.ValidateCreateUserRequestData(request); err != nil {
		return nil, errorhandler.NewInvalidArgumentError(err.Error())
	}
	return postgres.CreateUser(request)
}
func (s *GRPCServer) UpdateUser(context context.Context, request *api.UpdateUserRequest) (*api.User, error) {
	if err := validation.ValidateUserRequestData(request); err != nil {
		return nil, errorhandler.NewInvalidArgumentError(err.Error())
	}
	return postgres.UpdateUser(request)
}
func (s *GRPCServer) DeleteUser(context context.Context, request *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	if err := validation.ValidateId(request); err != nil {
		return nil, errorhandler.NewInvalidArgumentError(err.Error())
	}
	return postgres.DeleteUser(request)
}
func (s *GRPCServer) ListUser(context context.Context, request *api.ListUserRequest) (*api.ListUserResponse, error) {
	if err := validation.ValidatePageFilter(request); err != nil {
		return nil, errorhandler.NewInvalidArgumentError(err.Error())
	}
	return postgres.ListUser(request)
}
func (s *GRPCServer) GetUser(context context.Context, request *api.GetUserRequest) (*api.User, error) {
	if err := validation.ValidateId(request); err != nil {
		return nil, errorhandler.NewInvalidArgumentError(err.Error())
	}
	return postgres.GetUser(request)
}
