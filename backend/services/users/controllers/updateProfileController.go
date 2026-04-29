package controllers

import (
	"context"
	"users/types"
)

func (c *userController) UpdateProfile(ctx context.Context, req *types.ReqUpdateProfile) (res *types.ResUpdateProfile, err error) {
	if err := c.repo.Update(ctx, req.UserID, req.FullName, req.AvatarURL); err != nil {
		return nil, err
	}

	u, err := c.repo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &types.ResUpdateProfile{
		ID: u.ID,
		Email: u.Email,
		FullName: u.FullName,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}