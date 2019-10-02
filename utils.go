package main

import (
	"log"
	"regexp"
)

func inputValidation(inputString string) ResultResponse {
	// input validation. Will need to look how else this could be done
	response.Error = ""
	if m, _ := regexp.MatchString("^[a-zA-Z, ]{1,50}$", inputString); !m {
		response.Error = "Validation failed for: " + inputString + ". String should contain no illigal characters and be no longer than 50 characters"
		log.Println("Validation Input String: ", inputString)
		return response
	}
	return response
}
