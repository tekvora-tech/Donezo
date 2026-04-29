package controllers

import (
	"platform/go-pkg/jwt"
	"users/repositories"
)

type userController struct {
	repo repositories.IUserRepository
	jwt jwt.Service
}

func NewUserController(repo repositories.IUserRepository, jwt jwt.Service) IUserController {
	return &userController{repo: repo, jwt: jwt}
}