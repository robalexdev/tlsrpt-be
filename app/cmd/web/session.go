package main

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	model "github.com/robalexdev/tlsrpt-be/model"
)

const SESSION_COOKIE_NAME = "sess"
const COOKIE_DOMAIN = "tlsrpt.alexsci.com"

func StartSession(user model.User, c *gin.Context, db *gorm.DB) error {
	// >= 128 bits of entropy
	sessionToken := rand.Text()
	hash := sha256.Sum256([]byte(sessionToken))

	session := model.Session{
		User:          user,
		SessionIdHash: hash[:],
	}
	result := db.Create(&session)
	if result.Error != nil {
		log.Println("Unable to save session: ", result.Error.Error())
		return result.Error
	}

	c.SetCookie(
		SESSION_COOKIE_NAME,
		sessionToken,
		0, // Session cookie
		"/",
		COOKIE_DOMAIN,
		true, // Secure transport only
		true, // HTTP only
	)

	// TODO: don't log PII
	log.Println("Logged in: ", user.Email)
	return nil
}

func IsLoggedIn(c *gin.Context, db *gorm.DB) (user *model.User) {
	sessionToken, err := c.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return
	}

	hash := sha256.Sum256([]byte(sessionToken))

	var session model.Session
	result := db.
		Where("session_id_hash = ?", hash[:]).
		Preload("User").
		Take(&session)
	if result.Error != nil {
		log.Println("Unknown session " + result.Error.Error())
		return
	}
	user = &session.User
	return
}

func LogOut(c *gin.Context, db *gorm.DB) {
	sessionToken, err := c.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		log.Println("Not logged in" + err.Error())
		return
	}

	// Remove from database
	hash := sha256.Sum256([]byte(sessionToken))
	result := db.
		Where("session_id_hash = ?", hash[:]).
		Delete(&model.Session{})
	if result.Error != nil {
		log.Println("Unknown session " + result.Error.Error())
	}

	// Clear cookie
	c.SetCookie(
		SESSION_COOKIE_NAME,
		"",
		-1, // Delete
		"/",
		COOKIE_DOMAIN,
		true, // Secure transport only
		true, // HTTP only
	)
}
