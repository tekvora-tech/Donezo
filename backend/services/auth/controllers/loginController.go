package controllers

import (
	"auth/types"
	"context"
	"fmt"
	"platform/go-pkg/hash"

	"github.com/google/uuid"
)

func (c *authController) Login(ctx context.Context, req *types.ReqLogin) (res *types.ResLogin, err error) {
	user, err := c.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !hash.ComparePasswordBool(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("unauthorized: email or password incorect")
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, err
	}

	tokenPair, err := c.jwt.GenerateTokenPair(userID, user.Email)
	if err != nil {
		return nil, err
	}

	return &types.ResLogin{
		User: types.User{
			ID: user.ID,
			Email: user.Email,
			FullName: user.FullName,
			AvatarURL: user.AvatarURL,
		},
		AccessToken: tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn: int(tokenPair.ExpiresIn),
	}, nil
}