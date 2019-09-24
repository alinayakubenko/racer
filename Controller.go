package main

import (
	"github.com/gin-gonic/gin"
)

func wikiRace(c *gin.Context) {
	//Contrloller uses gin to handle http request and maps the request and (eventially the response) to the data models
	var request StartEndRequest

	c.BindJSON(&request)
	resp, err := race(request.StartPage, request.EndPage)
	if err.Error == "" {
		c.JSON(200, gin.H{
			"message end": resp,
		})
	}
}
