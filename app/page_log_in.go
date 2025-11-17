package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func handleLogIn(c *gin.Context, db *gorm.DB) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		// TODO: obscure the error?
		rejectLogIn(c, err.Error())
		return
	}

	// Find the user
	var model User
	err := db.
		Where("email = ?", form.Email).
		Take(&model).
		Error

	if err != nil {
		log.Println("User not found: " + err.Error())
		rejectLogIn(c, "Invalid username or password")
		return
	}

	// Hash the provided password
	// TODO: timing attack leaks emails
	err = bcrypt.CompareHashAndPassword(model.PasswordHash, []byte(form.Password))
	if err != nil {
		rejectLogIn(c, "Invalid username or password")
		return
	}

	if err = StartSession(model, c, db); err != nil {
		rejectCreateAccount(c, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func rejectLogIn(c *gin.Context, reason string) {
	c.HTML(http.StatusBadRequest, "signin.html", gin.H{
		"error": reason,
	})
}
