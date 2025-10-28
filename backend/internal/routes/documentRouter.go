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
		documentGroup.POST("/:userId", middleware.Authentication(), controllers.CreateDocument())
		documentGroup.GET("/:userId", middleware.Authentication(), controllers.ListDocuments())
		documentGroup.GET("/:userId/:docId", middleware.Authentication(), controllers.GetDocumentByID())
		documentGroup.PUT("/:userId/:docId", middleware.Authentication(), controllers.UpdateDocument())
		documentGroup.DELETE("/:userId/:docId", middleware.Authentication(), controllers.DeleteDocument())
	}
}
