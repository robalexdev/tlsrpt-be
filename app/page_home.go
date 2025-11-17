package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func handleHomePage(c *gin.Context, db *gorm.DB) {
	view := View{c}
	user := IsLoggedIn(c, db)
	if user == nil {
		view.LoggedOutHome()
		return
	}

	var domains []Domain
	result := db.Where("user_id = ?", user.ID).Find(&domains)

	if result.Error != nil {
		// TODO: error handling
		log.Println("Domain Count Failed")
		log.Println(result.Error.Error())
		domains = []Domain{}
	}
	view.LoggedInHome(*user, domains)
}
