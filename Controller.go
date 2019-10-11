package main

import (
	"github.com/gin-gonic/gin"
)

func wikiRace(c *gin.Context) {
	var request StartEndRequest

	c.Header("Content-Type", "application/json")
	c.BindJSON(&request)
	httpResponse := race(request.StartPage, request.EndPage)
	if httpResponse.Error != "null" {
		c.JSON(400, gin.H{"error": httpResponse.Error, "responseTime": httpResponse.Time})
		return
	} else {
		c.JSON(200, gin.H{
			"path": httpResponse.Page, "responseTime": httpResponse.Time,
		})

	}

}
