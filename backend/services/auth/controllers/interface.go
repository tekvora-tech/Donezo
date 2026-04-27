package controllers

import (
	"auth/types"
	"context"
)

type (
	IAuthController interface {
		Register(ctx context.Context, req *types.ReqRegister) (res *types.ResRegister, err error)
		Login(ctx context.Context, req *types.ReqLogin) (res *types.ResLogin, err error)
		RefreshToken(ctx context.Context, req *types.ReqRefreshToken) (res *types.ResRefreshToken, err error)
	}
)