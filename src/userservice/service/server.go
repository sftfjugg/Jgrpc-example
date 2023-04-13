package service

import (
	"context"
	"errors"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	userV1 "userservice/genproto/go/v1"
)

// Server Server struct
type Server struct {
	userV1.UnimplementedUserServiceServer
	userClient userV1.UserServiceClient
	repo       *Repository
}

// NewServer New service grpc server
func NewServer(repo *Repository, userClient userV1.UserServiceClient) userV1.UserServiceServer {
	return &Server{
		repo:       repo,
		userClient: userClient,
	}
}

// Register 用户注册
func (s *Server) Register(ctx context.Context, req *userV1.RegisterRequest) (*emptypb.Empty, error) {

	isExists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	if isExists {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, "用户名已存在")
	}

	_, err = s.repo.Register(req)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil

}

// Login 用户登录
func (s *Server) Login(ctx context.Context, req *userV1.LoginRequest) (*userV1.LoginResponse, error) {

	result, err := s.repo.Login(ctx, req)
	loginResp := &userV1.LoginResponse{}

	if err != nil {
		return loginResp, status.Error(codes.FailedPrecondition, "账号或密码错误")
	}

	// 返回数据
	loginResp.AccessToken = result.AccessToken
	loginResp.Username = req.Username
	loginResp.ExpireIn = result.AccessTokenExpireTime

	return loginResp, nil
}

// Logout 用户退出登录
func (s *Server) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	accessToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, errors.New("退出登录失败，错误：获取头部 access token 失败")
	}
	result, err := s.repo.Logout(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("退出登录失败")
	}
	return &emptypb.Empty{}, nil

}

// Info 获取用户信息
func (s *Server) Info(ctx context.Context, _ *emptypb.Empty) (*userV1.UserDetail_Detail, error) {

	resp := &userV1.UserDetail_Detail{}
	accessToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Error(codes.Unknown, "获取 access token 失败")
	}
	info, err := s.repo.Info(accessToken)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	resp = &userV1.UserDetail_Detail{
		Id:                    info.ID,
		Username:              info.Username,
		Sex:                   info.Sex,
		IdNumber:              info.IDNumber,
		Email:                 info.Email,
		Phone:                 info.Phone,
		IsDisable:             info.IsDisable,
		AccessToken:           info.AccessToken,
		AccessTokenExpireTime: info.AccessTokenExpireTime,
		NickName:              info.NickName,
		RealName:              info.RealName,
		CreateTime:            info.CreateTime,
		UpdateTime:            info.UpdateTime,
	}
	return resp, nil

}
