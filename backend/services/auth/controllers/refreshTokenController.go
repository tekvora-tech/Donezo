package controllers

import (
	"auth/types"
	"context"
	"fmt"
	"time"
)

func (c *authController) RefreshToken(ctx context.Context, req *types.ReqRefreshToken) (res *types.ResRefreshToken, err error) {
	claims, err := c.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: refresh token is invalid or expired")
	}

	exists, err := c.repo.ExistsByID(ctx, claims.UserID.String())
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("unauthorized: refresh token is invalid or expired")
	}

	accessToken, err := c.jwt.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	return &types.ResRefreshToken{
		AccessToken: accessToken,
		ExpiresIn: int(15 * time.Minute),
	}, nil
}