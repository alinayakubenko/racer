package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//Entry point, configuring an endpoint and callin the controller method
	r := gin.Default()
	r.POST("/wikirace", wikiRace)
	r.Run()
}
