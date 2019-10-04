package main

import (
	"github.com/gin-gonic/gin"
)

func wikiRace(c *gin.Context) {
	var request StartEndRequest

	c.BindJSON(&request)
	if response.Error != "" {
		c.JSON(400, gin.H{"error": race(request.StartPage, request.EndPage)})
		return
	}
	c.JSON(200, gin.H{
		"message end": race(request.StartPage, request.EndPage),
	})

}
