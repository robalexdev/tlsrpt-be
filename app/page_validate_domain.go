package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

func handleValidateDomain(c *gin.Context, db *gorm.DB, user *User, domain *Domain) {
	view := View{c}

	isValid, txtValues := checkTxtRecord(*domain)
	log.Printf("%s: %b %v", domain.FQDN, isValid, txtValues)

	if isValid {
		// We only allow a single domain object to be valid at any time
		// Domain may belong to a previous customer, remove validation
		db.Model(&Domain{}).Where("fqdn = ?", domain.FQDN).Update("validated", false)
	}

	// Update the domain with (possibly new) status and timestamp
	domain.Validated = isValid
	domain.DNSLastChecked = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	db.Save(&domain)

	if isValid {
		view.RedirectDomain(*domain)
	} else {
		view.ValidateDomain(*user, *domain, txtValues)
	}
}
