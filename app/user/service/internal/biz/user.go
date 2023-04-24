package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	v1 "github.com/user-center/user-center-backend/api/user/service/v1"
	"github.com/user-center/user-center-backend/app/user/service/internal/conf"
	"time"
	//v1 "github.com/user-center/user-center-backend/api/user/service/v1"
)

type UserRepo interface {
	SearchUsers(ctx context.Context, userName string) ([]*User, error)
	DeleteUser(ctx context.Context, userName int32) error
	GetUserRoleById(ctx context.Context, userId int32) (int32, error)
	GetUserSession(ctx context.Context, userId int32) (*User, error)
	GetCurrentUser(ctx context.Context, userId int32) (*User, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
	re   Recovery
	tm   Transaction
	conf *conf.UserConstant
}

//easyjson:json
type User struct {
	Id           int32
	UserName     string
	UserAccount  string
	AvatarUrl    string
	Gender       int32
	UserPassword string
	Phone        string
	Email        string
	UserStatus   int32
	Role         int32
	CreateTime   time.Time
}

type SearchUser struct {
	UserName string
}

type DeleteUser struct {
	Id int32 `validate:"required,gt=0" comment:"用户Id"`
}

func NewUserUseCase(repo UserRepo, re Recovery, tm Transaction, logger log.Logger, conf *conf.UserConstant) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/userUseCase")),
		tm:   tm,
		re:   re,
		conf: conf,
	}
}

func (i *User) DoValidate(trans ut.Translator, validate *validator.Validate) error {
	err := validate.Struct(i)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return errors.New(e.Translate(trans))
		}
	}
	return nil
}

// SearchUsers 用户搜索逻辑
//1. 判断是否有管理员权限
//2. 根据用户名进行模糊查询
func (r *UserUseCase) SearchUsers(ctx context.Context, userName string) ([]*User, error) {
	adminId := ctx.Value("userId").(int32)
	err := r.isAdmin(ctx, adminId)
	if err != nil {
		return nil, err
	}

	users, err := r.repo.SearchUsers(ctx, userName)
	if err != nil {
		return nil, v1.ErrorUserSearchFailed("%s", err.Error())
	}
	return users, nil
}

// DeleteUser 用户删除逻辑
//1. 判断是否有管理员权限
//2. 根据用户id进行逻辑删除
func (r *UserUseCase) DeleteUser(ctx context.Context, userId int32) error {
	adminId := ctx.Value("userId").(int32)
	err := r.isAdmin(ctx, adminId)
	if err != nil {
		return err
	}

	err = r.repo.DeleteUser(ctx, userId)
	if err != nil {
		return v1.ErrorUserDeleteFailed("%s", err.Error())
	}
	return nil
}

// GetCurrentUser 当前登录用户获取逻辑
//1. 判断session是否存在
//2. 如果存在，从数据库中获取最新用户信息返回
func (r *UserUseCase) GetCurrentUser(ctx context.Context) (*User, bool, error) {
	userId := ctx.Value("userId").(int32)
	exist, err := r.isSessionExist(ctx, userId)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, true, err
	}
	user, err := r.repo.GetCurrentUser(ctx, userId)
	if err != nil {
		return nil, false, v1.ErrorUnknownError("%s", err.Error())
	}
	return user, false, nil
}

func (r *UserUseCase) isAdmin(ctx context.Context, userId int32) error {
	role, err := r.repo.GetUserRoleById(ctx, userId)
	if kerrors.IsNotFound(err) {
		return v1.ErrorLoginStateTimeout("")
	}

	if err != nil {
		return v1.ErrorUnknownError("%s", err.Error())
	}

	if role != r.conf.AdminRole {
		return v1.ErrorPermissionDeny("userId: %s, userRole: %v", userId, role)
	}
	return nil
}

func (r *UserUseCase) isSessionExist(ctx context.Context, userId int32) (bool, error) {
	_, err := r.repo.GetUserSession(ctx, userId)
	if kerrors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, v1.ErrorUnknownError("%s", err.Error())
	}
	return true, nil
}
