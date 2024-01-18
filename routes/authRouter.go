package routes

import ("github.com/gin-gonic/gin"
controllers "jwt/controllers"
)


func AUTHRoutes(incoming_routes *gin.Engine ){

incoming_routes.POST("users/login",controllers.Login())
incoming_routes.POST("users/signup",controllers.Signup())
}