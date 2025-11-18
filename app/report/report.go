package report

import (
	"database/sql"
	"encoding/json"
	"time"
	model "github.com/robalexdev/tlsrpt-be/model"
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

func Transform(body []byte) (output map[string][]model.Policy, err error) {
	output = map[string][]model.Policy{}
	var reportJSON ReportJSON
	err = json.Unmarshal(body, &reportJSON)
	if err != nil {
		return
	}

	for _, policyJSON := range reportJSON.Policies {
		policyModel := model.Policy{
			// ManualUpload: TODO
			// DomainID: TODO
			OrganizationName:            reportJSON.OrganizationName,
			ContactInfo:                 reportJSON.ContactInfo,
			ReportId:                    reportJSON.ReportId,
			StartDateTime:               processDatetime(reportJSON.DateRange.StartDatetime),
			EndDateTime:                 processDatetime(reportJSON.DateRange.EndDatetime),
			TotalSuccessfulSessionCount: policyJSON.Summary.TotalSuccessfulSessionCount,
			TotalFailureSessionCount:    policyJSON.Summary.TotalFailureSessionCount,
			PolicyType:                  policyJSON.Policy.PolicyType,
		}
		// TODO: policy failure detail

		fqdn := policyJSON.Policy.PolicyDomain
		previous, found := output[fqdn]
		if ! found {
			previous = []model.Policy{}
		}
		output[fqdn] = append(previous, policyModel)
	}
	return
}
