package main

type ResultResponse struct {
	Page  string `json:"path"`
	Error string `json:"validationError"`
}

type ResultResponses []ResultResponse
