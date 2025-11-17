package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func handleShowDomain(c *gin.Context, db *gorm.DB, user *model.User, domain *model.Domain) {
	view := View{c}
	if !domain.Validated {
		view.ValidateDomain(*user, *domain, []string{})
		return
	}

	var policies []model.Policy
	result := db.Where("domain_id = ?", domain.ID).Find(&policies)
	if result.Error != nil {
		log.Printf("Unable to load policies %v\n", result.Error.Error())
		policies = []model.Policy{}
	}

	view.ShowDomain(*user, *domain, policies)
}
