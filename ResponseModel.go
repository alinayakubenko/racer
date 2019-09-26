package main

type ResultResponse struct {
	Page  string `json:"page"`
	Error string `json:"error"`
}

type ResultResponses []ResultResponse
