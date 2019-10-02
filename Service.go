package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Service contains the main logic for /wikirace endpoint

// Method race that is called from the controller.

var (
	lock     = sync.RWMutex{}
	response ResultResponse
)

const (
	TILE_TEXT  = "title"
	QUERY_TEXT = "query"
	PAGE_TEXT  = "pages"
)

func race(requestStartString string, requestEndString string) ResultResponse {
	var response ResultResponse
	//Map to record visited pages
	var visited = make(map[string]bool)

	var res string = ""
	pool := make(chan int, 10)
	res_chan := make(chan string)

	//Setting default values to the response fields
	response.Page = "null"
	response.Error = "null"

	ctx, cancel := context.WithCancel(context.Background())

	//Calling input validation method to verify the request body
	if inputValidation(requestStartString).Error != "" {
		return inputValidation(requestStartString)
	} else if inputValidation(requestEndString).Error != "" {
		return inputValidation(requestEndString)
	}

	res = searchTitles(pool, res_chan, ctx, requestStartString, requestEndString, requestStartString, visited, 30)
	log.Println("Search started... ")

	res = <-res_chan
	if res != "" {
		response.Page = res
	} else {
		response.Error = "Page not found."
	}

	log.Println("And the very final one is ", res)
	cancel()
	//close(res_chan)
	return response
}

// Method to perform searching using start and end strings from the request JSON.
func searchTitles(pool chan int, res_chan chan string, ctx context.Context, currentTitleSting string, endString string, path string, visitedMap map[string]bool, max int) string {
	defer recoverRequest(ctx, res_chan, path)
	//Concurently writing the current title to the map of visited pages
	lock.Lock()
	visitedMap[currentTitleSting] = true
	lock.Unlock()

	if max == 0 {
		return ""
	} else {
		max -= 1
	}

	listenToCancel(ctx)
	// HTTP client setup
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}

	// Building the query to request data from MediaWiki API
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
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
		pool <- 1
		// Perform the HTTP call
		// Need to add some retry functionality for to handle failed request
		listenToCancel(ctx)
		resp, _ := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		<-pool
		// Mapping data
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		//Creating the variable to hold the map from the result
		var query map[string]interface{}
		var pages map[string]interface{}

		//Checking if the specific field exists before writing it
		query, ok := result[QUERY_TEXT].(map[string]interface{})
		if ok {
			pages, ok = query[PAGE_TEXT].(map[string]interface{})
			if !ok {
				return ""
			}
		} else {
			return ""
		}

		parr := make([]string, len(pages))
		var index = 0
		for key := range pages {

			listenToCancel(ctx)
			var title = pages[key].(map[string]interface{})[TILE_TEXT].(string)
			parr[index] = title
			index++

			if strings.ToLower(title) == strings.ToLower(endString) {
				log.Println("found the end page!!! ", title)
				path = path + " -> " + title
				log.Println("path ", path)
				res_chan <- path
				return path
			}

		}
		for _, item := range parr {
			listenToCancel(ctx)
			//log.Println(item)
			lock.RLock()
			if !visitedMap[item] {
				go searchTitles(pool, res_chan, ctx, item, endString, path+" -> "+item, visitedMap, max)
			}
			lock.RUnlock()
		}
	}
	return ""
}

// Method for recovering form the panic in case of failed http request. It then returns the current path and a message as a result.
func recoverRequest(ctx context.Context, res_chan chan string, path string) {
	if err := recover(); err != nil {
		res_chan <- " Network error. Unfinished search! Current path: " + path
		ctx.Done()
		log.Println("recovered from ", err)

	}
}

func listenToCancel(ctx context.Context) string {
	select {
	case <-ctx.Done():
		return ""
	default:
	}
	return ""
}
