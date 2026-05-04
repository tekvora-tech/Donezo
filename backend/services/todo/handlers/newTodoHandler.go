package handlers

import (
	"platform/go-pkg/validator"
	"todo/controllers"
)

type TodoHandler struct {
	controller controllers.ITodoController
	validator *validator.Validator
}

func NewTodoHandler(ctrl controllers.ITodoController, v *validator.Validator) *TodoHandler {
	return &TodoHandler{controller: ctrl, validator: v}
}