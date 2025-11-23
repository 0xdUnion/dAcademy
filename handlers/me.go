package handlers

import (
	"dAcademy/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MeHandler(c *gin.Context) {
	// Get session token from cookie
	sessionToken, err := c.Cookie("d_session")
	if err != nil || sessionToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no session"})
		return
	}

	db, err := database.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// session->user_id
	var (
		userId   uint
		userName string
	)

	err = db.QueryRowx(`SELECT user_id FROM sessions WHERE session = $1`, sessionToken).Scan(&userId)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid session token"})
		return
	}

	// user_id->username
	err = db.QueryRowx(`SELECT username FROM users WHERE id = $1`, userId).Scan(&userName)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid user id"})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"username": userName,
	})
}
