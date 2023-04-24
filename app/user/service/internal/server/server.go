package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"
	"strconv"
)

// ProviderSet is user providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)
var (
	ErrorsMsgMap = map[string]string{
		"UNKNOWN_ERROR":        "未知错误",
		"ACCOUNT_EXIST":        "账号已存在",
		"ACCOUNT_ILLEGAL":      "账号只能包含字母数字下划线",
		"USER_REGISTER_FAILED": "用户注册失败",
		"USER_LOGIN_FAILED":    "用户登录失败或账号不存在",
		"USER_DELETE_FAILED":   "用户删除失败",
		"PERMISSION_DENY":      "没有权限",
		"LOGIN_STATE_TIMEOUT":  "登录已过期，请重新登录",
		"USER_LOGOUT_FAILED":   "用户注销失败",
	}
)

func responseServer() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if err != nil {
				e, ok := err.(*errors.Error)
				if ok {
					if m, ok := ErrorsMsgMap[e.Reason]; ok {
						e.Message = m
					}
					return
				}
				e.Message = ErrorsMsgMap["UNKNOWN_ERROR"]
			}
			return
		}
	}
}

func requestServer() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				var userId int
				s := header.RequestHeader().Get("userId")
				if s != "" {
					userId, err = strconv.Atoi(s)
					if err != nil {
						return nil, err
					}
				}
				ctx = context.WithValue(ctx, "userId", int32(userId))
			}
			reply, err = handler(ctx, req)
			return
		}
	}
}
