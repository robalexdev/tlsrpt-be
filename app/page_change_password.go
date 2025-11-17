package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func handleShowChangePassword(c *gin.Context, db *gorm.DB, user *User) {
	view := View{c}
	view.ChangePassword(*user, "")
}

func handleChangePassword(c *gin.Context, db *gorm.DB, user *User) {
	view := View{c}

	var form ChangePasswordForm
	if err := c.ShouldBind(&form); err != nil {
		// TODO: need better error message for user
		// TODO: should do client side validation
		view.ChangePassword(*user, err.Error())
		return
	}
	if form.Password1 != form.Password2 {
		view.ChangePassword(*user, "Passwords did not match")
		return
	}

	// The old password must match
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(form.OldPassword)); err != nil {
		view.ChangePassword(*user, "Invalid username or password")
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password1), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "change_password.tmpl", gin.H{
			"error": "Server Error",
		})
		return
	}

	// Save new password
	user.PasswordHash = hash
	if err := db.Save(&user).Error; err != nil {
		log.Println(err.Error())
		view.ChangePasswordServerError(*user)
		return
	}

	// TODO: message the user about password change (in web, via email)
	view.RedirectHome()
}
