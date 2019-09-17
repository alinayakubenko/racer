package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"
)

func race(requestStartString string, requestEndString string) ResultResponses {

	var response ResultResponses

	fmt.Println("test")
	go searchTitles(requestStartString, requestEndString)
	fmt.Println("Terminating the application...")

	return response
}

func searchTitles(currentTitleSting string, endString string) {

	client := &http.Client{}
	client.Timeout = time.Second * 15

	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		q := req.URL.Query()
		q.Add("action", "query")
		q.Add("format", "json")
		q.Add("titles", currentTitleSting)
		q.Add("generator", "links")
		q.Add("redirects", "")
		q.Add("gpllimit", "500")
		q.Add("prop", "info")
		q.Add("inprop", "url")
		q.Add("alnamespace", "0")

		req.URL.RawQuery = q.Encode()
		resp, _ := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(req.URL.String())
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		log.Println(result)
	}
}
