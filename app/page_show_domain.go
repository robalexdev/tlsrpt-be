package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func handleShowDomain(c *gin.Context, db *gorm.DB, user *User, domain *Domain) {
	view := View{c}
	if !domain.Validated {
		view.ValidateDomain(*user, *domain, []string{})
		return
	}

	var policies []Policy
	result := db.Where("domain_id = ?", domain.ID).Find(&policies)
	if result.Error != nil {
		log.Printf("Unable to load policies %v\n", result.Error.Error())
		policies = []Policy{}
	}

	view.ShowDomain(*user, *domain, policies)
}
