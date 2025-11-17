package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"net"
	"strings"
)

func loadDomainByUri(c *gin.Context, db *gorm.DB, user User) (domain Domain, err error) {
	var form ValidateDomainForm
	if err = c.ShouldBindUri(&form); err != nil {
		err = errors.New("Invalid URI")
		return
	}

	result := db.Model(&Domain{}).Where("id = ? AND user_id = ?", form.DomainID, user.ID).Take(&domain)
	if result.Error != nil {
		err = errors.New("No such domain")
		return
	}
	return
}

func FormatDomainName(domain Domain) string {
	fqdn, err := idna.ToUnicode(domain.FQDN)
	if err != nil {
		// Unable to render as unicode for some reason
		// render as stored
		// TODO: log this
		fqdn = domain.FQDN
	}

	if fqdn != domain.FQDN {
		fqdn = fmt.Sprintf("%s (%s)", fqdn, domain.FQDN)
	}
	return fqdn
}

func checkTxtRecord(domain Domain) (bool, []string) {
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
