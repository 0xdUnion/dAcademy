package handlers

import (
	"dAcademy/database"
	"dAcademy/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CourseListHandler(c *gin.Context) {
	var courses []models.CourseData

	db, err := database.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Read data from database
	err = db.Select(&courses, `
        SELECT slug, name, description, tags, folder, chapter_count
        FROM courses
        ORDER BY id ASC
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Output JSON response
	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}
