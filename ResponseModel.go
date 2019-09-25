package main

type ResultResponse struct {
	Page chan string `json:"page"`
}

type ResultResponses []ResultResponse
