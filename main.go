package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/admin"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"

	"github.com/qor/assetfs"
	_ "github.com/sergolius/qor_bindatafs_example/config/bindatafs"
)

type User struct {
	gorm.Model
	Name string
}

type Article struct {
	gorm.Model
	Title   string
	Content string
}

type TrialRequest struct {
	gorm.Model
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func (trial TrialRequest) Validate(DB *gorm.DB) {
	if len(trial.Name) < 5 || len(trial.Name) > 10 {
		DB.AddError(validations.NewError(trial, "Name", "name should be between 5 and 10"))
	}

	if len(trial.Email) < 5 || len(trial.Email) > 10 {
		DB.AddError(validations.NewError(trial, "Email", "email should be between 5 and 10"))
	}
}

func main() {
	ginEngine := gin.Default()

	DB, _ := gorm.Open(
		"sqlite3",
		"db.db",
	)
	DB.AutoMigrate(&User{}, &Article{}, &TrialRequest{})
	validations.RegisterCallbacks(DB)

	// Initialize AssetFS
	AssetFS := assetfs.AssetFS().NameSpace("admin")
	// Register custom paths to manually saved views
	AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "qor/admin/views"))

	Admin := admin.New(&admin.AdminConfig{
		DB:      DB,
		AssetFS: AssetFS,
	})

	Admin.AddResource(&User{})
	Admin.AddResource(&Article{})
	Admin.AddResource(&TrialRequest{})

	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)

	ginEngine.Any("/admin/*resources", gin.WrapH(mux))

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
		var articles = []Article{}
		var result = DB.Find(&articles)

		if result.Error != nil {
			returnPage500(context, result.Error.Error())
		}

		context.HTML(http.StatusOK, "blog.html", gin.H{
			"meta": gin.H{
				"title": "Blogus Listus",
			},
			"articles": articles,
		})
	})
	ginEngine.GET("/blog/:id", func(context *gin.Context) {
		var articleId = context.Param("id")
		var article = Article{}

		var result = DB.Where("ID = ?", articleId).First(&article)

		if result.Error != nil {
			returnPage500(context, result.Error.Error())
		}

		context.HTML(http.StatusOK, "content.html", gin.H{
			"meta": gin.H{
				"title": article.Title,
			},
			"content": article,
		})
	})

	ginEngine.POST("/trial", func(context *gin.Context) {
		var trialReqiest = TrialRequest{
			Name:  context.PostForm("name"),
			Email: context.PostForm("email"),
		}

		var result = DB.Create(&trialReqiest)

		var errors = result.GetErrors()

		if len(errors) > 0 {
			context.JSON(http.StatusInternalServerError, gin.H{
				"errors": errors,
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"message": "Ваша заявка принята",
			})
		}
	})

	ginEngine.NoRoute(func(context *gin.Context) {
		returnPage404(context, "Нет такой страницы")
	})

	ginEngine.Run(":9000")
}

func strToHTML(str string) template.HTML {
	return template.HTML(str)
}

func returnPage404(context *gin.Context, message string) {
	context.HTML(http.StatusNotFound, "404.html", gin.H{
		"message": message,
	})
}

func returnPage500(context *gin.Context, message string) {
	context.HTML(http.StatusInternalServerError, "500.html", gin.H{
		"message": message,
	})
}
