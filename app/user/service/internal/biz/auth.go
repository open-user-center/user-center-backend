package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	v1 "github.com/user-center/user-center-backend/api/user/service/v1"
	"regexp"
)

type AuthRepo interface {
	AccountExist(ctx context.Context, userAccount string) (bool, error)
	UserRegister(ctx context.Context, userAccount, passwordHash string) (int32, error)
	UserLogin(ctx context.Context, userAccount, passwordHash string) (*User, error)
	UserLogout(ctx context.Context, userId int32) error
	SetLoginSession(ctx context.Context, userInfo *User) error
}

type AuthRepoUseCase struct {
	repo AuthRepo
	log  *log.Helper
	re   Recovery
	tm   Transaction
}

// UserRegister DO对象，带简单校验
type UserRegister struct {
	UserAccount   string `validate:"required,min=4" comment:"用户名"`
	UserPassword  string `validate:"required,min=4,max=8" comment:"用户密码"`
	CheckPassword string `validate:"required,min=4,max=8" comment:"重复密码"`
}

// UserLogin DO对象，带简单校验
type UserLogin struct {
	UserAccount  string `validate:"required,min=4" comment:"用户名"`
	UserPassword string `validate:"required,min=4,max=8" comment:"用户密码"`
}

func NewAuthRepoUseCase(repo AuthRepo, re Recovery, tm Transaction, logger log.Logger) *AuthRepoUseCase {
	return &AuthRepoUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/AuthRepoUseCase")),
		tm:   tm,
		re:   re,
	}
}

// UserRegister 注册逻辑设计
//
//1. 用户在前端输入账户和密码、以及校验码（todo）
//2. 校验用户的账户、密码、校验密码，是否符合要求
//1. 非空
//2. 账户长度 **不小于** 4 位
//3. 密码就 **不小于** 8 位吧
//4. 账户不能重复
//5. 账户不包含特殊字符
//6. 密码和校验密码相同
//3. 对密码进行加密（密码千万不要直接以明文存储到数据库中）
//4. 向数据库插入用户数据
func (r *AuthRepoUseCase) UserRegister(ctx context.Context, userAccount, userPassword, checkPassword string) (int32, error) {
	// 1、密码一致性校验
	err := r.isPasswordEqlCheckPassword(userPassword, checkPassword)
	if err != nil {
		return 0, err
	}

	// 2、账户合法性校验
	err = r.validateAccountBeforeRegister(ctx, userAccount)
	if err != nil {
		return 0, err
	}

	// 3、加密
	passwordHash := r.passwordMD5Hash(userPassword)

	// 4、插入数据
	id, err := r.repo.UserRegister(ctx, userAccount, passwordHash)
	if err != nil {
		return 0, v1.ErrorUserRegisterFailed("%s", err.Error())
	}
	return id, nil
}

func (r *AuthRepoUseCase) isPasswordEqlCheckPassword(userPassword, checkPassword string) error {
	if userPassword != checkPassword {
		return v1.ErrorValidateError("两次密码不一致")
	}
	return nil
}

func (r *AuthRepoUseCase) validateAccountBeforeRegister(ctx context.Context, userAccount string) error {
	// 账户不能含有特殊字符
	pass, err := r.isAccountWordsValidate(userAccount)
	if err != nil {
		return v1.ErrorUnknownError("%s", err.Error())
	}
	if pass {
		return v1.ErrorAccountIllegal("account(%s) illegal!", userAccount)
	}

	// 账户不能重复
	exist, err := r.isAccountExist(ctx, userAccount)
	if err != nil {
		return v1.ErrorUnknownError("%s", err.Error())
	}
	if exist {
		return v1.ErrorAccountExist("account(%s) exist!", userAccount)
	}
	return nil
}

func (r *AuthRepoUseCase) isAccountExist(ctx context.Context, userAccount string) (bool, error) {
	return r.repo.AccountExist(ctx, userAccount)
}

func (r *AuthRepoUseCase) isAccountWordsValidate(userAccount string) (bool, error) {
	targetString := userAccount
	matchString := "[^a-zA-Z0-9_]+"
	match, err := regexp.MatchString(matchString, targetString)
	if err != nil {
		return false, errors.Wrapf(err, "strings match error")
	}
	return match, nil
}

func (r *AuthRepoUseCase) passwordMD5Hash(userAccount string) string {
	m := md5.New()
	m.Write([]byte(userAccount))
	return hex.EncodeToString(m.Sum(nil))

}

// UserLogin 登录逻辑
//
//1. 校验用户账户和密码是否合法
//	1. 非空
//	2. 账户长度不小于 4 位
//	3. 密码就不小于 8 位
//	4. 账户不包含特殊字符
//2. 校验密码是否输入正确，要和数据库中的密文密码去对比
//3. 用户信息脱敏，隐藏敏感信息，防止数据库中的字段泄露
//4. 我们要记录用户的登录态（session），将其存到服务器上（redis）
// 		cookie
//5. 返回脱敏后的用户信息
func (r *AuthRepoUseCase) UserLogin(ctx context.Context, userAccount, userPassword string) (*User, error) {
	// 1、账户合法性校验
	err := r.validateAccountBeforeLogin(ctx, userAccount)
	if err != nil {
		return nil, err
	}

	// 3、加密
	passwordHash := r.passwordMD5Hash(userPassword)

	// 4、登录
	user, err := r.repo.UserLogin(ctx, userAccount, passwordHash)
	if err != nil {
		return nil, v1.ErrorUserLoginFailed("%s", err.Error())
	}

	// 5、存储登录的session
	err = r.repo.SetLoginSession(ctx, user)
	if err != nil {
		return nil, v1.ErrorUserLoginFailed("set user login session failed: %s", err.Error())
	}

	return user, nil
}

// UserLogout 注销逻辑
//1. 移除redis中的session即可
func (r *AuthRepoUseCase) UserLogout(ctx context.Context) error {
	userId := ctx.Value("userId").(int32)
	err := r.repo.UserLogout(ctx, userId)
	if err != nil {
		return v1.ErrorUserLogoutFailed("%s", err.Error())
	}
	return nil
}

func (r *AuthRepoUseCase) validateAccountBeforeLogin(ctx context.Context, userAccount string) error {
	// 账户不能含有特殊字符
	pass, err := r.isAccountWordsValidate(userAccount)
	if err != nil {
		return v1.ErrorUnknownError("%s", err.Error())
	}
	if pass {
		return v1.ErrorAccountIllegal("account(%s) illegal!", userAccount)
	}
	return nil
}
