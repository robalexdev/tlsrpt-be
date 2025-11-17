package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GinRoute func(c *gin.Context, db *gorm.DB)
type GinUserRoute func(c *gin.Context, db *gorm.DB, user *User)
type GinDomainRoute func(c *gin.Context, db *gorm.DB, user *User, domain *Domain)

func decorateRequireUser(wrapped GinUserRoute) GinRoute {
	return func(c *gin.Context, db *gorm.DB) {
		view := View{c}
		user := IsLoggedIn(c, db)
		if user == nil {
			view.RedirectLogin()
			return
		}
		wrapped(c, db, user)
	}
}

func decorateRequireDomain(wrapped GinDomainRoute) GinRoute {
	return func(c *gin.Context, db *gorm.DB) {
		decorateRequireUser(func(c *gin.Context, db *gorm.DB, user *User) {
			view := View{c}
			domain, err := loadDomainByUri(c, db, *user)
			if err != nil {
				view.Domain404(*user)
				return
			}
			wrapped(c, db, user, &domain)
		})(c, db)
	}
}
