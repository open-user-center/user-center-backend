package service

import (
	"context"
	v1 "github.com/user-center/user-center-backend/api/user/service/v1"
	"github.com/user-center/user-center-backend/app/user/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) SearchUsers(ctx context.Context, req *v1.SearchUsersReq) (*v1.SearchUsersReply, error) {
	usersList, err := s.uc.SearchUsers(ctx, req.UserName)
	if err != nil {
		return nil, err
	}

	reply := &v1.SearchUsersReply{
		Data: make([]*v1.User, 0, len(usersList)),
	}
	for _, item := range usersList {
		reply.Data = append(reply.Data, &v1.User{
			Id:          item.Id,
			UserName:    item.UserName,
			UserAccount: item.UserAccount,
			AvatarUrl:   item.AvatarUrl,
			Phone:       item.Phone,
			Email:       item.Email,
			UserStatus:  item.UserStatus,
			Gender:      item.Gender,
			UserRole:    item.Role,
			CreateTime:  item.CreateTime.String(),
		})
	}
	return reply, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserReq) (*emptypb.Empty, error) {
	search := &biz.DeleteUser{
		Id: req.Id,
	}
	err := s.vc.ParamsValidate(search)
	if err != nil {
		return nil, err
	}

	err = s.uc.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetCurrentUser(ctx context.Context, _ *emptypb.Empty) (*v1.GetCurrentReply, error) {
	user, empty, err := s.uc.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if empty {
		return &v1.GetCurrentReply{
			Data: &v1.User{
				Empty: empty,
			},
		}, nil
	}
	return &v1.GetCurrentReply{
		Data: &v1.User{
			Empty:       false,
			Id:          user.Id,
			UserName:    user.UserName,
			UserAccount: user.UserAccount,
			AvatarUrl:   user.AvatarUrl,
			Phone:       user.Phone,
			Email:       user.Email,
			UserStatus:  user.UserStatus,
			Gender:      user.Gender,
			UserRole:    user.Role,
		},
	}, nil
}
