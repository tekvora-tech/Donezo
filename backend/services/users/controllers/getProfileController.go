package controllers

import (
	"context"
	"users/types"
)

func (c *userController) GetProfile(ctx context.Context, req *types.ReqGetProfile) (res *types.ResGetProfile, err error) {
	u, err := c.repo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &types.ResGetProfile{
		ID: u.ID,
		Email: u.Email,
		FullName: u.FullName,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}