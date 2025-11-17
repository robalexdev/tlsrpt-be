package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func handleShowAddDomain(c *gin.Context, db *gorm.DB, user *User) {
	c.HTML(http.StatusOK, "add_domain.tmpl", gin.H{
		"user": user,
	})
}

func handleAddDomain(c *gin.Context, db *gorm.DB, user *User) {
	var form AddDomainForm
	if err := c.ShouldBind(&form); err != nil {
		rejectAddDomain(c, err.Error())
		return
	}

	normalized, err := idna.Registration.ToASCII(form.FQDN)
	if err != nil {
		rejectAddDomain(c, err.Error())
		return
	}

	if !strings.Contains(normalized, ".") {
		rejectAddDomain(c, "Please enter a valid domain name")
		return
	}

	// Don't add the same domain twice per user
	var existingUserDomainCount int64
	result := db.Model(&Domain{}).
		Where("fqdn = ? AND user_id = ?", normalized, user.ID).
		Count(&existingUserDomainCount)
	if result.Error != nil {
		rejectAddDomain(c, result.Error.Error())
		return
	}
	if existingUserDomainCount >= 1 {
		rejectAddDomain(c, "You've already added this domain")
		return
	}

	// Note: FQDN can be associated with another user at this point
	// hence, Validated is false
	model := Domain{
		// User
		User:      *user,
		FQDN:      normalized,
		Validated: false,
	}
	result = db.Create(&model)
	if result.Error != nil {
		rejectAddDomain(c, result.Error.Error())
		return
	}

	// TODO: show the added domain?
	c.Redirect(http.StatusFound, "/")
}

func rejectAddDomain(c *gin.Context, reason string) {
	c.HTML(http.StatusBadRequest, "add_domain.tmpl", gin.H{
		"error": reason,
	})
}
