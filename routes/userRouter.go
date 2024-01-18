package routes

import (
	"github.com/gin-gonic/gin"

	"jwt/controllers"
	middleware "jwt/middleware"
)

func USERRoutes(incoming_routes *gin.Engine) {

	incoming_routes.Use(middleware.Authenticate())
	// incoming_routes.GET("/users", controllers.GetUsers())
	incoming_routes.GET("/user/:id", controllers.GetUser())
}
