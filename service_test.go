package main

import (
	"strings"
	"testing"
)

func TestRaceSamePage(t *testing.T) {
	var StartPage = "David Tennant"
	var EndPage = StartPage
	var result = (race(StartPage, EndPage)).Page

	if !strings.HasSuffix(result, "-> "+StartPage) {
		t.Errorf("Wrong result, expected ends with \"%s\", got \"%s\"", "-> "+StartPage, result)
	}
}
