package controller

import (
	"context"

	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/repo"
)

type ServiceServer struct {
	sso.UnimplementedSSOServiceServer
}

func NewSSOServiceServer() sso.SSOServiceServer {
	return &ServiceServer{}
}

func (*ServiceServer) GetUserBasicInfo(ctx context.Context, req *sso.QueryUserInfoRequest) (*sso.User, error) {
	user, err := repo.GetBasicUserById(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	return user, nil
}
