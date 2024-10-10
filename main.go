package main

import (
	"avidlogic/controllers"
	"avidlogic/database"
	_ "avidlogic/docs"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AvidLogic API
// @version 1.0
// @description This is a user management API for AvidLogic.
// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	database.ConnectDB()
	defer database.CloseDB()

	router := gin.Default()

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes
	router.POST("/users", controllers.CreateUser)

	router.Run(":8080")
}
