package main

import (
	"fmt"
	"net/http"
)

type Resource interface {
	Get() (int, string)
}

type Greeting struct{}

func (g Greeting) Get() (int, string) {
	return 200, "hello"
}

type NotFound struct{}

func (n NotFound) Get() (int, string) {
	return 404, "not found"
}

func route(p string) Resource {
	switch p {
	case "/greeting":
		return Greeting{}
	}
	return NotFound{}
}

func handler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	// Retrieve method
	method := r.Method
	fmt.Println("method:", method)

	// Parse path
	path := r.URL.Path
	fmt.Println("path:", path)

	// Parse params from Query string
	r.ParseForm()
	params := r.Form // type: url.Values

	fmt.Println("params:", params)

	var statusCode int
	var resp string

	switch method {
	case "GET":
		statusCode, resp = route(path).Get()
	}

	switch statusCode {
	case 404:
		http.Error(w, resp, 404)
	}
	if statusCode == 200 {
		fmt.Fprintf(w, "%s", resp)
	}
}

func main() {
	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
