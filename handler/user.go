package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserhandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		log.Println(err)
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, true, result)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loginUser, err := h.userService.Login(input)

	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := h.authService.GenerateToken(loginUser.ID)

	if err != nil {
		fmt.Println(err)
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	result := user.FormatUser(loginUser, token)

	response := helper.APIResponseLogin("Login Success", http.StatusOK, true, result)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailIsAvailable(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Email has been registered", http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		response := helper.APIResponse("Email has been registered", http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	meta := "Email has been registered"

	if isEmailAvailable {
		meta = "Email is Available"
	}

	response := helper.APIResponseLogin(meta, http.StatusBadRequest, false, data)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {

	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)
	MyBucket := helper.GetEnvWithKey("BUCKET_NAME")
	file, header, err := c.Request.FormFile("avatar")
	filename := "images/" + header.Filename

	fmt.Println(file)
	fmt.Println(filename)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": up,
		})
		return
	}

	// file, err := c.FormFile("avatar")

	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := filename

	userID := 1

	user, err := h.userService.SaveAvatar(userID, path)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, false, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponseLogin("Avatar succesfully uploaded", http.StatusOK, true, user)

	c.JSON(http.StatusOK, response)
}
