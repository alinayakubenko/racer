// TestAssgn project TestAssgn.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/wikirace", wikiRace)
	r.Run()
}
