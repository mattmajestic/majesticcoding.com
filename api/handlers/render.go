package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

func RenderWithClerk(tmpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, tmpl, gin.H{
			"CLERK_PUBLISHABLE_KEY": os.Getenv("CLERK_PUBLISHABLE_KEY"),
		})
	}
}

func RenderTemplate(tmpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, tmpl, gin.H{})
	}
}

func DocsHandler(r *gin.Engine) {
	r.GET("/docs/:section", func(c *gin.Context) {
		section := c.Param("section")
		for _, s := range models.DocsList {
			if s == section {
				c.HTML(http.StatusOK, section, nil)
				return
			}
		}
		c.String(http.StatusNotFound, "Not found")
	})
}

func AboutHandler(r *gin.Engine) {
	r.GET("/about/:section", func(c *gin.Context) {
		section := c.Param("section")
		for _, s := range models.AboutList {
			if s == section {
				c.HTML(http.StatusOK, section, nil)
				return
			}
		}
		c.String(http.StatusNotFound, "Not found")
	})
}

func ChatWidget(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.HTML(http.StatusOK, "chat-widget.tmpl", nil)
}
