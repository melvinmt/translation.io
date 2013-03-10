package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

// All Resources need to implement this interface
type Resource interface {
	Get(*url.Values) (int, APIResponse)
	Post(*url.Values) (int, APIResponse)
	Put(*url.Values) (int, APIResponse)
	Delete(*url.Values) (int, APIResponse)
}

// matchRoute is a helper method for the Router to interprete paths and to parse optional params
func matchRoute(r string, p string) (bool, []string) {
	regex, err := regexp.Compile(r)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	if regex.MatchString(p) {
		m := regex.FindStringSubmatch(p)
		return true, m
	}

	return false, nil
}

// The Router method routes requests to the appropriate Resource
func Router(path string) Resource {
	if ok, params := matchRoute("/strings/?([a-zA-Z0-9]{0,12})?", path); ok {
		if len(params[0]) <= 1 {
			return &Strings{}
		}
		return &Strings{
			Id: params[1],
		}
	}
	return &NotFound{}
}

// Handler parses all HTTP requests and retrieves responses
func Handler(w http.ResponseWriter, req *http.Request) {

	// Cast APIResponse interface
	var response APIResponse

	// Parse Request
	req.ParseForm()
	values := &req.Form
	path := req.URL.Path

	// Default error when method is not found
	statusCode := 405
	response = &APIError{
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
	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
