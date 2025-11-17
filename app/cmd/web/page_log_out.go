package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func handleLogOut(c *gin.Context, db *gorm.DB, _ *model.User) {
	view := View{c}
	LogOut(c, db)
	view.RedirectHome()
}
