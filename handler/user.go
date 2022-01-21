package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserhandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// fmt.Errorf("pack %v", err)
		log.Println(err.Error())

		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		log.Println(err)
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result := user.FormatUser(newUser, "tokentokentoken")

	response := helper.APIResponse("Account has been registered", http.StatusOK, true, result)

	c.JSON(http.StatusOK, response)
}
