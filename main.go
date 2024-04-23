package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()

	ginEngine.LoadHTMLGlob("templates/*")

	ginEngine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "home.html", nil)
	})
	ginEngine.GET("/about", func(context *gin.Context) {
		context.HTML(http.StatusOK, "about.html", nil)
	})

	ginEngine.Run(":9000")
}
