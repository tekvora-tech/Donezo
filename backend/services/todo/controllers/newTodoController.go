package controllers

import (
	"platform/go-pkg/jwt"
	"todo/repositories"
)

type todoController struct {
	todoRepo repositories.ITodoRepository
	tagRepo repositories.ITagRepository
	jwt jwt.Service
}

func NewTodoController(todoRepo repositories.ITodoRepository, tagRepo repositories.ITagRepository, jwt jwt.Service) ITodoController {
	return &todoController{todoRepo: todoRepo, tagRepo: tagRepo, jwt: jwt}
}