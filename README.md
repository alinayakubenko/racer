## Project name  WikiRacer
## Overview
### Problem statement
   As a user I want to be able to to perform search in wikipedia  by providing a start and end article title. And receive the list of page titles linking start and end pages. 
   `  `

### Technologies 
   Go, Gin (RESTful)
## Technical approach
 ### Proposed design
` `
   This app was tested in the local environment but could also be installed on dedicated app server once stable. It wasn't tested for high accessibility or stability. 
   ` ` 
   
To let consumer perform the request, this application implemented as RESTful. The endpoint accepts JSON request (POST). The example request will be provided in this document.
`  `

   The start and end page titles then passed to the Service layer and used to request data from the MediaWiki API . The Service layer method then searches the end page title and returns the path (from Start to End page). Result then matches to the data model and returns to the consumer as JSON body. 
` `

   In case of failed HTTP request the partial result is returned with corresponding message.
   
### Implementation
` `
Diagram
![]({{site.baseurl}}/wikiRacerAPI_Diagram.png)
` `

I chose Go because it was new to me and because it's presumably very suitable for tasks that require concurrency. Gin was chosen for its simplicity and speed. 
`  `

Application is currently designed to have these main modules:

- WikiRacer. Main function that configures the endpoint and calls the controller
- Controller. Binds the request with the data model and passes the parameters to "race" method.
- Service. Contains logic, search function and composes the response. 
- Models. Describes data models for request, response and errors.
- Unit Tests
- Utils.

### Problems and things yet to implement:
` `
- Add retry functionality.  
- Add more Unit tests.   
- Find another the way how add validation to the data model or something similar to javax.
- Learn more how to create less coupled design in GO (packaging, layering and stuff).
- Learn more about links and pointers and if possible improve memory usage by the app.

### How to test this app on local environment. 
` `
1. The simple way to try it is to clone it from the GitHub repo to the local environment, install Go and Gin package.
2. in the console, from the app folder, execute command "go run WikiRacer" and confirm the app is running 
3. In Postman configure POST request to http://localhost:8081/wikirace
4. For header setup key: Content-Type with value: application/json
` `

Body example: 
` `

> {
	"startPage" : "Mike Tyson",
	"endPage"   : "Cannabis"
}

### More examples:
` `

![]({{site.baseurl}}/validation_error_message.png)

` `

![]({{site.baseurl}}/WikiRace_example1.png)

` `

![]({{site.baseurl}}/no_host_available_message.png)
