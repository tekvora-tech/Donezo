package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"users/types"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetProfileHandler(c *gin.Context) {
	log.Printf("hit service get profile with request: %v", c.Request)

	userID := c.GetString("user_id")

	log.Println("running controller")
	result, err := h.controller.GetProfile(c.Request.Context(), &types.ReqGetProfile{
		UserID: userID,
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

	log.Println("success get profile")
	response.OK(c, "success", result)
}