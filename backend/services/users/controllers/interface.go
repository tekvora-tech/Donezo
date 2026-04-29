package controllers

import (
	"context"
	"users/types"
)

type (
	IUserController interface {
		GetProfile(ctx context.Context, req *types.ReqGetProfile) (res *types.ResGetProfile, err error)
		UpdateProfile(ctx context.Context, req *types.ReqUpdateProfile) (res *types.ResUpdateProfile, err error)
	}
)