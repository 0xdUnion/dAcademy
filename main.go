package main

import (
	"dAcademy/database"
	"dAcademy/handlers"
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
	api := r.Group("/api")
	{
		api.GET("/course/list", func(c *gin.Context) {
			handlers.CourseListHandler(c)
		})
		api.GET("/course/scan", func(c *gin.Context) {
			handlers.CourseScanHandler(c)
		})
		api.GET("/course/:slug", func(c *gin.Context) {
			handlers.CourseDetailHandler(c)
		})
		api.GET("/chapter/:courseSlug/:chapterID", func(c *gin.Context) {
			handlers.ChapterDetailHandler(c)
		})

	}

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
