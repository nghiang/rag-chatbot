package routes

import (
	"github.com/gin-gonic/gin"
	"backend/internal/controllers"
	"backend/internal/middleware"
)

// Placeholder: register document routes here
func DocumentRoutes(router *gin.Engine) {
	documentGroup := router.Group("/documents")
	{
		documentGroup.POST("/:knowledgeBaseId", middleware.Authentication(), controllers.CreateDocument())
		documentGroup.GET("/:knowledgeBaseId", middleware.Authentication(), controllers.ListDocuments())
		documentGroup.GET("/:knowledgeBaseId/:docId", middleware.Authentication(), controllers.GetDocumentByID())
		documentGroup.PUT("/:knowledgeBaseId/:docId", middleware.Authentication(), controllers.UpdateDocument())
		documentGroup.DELETE("/:knowledgeBaseId/:docId", middleware.Authentication(), controllers.DeleteDocument())
	}
}
