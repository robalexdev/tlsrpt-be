package main

import (
	"errors"
	"gorm.io/gorm"
	"log"
	model "github.com/robalexdev/tlsrpt-be/model"
	report "github.com/robalexdev/tlsrpt-be/report"
)


func processReportForUser(user *model.User, body []byte, db *gorm.DB) error {
	policyMap, err := report.Transform(body)
	if err != nil {
		log.Printf("Unable to process report: %v\n", err)
		return errors.New("Unable to process report")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for fqdn, policies := range policyMap {
			var domain model.Domain
			result := tx.Where("fqdn = ? AND user_id = ?", fqdn, user.ID).Take(&domain)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// Be helpful, create the domain for the user
					domain = model.Domain{
						UserID: user.ID,
						FQDN:   fqdn,
						Validated: false,
					}
					result := tx.Create(&domain)
					if result.Error != nil {
						log.Printf("Database error: %v\n", result.Error)
						return errors.New("Database error")
					}
				} else {
					log.Printf("Database error: %v\n", result.Error)
					return errors.New("Database error")
				}
			}

			for _, policy := range policies {
				policy.ManualUpload = true
				policy.DomainID = domain.ID
				result := tx.Create(&policy)
				if result.Error != nil {
					log.Printf("Database error: %v\n", result.Error)
					return errors.New("Database error")
				}
			}
		}
		return nil
	})
}
