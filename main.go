// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt

package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/http"
	"runtime"
	"translation.io/rest"
)

// The Router method routes requests to the appropriate Resource
func Router(path string) rest.Resource {
	if match, params := rest.MatchRoute("/collections/([a-z0-9]+)", path); match {
		if bson.IsObjectIdHex(params[1]) {
			return &Collection{
				Id: bson.ObjectIdHex(params[1]),
			}
		} else {
			return &rest.NotFound{}
		}
	} else if match, _ := rest.MatchRoute("/collections/?", path); match {
		return &Collection{}
	}
	return &rest.NotFound{}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s", "hi!")
}

// Handler parses all HTTP requests and retrieves responses
func APIHandler(w http.ResponseWriter, req *http.Request) {

	// Cast APIResponse interface
	var response rest.APIResponse

	// Parse Request
	req.ParseForm()
	values := &req.Form
	path := req.URL.Path

	if path == "/" {
		RootHandler(w, req)
		return
	}

	// Default error when method is not found
	statusCode := 405
	response = rest.InvalidMethodError(&[]rest.Rel{
		rest.Rel{"POST": "/collections"},
		rest.Rel{"GET": "/collections/{CollectionId}"},
		rest.Rel{"PUT": "/collections/{CollectionId}"},
		rest.Rel{"DELETE": "/collections/{CollectionId}"},
		rest.Rel{"POST": "/collections/{CollectionId}/strings"},
		rest.Rel{"DELETE": "/collections/{CollectionId}/strings/{StringId}"},
	})

	// Retrieve response on allowed methods
	switch req.Method {
	case "GET":
		statusCode, response = Router(path).Get(values)
	case "POST":
		statusCode, response = Router(path).Post(values)
	case "PUT":
		statusCode, response = Router(path).Put(values)
	case "DELETE":
		statusCode, response = Router(path).Delete(values)
	}

	// Return response
	if statusCode != 200 {
		http.Error(w, response.ToJSON(), statusCode)
	} else {
		fmt.Fprintf(w, "%s", response.ToJSON())
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", APIHandler)
	http.ListenAndServe(":8080", nil)
}
