package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//GetEnvWithKey : get env value

func main() {

	helper.LoadEnv()

	dsn := "host=localhost user=postgres password=12345 dbname=bwastartup port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserhandler(userService, authService)

	sess := helper.ConnectAws()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	})

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/check-email", userHandler.CheckEmailIsAvailable)
	api.POST("/upload-photo", userHandler.UploadAvatar)

	router.Run()
}
