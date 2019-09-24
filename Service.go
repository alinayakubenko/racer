package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	mutex = &sync.Mutex{}
)

func race(requestStartString string, requestEndString string) (ResultResponses, ErrorModel) {
	var response ResultResponses
	var errResponse ErrorModel
	//var res string
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	if m, _ := regexp.MatchString("^[a-zA-Z_ ]{1,50}$", requestStartString); !m {

		errResponse.Error = "Validation error!"

		log.Println("Validation response: ", errResponse.Error)
		log.Println("Validation Input String: ", requestStartString)
		return nil, errResponse
	}

	defer wg.Done()
	wg.Add(1)
	searchTitles(ctx, requestStartString, requestEndString, requestStartString, 20)
	fmt.Println("Terminating the application... completed ")

	cancel()
	wg.Wait()

	return response, errResponse
}

func searchTitles(ctx context.Context, currentTitleSting string, endString string, path string, maxDepth int) string {

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
		resp, err := client.Do(req)

		if err != nil {
			log.Println("Unable to perform request!", err)
			return ""
		}
		//log.Println(req.URL.String())

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		var query map[string]interface{}
		var pages map[string]interface{}

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
			var title = pages[key].(map[string]interface{})["title"].(string)
			pagesArr[index] = title
			index++

			if strings.ToLower(title) == strings.ToLower(endString) {
				log.Println("found the end page!!! ", title)
				path = path + " -> " + title
				log.Println("path ", path)
				ctx.Done()
				return path
			}
		}

		for _, item := range pagesArr {
			select {
			case <-ctx.Done():
				return ""
			default:
			}

			go searchTitles(ctx, item, endString, path+" -> "+item, maxDepth)
			log.Println(len(pagesArr))
		}

	}
	return ""
}
