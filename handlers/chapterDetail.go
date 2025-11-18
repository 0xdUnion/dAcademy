package handlers

import (
	"dAcademy/database"
	"dAcademy/models"
	"dAcademy/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ChapterDetailHandler(c *gin.Context) {
	courseSlug := c.Param("courseSlug")
	chapterID := c.Param("chapterID")

	var course models.CourseData

	db, err := database.Run()
	err = db.Get(&course, `
        SELECT slug, name, description, tags, folder, chapter_count
        FROM courses
        WHERE slug = ?
        LIMIT 1
    `, courseSlug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err})
	}

	var courseFolder = course.Folder
	var courseName = course.Name

	// Read _chapters.yaml in that course folder
	chaptersFile := filepath.Join("./courses", courseFolder, "_chapters.yaml")
	var chapters []models.ChapterData
	if err := utils.ReadYAML(chaptersFile, &chapters); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot read %s: %v", chaptersFile, err)})
		return
	}

	// Find chapter folder by ID
	var chapterFolder string
	var chapterTitle string
	for _, ch := range chapters {
		if fmt.Sprintf("%d", ch.ID) == chapterID {
			chapterFolder = ch.Folder
			chapterTitle = ch.Title
			break
		}
	}
	if chapterFolder == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "chapter not found"})
		return
	}

	// Read sections.yaml
	sectionsFile := filepath.Join("./courses", courseFolder, chapterFolder, "sections.yaml")
	var sections []models.SectionData
	if err := utils.ReadYAML(sectionsFile, &sections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot read %s: %v", sectionsFile, err)})
		return
	}

	// Read quiz.yaml
	quizFile := filepath.Join("./courses", courseFolder, chapterFolder, "quiz.yaml")
	var quiz []models.QuizData
	if err := utils.ReadYAML(quizFile, &quiz); err != nil {
		// If quiz file doesn't exist or can't be read, use an empty slice
		quiz = []models.QuizData{}
	}

	// Respond ✨
	c.JSON(http.StatusOK, gin.H{
		"course":        courseName,
		"chapter_id":    chapterID,
		"chapter_title": chapterTitle,
		"sections":      sections,
		"quiz":          quiz,
	})
}
