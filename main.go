package main

import "github.com/gin-gonic/gin"
import ("jwt/routes")

func main() {
	PORT := ":9000"

	router := gin.New()
	router.Use(gin.Logger())
	
	routes.AUTHRoutes(router)
	routes.USERRoutes(router)

	// router.GET("/api-1", func(ctx *gin.Context) {
	// 	ctx.JSON(200,gin.H{"success":"access granted for api-1"})
	// })
	// router.GET("/api-2", func(ctx *gin.Context) {
	// 	ctx.JSON(200,gin.H{"success":"access granted for api-2"})
	// })
router.Run(PORT)



}