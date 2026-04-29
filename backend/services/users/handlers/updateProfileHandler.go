package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"users/types"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) UpdateProfileHandler(c *gin.Context) {
	log.Printf("hit service update profile with request: %v", c.Request)

	userID := c.GetString("user_id")

	var parsedBody types.DTOUpdateProfile
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
	result, err := h.controller.UpdateProfile(c.Request.Context(), &types.ReqUpdateProfile{
		UserID: userID,
		FullName: parsedBody.FullName,
		AvatarURL: parsedBody.AvatarURL,
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

	log.Println("success update profile")
	response.OK(c, "success", result)
}