package handlers

import (
	"dAcademy/database"
	"dAcademy/utils"
	"log"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var usernameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Basic validation
	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password cannot be empty"})
		return
	}

	// Username format check
	if !usernameRegex.MatchString(username) {
		c.JSON(400, gin.H{"error": "invalid username format"})
		return
	}

	// Password format check
	if !utils.IsSha256Hex(password) {
		c.JSON(400, gin.H{"error": "invalid password format"})
		return
	}

	db, err := database.Run()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	// Query DB for user data
	var (
		storedID       int
		storedPassword string
	)
	err = db.QueryRowx(`SELECT id, password FROM users WHERE username = $1`, username).Scan(&storedID, &storedPassword)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	// Check password
	if storedPassword != password {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	var sessionID string

	for {
		sessionID, err = utils.RandomSecureString()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}

		// check collision
		var exists int
		err = db.Get(&exists, `SELECT COUNT(*) FROM sessions WHERE session = $1`, sessionID)
		if err != nil {
			log.Println("error:", err)
			c.JSON(500, gin.H{"error": "internal error"})
			return
		}

		if exists == 0 {
			break // no collision, safe to use
		}
		// loop again to generate a new session
	}

	sessionData := map[string]interface{}{
		"user_id":     storedID,
		"session":     sessionID,
		"create_time": time.Now().Unix(),
	}

	// Insert into sessions table
	_, err = db.NamedExec(`
    INSERT INTO sessions (user_id, session, create_time)
    VALUES (:user_id, :session, :create_time)
`, sessionData)

	if err != nil {
		log.Println("error:", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"status":   "ok",
		"username": username,
		"session":  sessionID,
	})
}

func SignupHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check empty
	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password cannot be empty"})
		return
	}

	// Username regex check
	if !usernameRegex.MatchString(username) || len(username) < 1 || len(username) > 8 {
		c.JSON(400, gin.H{"error": "invalid username format"})
		return
	}

	// SHA256 hex password check
	if !utils.IsSha256Hex(password) {
		c.JSON(400, gin.H{"error": "password must be a valid SHA256 hex string"})
		return
	}

	// Connect DB
	db, err := database.Run()
	if err != nil {
		log.Println("DB connect error:", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	// Check username existing (safe prepared query)
	var count int
	err = db.Get(&count, `SELECT COUNT(*) FROM users WHERE username = $1`, username)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if count > 0 {
		c.JSON(409, gin.H{"error": "username already exists"})
		return
	}

	user := map[string]interface{}{
		"username":  username,
		"password":  password,
		"join_time": time.Now().Unix(),
	}

	// Insert
	_, err = db.NamedExec(`
        INSERT INTO users (username, password, join_time)
        VALUES (:username, :password, :join_time)
    `, user)

	if err != nil {
		log.Println("insert error:", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
