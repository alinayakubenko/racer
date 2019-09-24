package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

/*
TODO
Fix the issue with timeout request,
add tracking of visited pages,
map and return the result,
add unit tests,
updated valiidation,
see what's elese to improve
*/

// Service contains the main logic for /wikirace endpoint
var (
	mutex = &sync.Mutex{}
)

// Method race that is called from the controller.
func race(requestStartString string, requestEndString string) (ResultResponses, ErrorModel) {
	var response ResultResponses
	var errResponse ErrorModel
	wg := &sync.WaitGroup{}
	ch := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	if m, _ := regexp.MatchString("^[a-zA-Z_, ]{1,50}$", requestStartString); !m {

		errResponse.Error = "Validation error!"

		log.Println("Validation response: ", errResponse.Error)
		log.Println("Validation Input String: ", requestStartString)
		return nil, errResponse
	}

	defer wg.Done()
	wg.Add(1)
	searchTitles(ctx, ch, requestStartString, requestEndString, requestStartString, 20)
	fmt.Println("Start waiting for result from the channel: ")

	log.Println("Channel log result: ", <-ch)

	cancel()
	wg.Wait()
	//close(ch)
	return response, errResponse
}

// Method to perform searching using start and end strings from the request JSON
func searchTitles(ctx context.Context, ch chan<- string, currentTitleSting string, endString string, path string, maxDepth int) string {

	if maxDepth == 0 {
		return ""
	} else {
		maxDepth -= 1
	}
	select {
	case <-ctx.Done():
		return ""
	default:
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		// Forming the parametrized request to call the MediaWiki API
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
		// Perform the HTTP call
		resp, err := client.Do(req)

		if err != nil {
			log.Println("Unable to perform request!", err)
			return ""
		}
		//log.Println(req.URL.String())
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		//Creating the variable to hold the map from the result
		var query map[string]interface{}
		var pages map[string]interface{}

		//Checking if the specific field exists before writing it
		query, ok := result["query"].(map[string]interface{})
		if ok {
			pages, ok = query["pages"].(map[string]interface{})
			if !ok {
				return ""
			}
		} else {
			return ""
		}

		pagesArr := make([]string, len(pages))
		var index = 0
		for key := range pages {
			select {
			case <-ctx.Done():
				return ""
			default:
			}
			var title = pages[key].(map[string]interface{})["title"].(string)
			pagesArr[index] = title
			index++

			if strings.ToLower(title) == strings.ToLower(endString) {
				log.Println("found the end page!!! ", title)
				path = path + " -> " + title
				log.Println("path ", path)
				ctx.Done()
				ch <- path
				return ""
			}
		}

		for _, item := range pagesArr {
			select {
			case <-ctx.Done():
				log.Println("Context done event:  ", currentTitleSting)
				return ""
			default:
			}
			go searchTitles(ctx, ch, item, endString, path+" -> "+item, maxDepth)
			log.Println(len(pagesArr))
		}

	}
	return ""
}
