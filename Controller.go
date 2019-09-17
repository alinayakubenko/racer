// TestAssgn project TestAssgn.go
package main

import (
	"github.com/gin-gonic/gin"
)

func wikiRace(c *gin.Context) {
	var request StartEndRequest

	c.BindJSON(&request)
	c.JSON(200, gin.H{
		"message end": race(request.StartPage, request.EndPage),
	})
}
