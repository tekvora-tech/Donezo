package handlers

import (
	"platform/go-pkg/validator"
	"users/controllers"
)

type (
	UserHandler struct {
		controller controllers.IUserController
		validator *validator.Validator
	}
)

func NewUserHandler(ctrl controllers.IUserController, v *validator.Validator) *UserHandler {
	return &UserHandler{
		controller: ctrl,
		validator: v,
	}
}