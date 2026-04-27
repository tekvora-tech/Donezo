package handlers

import (
	"log"
	"platform/go-pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	log.Printf("hit service logout with request: %v", c.Request)

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("failed get cookie: %v", err)
		response.BadRequest(c, "invalid request cookie", err.Error())
		return
	}

	log.Println("success logout")
	c.SetCookie("refresh_token", refreshToken, 0, "/", "", true, true)
	response.OK(c, "success logout", nil)
}