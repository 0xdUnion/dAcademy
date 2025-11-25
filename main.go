package main

import (
	"dAcademy/database"
	handlers2 "dAcademy/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	_, err := database.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("📚 Database ready ✔")

	r := gin.Default()
	authApi := r.Group("/api/auth")
	{
		authApi.POST("/signup", func(c *gin.Context) {
			handlers2.SignupHandler(c)
		})
		authApi.POST("/login", func(c *gin.Context) {
			handlers2.LoginHandler(c)
		})
	}
	api := r.Group("/api")
	{
		api.GET("/course/list", func(c *gin.Context) {
			handlers2.CourseListHandler(c)
		})
		api.GET("/course/scan", func(c *gin.Context) {
			handlers2.CourseScanHandler(c)
		})
		api.GET("/course/:slug", func(c *gin.Context) {
			handlers2.CourseDetailHandler(c)
		})
		api.GET("/chapter/:courseSlug/:chapterID", func(c *gin.Context) {
			handlers2.ChapterDetailHandler(c)
		})

		api.GET("/me", func(c *gin.Context) {
			handlers2.MeHandler(c)
		})

	}

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
