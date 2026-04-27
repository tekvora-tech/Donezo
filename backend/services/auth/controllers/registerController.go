package controllers

import (
	"auth/types"
	"context"
	"fmt"
	"platform/go-pkg/hash"
	"time"

	"github.com/google/uuid"
)

func (c *authController) Register(ctx context.Context, req *types.ReqRegister) (res *types.ResRegister, err error) {
	if err := hash.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("validation failed: minimal 8 karakter, harus mengandung minimal 1 huruf besar, 1 huruf kecil, 1 angka.")
	}

	passwordHash, err := hash.HashPasswordDefault(req.Password)
	if err != nil {
		return nil, err
	}

	userID := uuid.New()
	
	now := time.Now()
	if err := c.repo.Create(ctx, userID, req.Email, passwordHash, req.FullName); err != nil {
		return nil, err
	}

	tokenPair, err := c.jwt.GenerateTokenPair(userID, req.Email)
	if err != nil {
		return nil, err
	}

	return &types.ResRegister{
		User: types.User{
			ID: userID.String(),
			Email: req.Email,
			FullName: req.FullName,
			AvatarURL: nil,
			CreatedAt: &now,
		},
		AccessToken: tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn: int(tokenPair.ExpiresIn),
	}, nil
}