package handlers

import (
	"auth/types"
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	log.Printf("hit service refresh token with request: %v", c.Request)

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("failed get cookie: %v", err)
		response.BadRequest(c, "invalid request cookie", err.Error())
		return
	}

	log.Println("running controller")
	result, err := h.controller.RefreshToken(c.Request.Context(), &types.ReqRefreshToken{
		RefreshToken: refreshToken,
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

	log.Println("success refresh token")
	response.OK(c, "success refresh token", result)
}