package handlers

import (
	"auth/types"
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	log.Printf("hit service login with request: %v", c.Request)

	var parsedBody types.DTOLogin
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		log.Printf("failed unmarshal request body: %v", err)
		response.BadRequest(c, "failed unmarshal request body", err.Error())
		return
	}

	if err := h.validator.Validate(parsedBody); err != nil {
		log.Printf("failed format request body: %v", err)
		response.BadRequest(c, "failed format request body", err.Error())
		return
	}

	log.Println("running controller")
	result, err := h.controller.Login(c.Request.Context(), &types.ReqLogin{
		Email: parsedBody.Email,
		Password: parsedBody.Password,
	})

	if err != nil {
		log.Printf("error on controller: %v", err)
		if parser.IsPostgresError(err) {
			httpError := parser.ParsePostgresError(err)
			response.Error(c, httpError.Code, httpError.Message, httpError.Err)
			return
		}

		msg, code, detail := parser.ParseError(err)
		response.Error(c, code, msg, detail)
		return
	}

	log.Println("success login")
	c.SetCookie("refresh_token", result.RefreshToken, result.ExpiresIn, "/", "", true, true)
	response.OK(c, "success login", result)
}