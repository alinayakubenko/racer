package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Service contains the main logic for /wikirace endpoint

// Method race that is called from the controller.

var (
	lock = sync.RWMutex{}
)

const (
	TILE_TEXT  = "title"
	QUERY_TEXT = "query"
	PAGE_TEXT  = "pages"
)

func race(requestStartString string, requestEndString string) ResultResponse {
	var response ResultResponse
	var visited = make(map[string]bool)

	var res string
	pool := make(chan int, 2)
	res_chan := make(chan string)
	wg := &sync.WaitGroup{}

	response.Page = "null"
	response.Error = "null"

	ctx, cancel := context.WithCancel(context.Background())

	// input validation. Will look how else to do it
	if m, _ := regexp.MatchString("^[a-zA-Z, ]{1,50}$", requestStartString); !m {
		response.Error = "Validation failed for: " + requestStartString + ". String should contain no illigal characters and be no longer than 50 characters"
		log.Println("Validation Input String: ", requestStartString)
		return response
	}
	if m1, _ := regexp.MatchString("^[a-zA-Z, ]{1,50}$", requestEndString); !m1 {
		response.Error = "Validation failed for: " + requestEndString + ". String should contain no illigal characters and be no longer than 50 characters"
		log.Println("Validation Input String: ", requestEndString)
		return response
	}

	res = searchTitles(pool, res_chan, ctx, requestStartString, requestEndString, requestStartString, visited, 60)
	log.Println("Search started... ")

	res = <-res_chan
	response.Page = res
	log.Println("And the very final one is ", res)
	cancel()
	wg.Wait()
	close(res_chan)
	return response
}

// Method to perform searching using start and end strings from the request JSON
func searchTitles(pool chan int, res_chan chan string, ctx context.Context, currentTitleSting string, endString string, path string, visitedMap map[string]bool, max int) string {
	//Concurently writing the current title to the map of visited pages
	lock.Lock()
	defer lock.Unlock()
	visitedMap[currentTitleSting] = true

	if max == 0 {
		return ""
	} else {
		max -= 1
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
		//log.Println("before pool")
		pool <- 1
		//log.Println("inside pool")
		// Perform the HTTP call
		resp, _ := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		<-pool
		//log.Println("after pool")

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

			select {
			case <-ctx.Done():
				return ""
			default:
			}
			var title = pages[key].(map[string]interface{})[TILE_TEXT].(string)
			if !visitedMap[title] {
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
		}
		for _, item := range parr {
			select {
			case <-ctx.Done():
				return ""
			default:
			}
			//log.Println(item)
			go searchTitles(pool, res_chan, ctx, item, endString, path+" -> "+item, visitedMap, max)

		}
	}
	return ""
}
