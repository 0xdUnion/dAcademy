package handlers

import (
	"dAcademy/database"
	"dAcademy/internal/models"
	"dAcademy/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CourseDetailHandler(c *gin.Context) {
	slug := c.Param("slug") // /api/course/detail/:slug
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	var course models.CourseData

	db, err := database.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Get(&course, `
        SELECT slug, name, description, tags, folder, chapter_count
        FROM courses
        WHERE slug = ?
        LIMIT 1
    `, slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err})
	}

	var chapters []models.ChapterData
	chaptersFile := filepath.Join("./courses", course.Folder, "_chapters.yaml")

	if err := utils.ReadYAML(chaptersFile, &chapters); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "course found",
		"course":   course,
		"chapters": chapters,
	})

}
