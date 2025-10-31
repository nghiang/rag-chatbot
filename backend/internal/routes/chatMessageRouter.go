package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ChatMessageRoutes(router *gin.Engine) {
	messageGroup := router.Group("/chat-sessions/:id/messages")
	{
		messageGroup.POST("/", middleware.Authentication(), controllers.CreateChatMessage())
		messageGroup.GET("/", middleware.Authentication(), controllers.ListChatMessages())
		messageGroup.GET("/:messageId", middleware.Authentication(), controllers.GetChatMessageByID())
		messageGroup.PUT("/:messageId", middleware.Authentication(), controllers.UpdateChatMessage())
		messageGroup.DELETE("/:messageId", middleware.Authentication(), controllers.DeleteChatMessage())
	}
}
