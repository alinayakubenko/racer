package main

type ResultResponse struct {
	Page  string `json:"path"`
	Error string `json:"error"`
	Time  string `json:"requestTime"`
}

type ResultResponses []ResultResponse
