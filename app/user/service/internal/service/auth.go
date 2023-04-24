package service

import (
	"context"
	v1 "github.com/user-center/user-center-backend/api/user/service/v1"
	"github.com/user-center/user-center-backend/app/user/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*v1.UserRegisterReply, error) {
	register := &biz.UserRegister{
		UserAccount:   req.UserAccount,
		UserPassword:  req.UserPassword,
		CheckPassword: req.CheckPassword,
	}
	err := s.vc.ParamsValidate(register)
	if err != nil {
		return nil, err
	}
	id, err := s.ac.UserRegister(ctx, register.UserAccount, register.UserPassword, register.CheckPassword)
	if err != nil {
		return nil, err
	}
	return &v1.UserRegisterReply{
		Data: &v1.User{
			Id: id,
		},
	}, nil
}

func (s *UserService) UserLogin(ctx context.Context, req *v1.UserLoginReq) (*v1.UserLoginReply, error) {
	login := &biz.UserLogin{
		UserAccount:  req.UserAccount,
		UserPassword: req.UserPassword,
	}
	err := s.vc.ParamsValidate(login)
	if err != nil {
		return nil, err
	}
	user, err := s.ac.UserLogin(ctx, login.UserAccount, login.UserPassword)
	if err != nil {
		return nil, err
	}
	// 脱敏处理，只返回必要的字段
	return &v1.UserLoginReply{
		Data: &v1.User{
			Id:          user.Id,
			UserName:    user.UserName,
			UserAccount: user.UserAccount,
			AvatarUrl:   user.AvatarUrl,
			Phone:       user.Phone,
			Email:       user.Email,
			UserStatus:  user.UserStatus,
			Gender:      user.Gender,
		},
	}, nil
}

func (s *UserService) UserLogout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.ac.UserLogout(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
