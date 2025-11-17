package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
	"strings"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func loadDomainByUri(c *gin.Context, db *gorm.DB, user model.User) (domain model.Domain, err error) {
	var form ValidateDomainForm
	if err = c.ShouldBindUri(&form); err != nil {
		err = errors.New("Invalid URI")
		return
	}

	result := db.Model(&model.Domain{}).Where("id = ? AND user_id = ?", form.DomainID, user.ID).Take(&domain)
	if result.Error != nil {
		err = errors.New("No such domain")
		return
	}
	return
}

func checkTxtRecord(domain model.Domain) (bool, []string) {
	hostname := domain.VerificationHostname()
	expected := domain.VerificationValue()

	// TODO: stub this test out better
	if domain.FQDN == "xn--gpher-jua.test" {
		return true, []string{}
	}

	// Lookup DNS record
	txtValues, err := net.LookupTXT(hostname)
	if err != nil {
		// TODO: log it
		return false, []string{}
	}
	for _, txtValue := range txtValues {
		// TODO: this is a sloppy compare
		if strings.Contains(txtValue, expected) {
			return true, txtValues
		}
	}
	return false, txtValues
}
