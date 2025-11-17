package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type ReportJSON struct {
	OrganizationName string        `json:"organization-name"`
	DateRange        DateRangeJSON `json:"date-range"`
	ContactInfo      string        `json:"contact-info"`
	ReportId         string        `json:"report-id"`
	Policies         []PolicyJSON  `json:"policies"`
}

type DateRangeJSON struct {
	StartDatetime string `json:"start-datetime"`
	EndDatetime   string `json:"end-datetime"`
}

type PolicyJSON struct {
	Policy  PolicyInfoJSON `json:"policy"`
	Summary SummaryJSON    `json:"summary"`
}

type PolicyInfoJSON struct {
	PolicyType   string `json:"policy-type"`
	PolicyDomain string `json:"policy-domain"`
}

type SummaryJSON struct {
	TotalSuccessfulSessionCount uint `json:"total-successful-session-count"`
	TotalFailureSessionCount    uint `json:"total-failure-session-count"`
}

func processDatetime(s string) sql.NullTime {
	var parsed sql.NullTime
	t, err := time.Parse(time.RFC3339, s)
	parsed.Time = t
	parsed.Valid = err != nil
	return parsed
}

func processReportForUser(user *User, reportStr []byte, db *gorm.DB) error {
	var reportJSON ReportJSON
	err := json.Unmarshal(reportStr, &reportJSON)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to unmarshal JSON")
	}

	log.Println(reportJSON)

	db.Transaction(func(tx *gorm.DB) error {

		for _, policyJSON := range reportJSON.Policies {
			var domain Domain
			result := db.Where("fqdn = ? AND user_id = ?", policyJSON.Policy.PolicyDomain, user.ID).Take(&domain)
			if result.Error != nil {
				log.Println(err.Error())
				// TODO: just add the domain on-the-fly?
				return errors.New(fmt.Sprintf("Domain is not verified for user: %s", policyJSON.Policy.PolicyDomain))
			}

			policyModel := Policy{
				ManualUpload:                true,
				OrganizationName:            reportJSON.OrganizationName,
				ContactInfo:                 reportJSON.ContactInfo,
				ReportId:                    reportJSON.ReportId,
				StartDateTime:               processDatetime(reportJSON.DateRange.StartDatetime),
				EndDateTime:                 processDatetime(reportJSON.DateRange.EndDatetime),
				DomainID:                    domain.ID,
				TotalSuccessfulSessionCount: policyJSON.Summary.TotalSuccessfulSessionCount,
				TotalFailureSessionCount:    policyJSON.Summary.TotalFailureSessionCount,
				PolicyType:                  policyJSON.Policy.PolicyType,
			}
			err = tx.Create(&policyModel).Error
			if err != nil {
				return err
			}
			// TODO: policy failure detail
		}
		return nil
	})
	return nil
}
