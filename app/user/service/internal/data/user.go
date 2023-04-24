package data

import (
	"context"
	"fmt"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/user-center/user-center-backend/app/user/service/internal/biz"
	"github.com/user-center/user-center-backend/app/user/service/internal/pkg/util"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/user")),
	}
}

func (r *userRepo) GetUserRoleById(ctx context.Context, userId int32) (int32, error) {
	user, err := r.getUserFromCache(ctx, userId)
	if err != nil {
		return 0, err
	}
	return user.Role, nil
}

func (r *userRepo) GetUserSession(ctx context.Context, userId int32) (*biz.User, error) {
	user, err := r.getUserFromCache(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := &biz.User{}
	util.StructAssign(result, user)
	return result, nil
}

func (r *userRepo) GetCurrentUser(ctx context.Context, userId int32) (*biz.User, error) {
	user := &User{
		Id: userId,
	}
	err := r.data.db.WithContext(ctx).Where("id = ?", userId).First(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("get user failed: userId(%v)", userId))
	}

	result := &biz.User{}
	util.StructAssign(result, user)
	return result, nil
}

// SearchUsers 查询用户（允许根据用户名查询，仅管理员可查询）
func (r *userRepo) SearchUsers(ctx context.Context, userName string) ([]*biz.User, error) {
	list := make([]*User, 0)
	var err error
	switch userName {
	case "":
		err = r.data.db.WithContext(ctx).Where("isDelete = 0").Find(&list).Error
	default:
		err = r.data.db.WithContext(ctx).Where("userName like ? and isDelete = 0", userName).Find(&list).Error
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to search users: userName(%s)", userName))
	}

	var search []*biz.User
	for _, item := range list {
		user := &biz.User{}
		util.StructAssign(user, item)
		search = append(search, user)
	}
	return search, nil
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, userId int32) error {
	user := &User{}
	user.Id = userId
	user.IsDelete = 1
	err := r.data.db.WithContext(ctx).Where("id = ? and isDelete = 0", userId).Delete(user).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete user: userId(%s)", userId))
	}
	return nil
}

func (r *userRepo) getUserFromCache(ctx context.Context, userId int32) (*User, error) {
	result, err := r.data.redisCli.Get(ctx, fmt.Sprintf("%s_%v", r.data.conf.UserLoginState, userId)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, kerrors.NotFound("user not found from cache", fmt.Sprintf("userId(%v)", userId))
	}

	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user from cache: userId(%v)", userId))
	}
	var cacheUser = &User{}
	err = cacheUser.UnmarshalJSON([]byte(result))
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: user(%v)", result))
	}
	return cacheUser, nil
}
