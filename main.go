package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()

	ginEngine.Static("/public", "./public")
	ginEngine.LoadHTMLGlob("templates/**/*")

	ginEngine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "home.html", gin.H{
			"meta": gin.H{
				"title": "Homus",
			},
		})
	})
	ginEngine.GET("/about", func(context *gin.Context) {
		context.HTML(http.StatusOK, "about.html", gin.H{
			"meta": gin.H{
				"title": "Aboutus",
			},
		})
	})

	ginEngine.Run(":9000")
}
