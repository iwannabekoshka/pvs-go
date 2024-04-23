package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type articleType struct {
	Id          string
	Title       string
	Description string
	Content     string
}

var articles = []articleType{
	{Id: "1", Title: "Articlus 1", Description: "Coolus articlus aboutus somethingus goose", Content: "<p>lorem <b>pipsum</b></p>"},
	{Id: "2", Title: "AArticlus 2", Description: "Coolus articlus aboutus somethingus goose", Content: "<p>lorem <b>pipsumpipsum</b></p>"},
	{Id: "3", Title: "Arrticlus 3", Description: "Coolus articlus aboutus somethingus goose", Content: "<p>lorem <b>pipsum</b>pipsum</p>"},
	{Id: "4", Title: "Artticlus 4", Description: "Coolus articlus aboutus somethingus goose", Content: "<p>lorem <b>pipsum</b>pipsumpipsum</p>"},
}

func main() {
	ginEngine := gin.Default()

	ginEngine.Static("/public", "./public")

	ginEngine.SetFuncMap(template.FuncMap{
		"strToHTML": strToHTML,
	})
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
	ginEngine.GET("/blog", func(context *gin.Context) {
		context.HTML(http.StatusOK, "blog.html", gin.H{
			"meta": gin.H{
				"title": "Blogus Listus",
			},
			"articles": articles,
		})
	})
	ginEngine.GET("/blog/:id", func(context *gin.Context) {
		var articleId = context.Param("id")
		var article articleType

		for _, _article := range articles {
			if _article.Id == articleId {
				article = _article
				break
			}
		}

		context.HTML(http.StatusOK, "content.html", gin.H{
			"meta": gin.H{
				"title": article.Title,
			},
			"content": article,
		})
	})

	ginEngine.Run(":9000")
}

func strToHTML(str string) template.HTML {
	return template.HTML(str)
}
