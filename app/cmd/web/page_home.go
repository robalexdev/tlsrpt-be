package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func handleHomePage(c *gin.Context, db *gorm.DB) {
	view := View{c}
	user := IsLoggedIn(c, db)
	if user == nil {
		view.LoggedOutHome()
		return
	}

	var domains []model.Domain
	result := db.Where("user_id = ?", user.ID).Find(&domains)

	if result.Error != nil {
		// TODO: error handling
		log.Println("Domain Count Failed")
		log.Println(result.Error.Error())
		domains = []model.Domain{}
	}
	view.LoggedInHome(*user, domains)
}
