package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pong": "true",
		})
	})

	r.GET("/", func(c *gin.Context) {
		handleHomePage(c, db)
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})
	r.POST("/signup", func(c *gin.Context) {
		handleCreateAccount(c, db)
	})

	r.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", gin.H{})
	})
	r.POST("/signin", func(c *gin.Context) {
		handleLogIn(c, db)
	})

	r.GET("/changePassword", func(c *gin.Context) {
		decorateRequireUser(handleShowChangePassword)(c, db)
	})
	r.POST("/changePassword", func(c *gin.Context) {
		decorateRequireUser(handleChangePassword)(c, db)
	})

	r.POST("/signout", func(c *gin.Context) {
		decorateRequireUser(handleLogOut)(c, db)
	})

	r.GET("/domain/add", func(c *gin.Context) {
		decorateRequireUser(handleShowAddDomain)(c, db)
	})
	r.POST("/domain/add", func(c *gin.Context) {
		decorateRequireUser(handleAddDomain)(c, db)
	})
	r.GET("/domain/:id/", func(c *gin.Context) {
		decorateRequireDomain(handleShowDomain)(c, db)
	})
	r.POST("/domain/:id/", func(c *gin.Context) {
		decorateRequireDomain(handleValidateDomain)(c, db)
	})
	// TODO remove domain

	r.GET("/uploadReport", func(c *gin.Context) {
		decorateRequireUser(handleShowUploadReport)(c, db)
	})
	r.POST("/uploadReport", func(c *gin.Context) {
		decorateRequireUser(handleUploadReport)(c, db)
	})
	// TODO delete report
	// TODO delete account / purge data
}
