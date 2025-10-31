package main

import (
	"backend/config"
	_ "backend/docs" // import để register docs
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/routes"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Chatbot RAG API
// @version 1.0
// @description API for Chatbot with RAG functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := services.InitDB(); err != nil {
		panic(err)
	}
	// run AutoMigrate for all models
	if err := services.DB.AutoMigrate(&models.User{}, &models.KnowledgeBase{}, &models.Document{}, &models.ChatSession{}, &models.ChatMessage{}); err != nil {
		panic(err)
	}
	defer func() {
		sqlDB, err := services.DB.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.Close()
	}()
	defer services.CloseDB()
	port := config.LoadConfig().Port
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.KnowledgeBaseRoutes(router)
	routes.ChatSessionRoutes(router)
	routes.ChatMessageRoutes(router)
	routes.DocumentRoutes(router)

	router.Run(":" + port)
}
