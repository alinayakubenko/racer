package main

import (
	"testing"
)

func TestInputValidation(t *testing.T) {
	inputString1 := "Mike Tyson&#"
	inputString2 := "Mike Tyson"
	response1 := inputValidation(inputString1)
	response2 := inputValidation(inputString2)
	switch {
	case response1.Error != "Validation failed for: "+inputString1+". String should contain no illigal characters and be no longer than 50 characters":
		t.Error("Test failed. Input validation: Error value should contain message")
	case response2.Error != "":
		t.Error("Test failed. Input validation: Error value should be empty")
	}
}
