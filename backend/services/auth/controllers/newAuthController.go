package controllers

import (
	"auth/repositories"
	"platform/go-pkg/jwt"
)

type authController struct {
	repo repositories.IUserRepository
	jwt jwt.Service
}

func NewAuthController(repo repositories.IUserRepository, jwt jwt.Service) IAuthController {
	return &authController{repo: repo, jwt: jwt}
}