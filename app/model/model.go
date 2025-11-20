package model

import (
	"database/sql"
	"fmt"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//
// User
//

type User struct {
	ID      uint  `gorm:"primaryKey"`
	Created int64 `gorm:"autoCreateTime"`

	Email          string `gorm:"uniqueIndex"`
	ValidatedEmail bool
	PasswordHash   sql.RawBytes
}

func (u *User) AfterDelete(tx *gorm.DB) error {
	result := tx.Clauses(clause.Returning{}).Where("user_id = ?", u.ID).Delete(&Session{})
	if result.Error != nil {
		return result.Error
	}
	result = tx.Clauses(clause.Returning{}).Where("user_id = ?", u.ID).Delete(&Domain{})
	return result.Error
}

type Session struct {
	ID      uint  `gorm:"primaryKey"`
	Created int64 `gorm:"autoCreateTime"`

	UserID        uint `gorm:"index"`
	User          User
	SessionIdHash sql.RawBytes
}

//
// Domain
//

type Domain struct {
	ID      uint  `gorm:"primaryKey"`
	Created int64 `gorm:"autoCreateTime"`

	UserID uint `gorm:"index"`
	User   User

	// Note, the same domain may exist on multiple users
	// But only one should be valid
	// TODO: Valdating one invalidates any others
	FQDN           string
	Validated      bool
	DNSLastChecked sql.NullTime
}

func (d *Domain) AfterDelete(tx *gorm.DB) error {
	result := tx.Clauses(clause.Returning{}).Where("domain_id = ?", d.ID).Delete(&Policy{})
	return result.Error
}

func (d *Domain) VerificationHostname() string {
	return fmt.Sprintf("_smtp._tls.%s", d.FQDN)
}

func (d *Domain) VerificationValue() string {
	return fmt.Sprintf("v=TLSRPTv1; rua=mailto:d-%d@tlsrpt.alexsci.com;", d.ID)
}

func (d *Domain) PrettyDomainName() string {
	return FormatDomainName(*d)
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

//
// TLSRPT data
//

type Policy struct {
	ID uint `gorm:"primaryKey"`
	Created int64 `gorm:"autoCreateTime"`

	// Metadata
	ManualUpload bool

	// Denormalized from Report
	// Report JSON
	OrganizationName string
	ContactInfo      string
	ReportId         string
	// JSON date-range
	StartDateTime sql.NullTime
	EndDateTime   sql.NullTime

	DomainID uint
	Domain   Domain

	// Summary
	TotalSuccessfulSessionCount uint
	TotalFailureSessionCount    uint

	// Policy
	PolicyType string
}

func (p *Policy) AfterDelete(tx *gorm.DB) error {
	result := tx.Clauses(clause.Returning{}).Where("policy_id = ?", p.ID).Delete(&PolicyFailureDetail{})
	return result.Error
}

type PolicyFailureDetail struct {
	ID uint `gorm:"primaryKey"`

	PolicyID uint
	Policy   Policy

	ResultType            string
	SendingMtaIp          string
	ReceivingMxHostname   string
	ReceivingMxHelo       string
	ReceivingIp           string
	FailedSessionCount    uint
	AdditionalInformation string
	FailureReasonCode     string
}
