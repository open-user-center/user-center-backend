package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/user-center/user-center-backend/app/user/service/internal/biz"
	"github.com/user-center/user-center-backend/app/user/service/internal/pkg/util"
	"gorm.io/gorm"
	"time"
)

var _ biz.AuthRepo = (*authRepo)(nil)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/register")),
	}
}

func (r *authRepo) AccountExist(ctx context.Context, userAccount string) (bool, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Select("id").Where("userAccount = ?", userAccount).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to get account: userAccount(%v)", userAccount))
	}
	return true, nil
}

func (r *authRepo) UserRegister(ctx context.Context, userAccount, passwordHash string) (int32, error) {
	user := &User{
		UserAccount:  userAccount,
		UserPassword: passwordHash,
	}
	err := r.data.db.WithContext(ctx).Select("userAccount", "userPassword").Create(user).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to register user: userAccount(%s), userPassword(%s)", userAccount, passwordHash))
	}
	return user.Id, nil
}

func (r *authRepo) UserLogin(ctx context.Context, userAccount, passwordHash string) (*biz.User, error) {
	user := &User{
		UserAccount:  userAccount,
		UserPassword: passwordHash,
	}
	err := r.data.db.WithContext(ctx).Where("userAccount = ? and userPassword = ? and isDelete = 0", userAccount, passwordHash).First(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("user login failed: userAccount(%s), userPassword(%s)", userAccount, passwordHash))
	}

	result := &biz.User{}
	util.StructAssign(result, user)
	return result, nil
}

func (r *authRepo) UserLogout(ctx context.Context, userId int32) error {
	_, err := r.data.redisCli.Del(ctx, fmt.Sprintf("%s_%v", r.data.conf.UserLoginState, userId)).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("user login failed: userId(%v)", userId))
	}
	return nil
}

func (r *authRepo) SetLoginSession(ctx context.Context, user *biz.User) error {
	marshal, err := user.MarshalJSON()
	if err != nil {
		r.log.Errorf("fail to set user info to json: json.Marshal(%v), error(%v)", user, err)
		return nil
	}
	err = r.data.redisCli.Set(ctx, fmt.Sprintf("%s_%v", r.data.conf.UserLoginState, user.Id), string(marshal), time.Second*time.Duration(r.data.conf.SessionTimeout)).Err()
	if err != nil {
		r.log.Errorf("fail to set user session to cache: redis.Set(%v), error(%v)", user, err)
	}
	return nil
}
