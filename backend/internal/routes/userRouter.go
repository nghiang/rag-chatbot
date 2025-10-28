package routes

import (
	"backend/internal/controllers"
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", controllers.SignUp())
		userGroup.POST("/login", controllers.Login())
		userGroup.GET("/", middleware.Authentication(), controllers.ListUsers())
		userGroup.GET("/:id", middleware.Authentication(), controllers.GetUserByID())
	}
}
