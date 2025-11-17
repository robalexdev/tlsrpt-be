package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	//"compress/gzip"
	model "github.com/robalexdev/tlsrpt-be/model"
)

func handleUploadReport(c *gin.Context, db *gorm.DB, user *model.User) {
	view := View{c}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("Report can't load form file")
		view.ShowUploadForm(*user, "No file provided")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		// TODO better debug logging
		log.Println("Report can't open")
		view.ShowUploadForm(*user, "No file provided")
		return
	}
	defer file.Close()

	// TODO: gzip handling

	content, err := io.ReadAll(file)
	if err != nil {
		log.Println("Report can't read")
		view.ShowUploadForm(*user, "No file provided")
		return
	}

	log.Println("Report: " + string(content[:]))
	err = processReportForUser(user, content, db)
	if err != nil {
		view.ShowUploadForm(*user, err.Error())
		return
	}

	// TODO: Show the report
	view.RedirectHome()
}

func rejectUploadReport(user *model.User, errorMsg string, c *gin.Context) {
	c.HTML(http.StatusOK, "report_upload_form.tmpl", gin.H{
		"username": user.Email,
		"error":    errorMsg,
	})
}

func handleShowUploadReport(c *gin.Context, db *gorm.DB, user *model.User) {
	view := View{c}
	view.ShowUploadForm(*user, "")
}
