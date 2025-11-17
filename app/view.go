package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (v View) RedirectDomain(domain Domain) {
	v.Context.Redirect(http.StatusFound, fmt.Sprintf("/domain/%d", domain.ID))
}

func (v View) ShowDomain(user User, domain Domain, policies []Policy) {

	v.Context.HTML(http.StatusOK, "show_domain.tmpl", gin.H{
		"user":     user,
		"domain":   domain,
		"fqdn":     domain.PrettyDomainName(),
		"policies": policies,
	})
}

func (v View) ValidateDomain(user User, domain Domain, txt []string) {

	v.Context.HTML(http.StatusOK, "validate_domain.tmpl", gin.H{
		"user":     user,
		"domain":   domain,
		"fqdn":     domain.PrettyDomainName(),
		"hostname": domain.VerificationHostname(),
		"expected": domain.VerificationValue(),
		"txt":      txt,
	})
}

func (v View) Domain404(user User) {
	v.Context.HTML(http.StatusNotFound, "validate_domain.tmpl", gin.H{
		"user": user,
	})
}

func (v View) ChangePassword(user User, errorMsg string) {
	code := http.StatusOK
	if errorMsg != "" {
		code = http.StatusBadRequest
	}

	v.Context.HTML(code, "change_password.tmpl", gin.H{
		"user":  user,
		"error": errorMsg,
	})
}

func (v View) ChangePasswordServerError(user User) {
	v.Context.HTML(http.StatusInternalServerError, "change_password.tmpl", gin.H{
		"error": "Server Error",
	})
}

func (v View) LoggedOutHome() {
	v.Context.HTML(http.StatusOK, "home.tmpl", gin.H{})
}

func (v View) LoggedInHome(user User, domains []Domain) {
	v.Context.HTML(http.StatusOK, "home_logged_in.tmpl", gin.H{
		"user":    user,
		"domains": domains,
	})
}

func (v View) ShowUploadForm(user User, errorMsg string) {
	code := http.StatusOK
	if errorMsg != "" {
		code = http.StatusBadRequest
	}
	v.Context.HTML(code, "report_upload_form.tmpl", gin.H{
		"user":  user,
		"error": errorMsg,
	})
}
