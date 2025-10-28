package main

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/routes"
	"github.com/gin-gonic/gin"
	"os"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "backend/docs"  // import để register docs
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
	if err := database.InitDB(); err != nil {
		panic(err)
	}
	// run AutoMigrate for all models
	if err := models.MigrateUser(); err != nil {
		panic(err)
	}
	if err := models.MigrateDocument(); err != nil {
		panic(err)
	}
	if err := models.MigrateChatSession(); err != nil {
		panic(err)
	}
	if err := models.MigrateChatMessage(); err != nil {
		panic(err)
	}
	if err := models.MigrateKnowledgeBase(); err != nil {
		panic(err)
	}
	defer database.CloseDB()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.KnowledgeBaseRoutes(router)
	routes.ChatMessageRoutes(router)
	routes.ChatSessionRoutes(router)
	routes.DocumentRoutes(router)

	router.Run(":" + port)
}
