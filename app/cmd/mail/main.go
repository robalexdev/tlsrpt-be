package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"log"
	"strings"
	"compress/gzip"
	"bytes"
	"io/ioutil"

	"github.com/mnako/letters"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	model "github.com/robalexdev/tlsrpt-be/model"
	report "github.com/robalexdev/tlsrpt-be/report"
)

const (
	CONTENT_TYPE_TLSRPT string = "multipart/report"
	CONTENT_TYPE_REPORT_TYPE string = "tlsrpt"
	MEDIA_TYPE_TLSRPT_JSON string = "application/tlsrpt+json"
	MEDIA_TYPE_TLSRPT_GZ string = "application/tlsrpt+gzip"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--test" {
		processEmail(nil, nil)
		return
	}

	// Parse the email as delivered by postfix
	recipient := os.Getenv("ORIGINAL_RECIPIENT")
	userId, _, success := strings.Cut(recipient, "@")
	if !success {
		// TODO: how does postfix interpret this error handling? Is it a bounce?
		log.Printf("Can't find user ID in: %s\n", recipient)
		return
	}

	// Very simple spam filter. User ID uses format "d-[int]"
	if ! strings.HasPrefix(userId, "d-") {
		log.Printf("Invalid 'to' address: %s\n", userId)
		return
	}

	domainIdStr := userId[2:]
	domainId, err := strconv.Atoi(domainIdStr)
	if err != nil {
		log.Printf("Invalid domain ID: %s\n", domainIdStr)
		return
	}

	db, err := setupDatabase()
	if err != nil {
		log.Printf("Unable to setup database: %v\n", err)
		return
	}

	// Find the domain matching this username (if any)
	var domain model.Domain
	result := db.First(&domain, domainId)
	if result.Error != nil {
		log.Printf("Unable to find domain for ID %s: %d\n", domainId, result.Error)
		return
	}
	processEmail(db, &domain)
}

func setupDatabase() (db *gorm.DB, err error) {
	// Postgres password is saved to a file (as ENVs aren't passed down)
	content, err := ioutil.ReadFile("/postgres-password.txt")
	if err != nil {
		return
	}
	password := string(content)

	// Database setup
	dsn := fmt.Sprintf("host=db user=postgres password=%s dbname=tlsrpt port=5432 sslmode=disable TimeZone=UTC", password)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}

func processEmail(db *gorm.DB, domain *model.Domain) {
	emailParser := letters.NewEmailParser(
		letters.
			WithFileFilter(
					func(
							cth letters.ContentTypeHeader,
							_ letters.ContentDispositionHeader,
					) bool {
						return slices.Contains([]string{
								MEDIA_TYPE_TLSRPT_JSON,
								MEDIA_TYPE_TLSRPT_GZ,
							}, strings.ToLower(cth.ContentType))
					},
			),
	)
	email, err := emailParser.Parse(os.Stdin)
	if err != nil {
		log.Printf("Unable to parse email for tag %v\n", err)
		return
	}

	// Make sure this is really a TLSRPT message
	if ! checkContentType(email) {
		return
	}

	// Make sure the domain matches
	if ! checkReportDomain(email, domain) {
		return
	}

	for _, attachment := range email.AttachedFiles {
		processAttachment(db, domain, attachment.ContentType, attachment.Data)
	}
	for _, attachment := range email.InlineFiles {
		processAttachment(db, domain, attachment.ContentType, attachment.Data)
	}
}

func checkContentType(email letters.Email) bool {
	if strings.ToLower(email.Headers.ContentType.ContentType) != CONTENT_TYPE_TLSRPT {
		log.Printf("Not a TLSRPT, ContentType: %v\n", email.Headers.ContentType)
		return false
	}
	reportType, found := email.Headers.ContentType.Params["report-type"]
	if ! found {
		log.Printf("Not a TLSRPT, ContentType: %v\n", email.Headers.ContentType)
		return false
	}
	if strings.ToLower(reportType) != CONTENT_TYPE_REPORT_TYPE {
		log.Printf("Not a TLSRPT, ContentType: %v\n", email.Headers.ContentType)
		return false
	}
	return true
}

func checkReportDomain(email letters.Email, domain *model.Domain) bool {
	reportDomain, found := email.Headers.ExtraHeaders["Tls-Report-Domain"]
	if ! found || len(reportDomain) != 1 {
		log.Printf("TLSRPT missing report domain: %v\n", email.Headers.ExtraHeaders)
		return false
	}

	// Can there be more than one?
	if domain != nil && domain.FQDN != reportDomain[0] {
		log.Printf("Report recv'd for wrong domain: %s %s\n", domain.FQDN, reportDomain)
		return false
	}
	return true
}

func processAttachment(db *gorm.DB, domain *model.Domain, contentType letters.ContentTypeHeader, data []byte) {
	if contentType.ContentType == MEDIA_TYPE_TLSRPT_GZ {
		r := bytes.NewReader(data)
		gz, err := gzip.NewReader(r)
		if err != nil {
			log.Printf("Unable to gzip: %v", err)
			return
		}
		data, err = ioutil.ReadAll(gz)
	}

	policyMap, err := report.Transform(data)
	if err != nil {
		log.Printf("JSON parsing error: %v\n", err)
		return
	}

	if domain != nil && db != nil {
		for fqdn, policies := range policyMap {
			if fqdn != domain.FQDN {
				log.Printf("Report received for the wrong domain: %s vs %s\n", fqdn, domain.FQDN)
				continue
			}
			for _, policy := range policies {
				policy.ManualUpload = false
				policy.DomainID = domain.ID
				result := db.Create(&policy)
				if result.Error != nil {
					log.Printf("Database error: %v\n", result.Error)
					return
				}
			}
		}
	} else {
		fmt.Printf("TESTING: Got: %v\n", policyMap)
	}
}

