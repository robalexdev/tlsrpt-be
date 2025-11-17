package main

type AddDomainForm struct {
	FQDN string `form:"domain" binding:"required"`
}

type CreateAccountForm struct {
	Email     string `form:"email" binding:"required,email"`
	Password1 string `form:"password1" binding:"required,min=8,max=100"`
	Password2 string `form:"password2" binding:"required,min=8,max=100"`
	// CSRF other?
}

type ValidateDomainForm struct {
	DomainID uint `uri:"id" binding:"required"`
}

type LoginForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=8,max=100"`
}

type ChangePasswordForm struct {
	OldPassword string `form:"oldPassword" binding:"required,min=8,max=100"`
	Password1   string `form:"password1" binding:"required,min=8,max=100"`
	Password2   string `form:"password2" binding:"required,min=8,max=100"`
}
