package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ChatSessionRoutes(router *gin.Engine) {
	sessionGroup := router.Group("/chat-sessions")
	{
		sessionGroup.POST("/", middleware.Authentication(), controllers.CreateChatSession())
		sessionGroup.GET("/", middleware.Authentication(), controllers.ListChatSessions())
		sessionGroup.GET("/:id", middleware.Authentication(), controllers.GetChatSessionByID())
		sessionGroup.PUT("/:id", middleware.Authentication(), controllers.UpdateChatSession())
		sessionGroup.DELETE("/:id", middleware.Authentication(), controllers.DeleteChatSession())
	}
}
