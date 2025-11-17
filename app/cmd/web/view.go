package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	model "github.com/robalexdev/tlsrpt-be/model"
)

type View struct {
	Context *gin.Context
}

func (v View) RedirectLogin() {
	v.Context.Redirect(http.StatusFound, "/signin")
}

func (v View) RedirectHome() {
	v.Context.Redirect(http.StatusFound, "/")
}

func (v View) RedirectDomain(domain model.Domain) {
	v.Context.Redirect(http.StatusFound, fmt.Sprintf("/domain/%d", domain.ID))
}

func (v View) ShowDomain(user model.User, domain model.Domain, policies []model.Policy) {

	v.Context.HTML(http.StatusOK, "show_domain.tmpl", gin.H{
		"user":     user,
		"domain":   domain,
		"fqdn":     domain.PrettyDomainName(),
		"policies": policies,
	})
}

func (v View) ValidateDomain(user model.User, domain model.Domain, txt []string) {

	v.Context.HTML(http.StatusOK, "validate_domain.tmpl", gin.H{
		"user":     user,
		"domain":   domain,
		"fqdn":     domain.PrettyDomainName(),
		"hostname": domain.VerificationHostname(),
		"expected": domain.VerificationValue(),
		"txt":      txt,
	})
}

func (v View) Domain404(user model.User) {
	v.Context.HTML(http.StatusNotFound, "validate_domain.tmpl", gin.H{
		"user": user,
	})
}

func (v View) ChangePassword(user model.User, errorMsg string) {
	code := http.StatusOK
	if errorMsg != "" {
		code = http.StatusBadRequest
	}

	v.Context.HTML(code, "change_password.tmpl", gin.H{
		"user":  user,
		"error": errorMsg,
	})
}

func (v View) ChangePasswordServerError(user model.User) {
	v.Context.HTML(http.StatusInternalServerError, "change_password.tmpl", gin.H{
		"error": "Server Error",
	})
}

func (v View) LoggedOutHome() {
	v.Context.HTML(http.StatusOK, "home.tmpl", gin.H{})
}

func (v View) LoggedInHome(user model.User, domains []model.Domain) {
	v.Context.HTML(http.StatusOK, "home_logged_in.tmpl", gin.H{
		"user":    user,
		"domains": domains,
	})
}

func (v View) ShowUploadForm(user model.User, errorMsg string) {
	code := http.StatusOK
	if errorMsg != "" {
		code = http.StatusBadRequest
	}
	v.Context.HTML(code, "report_upload_form.tmpl", gin.H{
		"user":  user,
		"error": errorMsg,
	})
}
