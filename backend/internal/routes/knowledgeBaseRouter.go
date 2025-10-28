package routes

import (
	"backend/internal/controllers"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func KnowledgeBaseRoutes(router *gin.Engine) {
	kbGroup := router.Group("/knowledge-bases")
	{
		kbGroup.POST("/", middleware.Authentication(), controllers.CreateKnowledgeBase())
		kbGroup.GET("/", middleware.Authentication(), controllers.ListKnowledgeBases())
		kbGroup.GET("/:id", middleware.Authentication(), controllers.GetKnowledgeBaseByID())
		kbGroup.PUT("/:id", middleware.Authentication(), controllers.UpdateKnowledgeBase())
		kbGroup.DELETE("/:id", middleware.Authentication(), controllers.DeleteKnowledgeBase())
	}
}
