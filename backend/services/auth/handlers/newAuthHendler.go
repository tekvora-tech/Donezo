package handlers

import (
	"auth/controllers"
	"platform/go-pkg/validator"
)

type (
	AuthHandler struct {
		controller controllers.IAuthController
		validator *validator.Validator
	}
)

func NewAuthHandler(ctrl controllers.IAuthController, v *validator.Validator) *AuthHandler {
	return &AuthHandler{
		controller: ctrl,
		validator: v,
	}
}