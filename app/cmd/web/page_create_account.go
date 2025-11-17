package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func handleCreateAccount(c *gin.Context, db *gorm.DB) {
	var form CreateAccountForm
	if err := c.ShouldBind(&form); err != nil {
		// TODO: better error message
		// TODO: client side validation
		rejectCreateAccount(c, err.Error())
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password1), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error": "Server Error",
		})
		return
	}

	model := model.User{
		Email:          form.Email,
		PasswordHash:   hash,
		ValidatedEmail: false,
	}
	result := db.Create(&model)
	if result.Error != nil {
		// TODO: Pretty rejection for taken username
		rejectCreateAccount(c, result.Error.Error())
		return
	}

	// User added to database, log them in immediately
	if err = StartSession(model, c, db); err != nil {
		rejectCreateAccount(c, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func rejectCreateAccount(c *gin.Context, reason string) {
	c.HTML(http.StatusBadRequest, "signup.html", gin.H{
		"error": reason,
	})
}
