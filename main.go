// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt

package main

import (
	"fmt"
	"net/http"
	"runtime"
	"translation.io/rest"
)

// The Router method routes requests to the appropriate Resource
func Router(path string) rest.Resource {
	if match, _ := rest.MatchRoute("/strings/?", path); match {
		return &String{}
	} else if match, params := rest.MatchRoute("/strings/([a-zA-Z0-9]{0,12})", path); match {
		return &String{
			Id: params[1],
		}
	} else if match, params := rest.MatchRoute("/strings/([a-zA-Z0-9]{0,12})/translations/?", path); match {
		return &Translations{
			StringId: params[1],
		}
	} else if match, params := rest.MatchRoute("/strings/([0-9]{0,12})/translations/([a-z]{2}(-[a-z]{2})?)", path); match {
		return &Translations{
			StringId:   params[1],
			TargetLang: params[2],
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
	response = &rest.APIError{
		Type:    "invalid-method",
		Message: "Sorry, this method is not allowed.",
		Code:    405,
		Param:   []string{},
	}

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
		http.Error(w, response.ToJSON(), 404)
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
