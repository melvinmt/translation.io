package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type Resource interface {
	Get(*http.Request) (int, string)
	Post(*http.Request) (int, string)
	Put(*http.Request) (int, string)
	Delete(*http.Request) (int, string)
}

type Greeting struct {
	Id string
}

func (g *Greeting) Get(r *http.Request) (int, string) {
	return 200, "hello " + g.Id
}

func (g *Greeting) Post(r *http.Request) (int, string) {
	return 200, "hello " + g.Id
}

func (g *Greeting) Put(r *http.Request) (int, string) {
	return 200, "hello " + g.Id
}

func (g *Greeting) Delete(r *http.Request) (int, string) {
	return 200, "hello " + g.Id
}

type NotFound struct{}

func (n *NotFound) Get(r *http.Request) (int, string) {
	return 404, "not found"
}

func (n *NotFound) Post(r *http.Request) (int, string) {
	return 404, "not found"
}

func (n *NotFound) Put(r *http.Request) (int, string) {
	return 404, "not found"
}

func (n *NotFound) Delete(r *http.Request) (int, string) {
	return 404, "not found"
}

func MatchRoute(r string, p string) (bool, []string) {
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

func Router(path string) Resource {
	if ok, vars := MatchRoute("/greetings/?([a-zA-Z0-9]{0,12})?", path); ok {
		if len(vars[0]) <= 1 {
			return &Greeting{}
		}
		return &Greeting{
			Id: vars[1],
		}
	}
	return &NotFound{}
}

func Handler(w http.ResponseWriter, req *http.Request) {

	statusCode := 405 // Method not allowed
	body := "Method not allowed"

	req.ParseForm()

	path := req.URL.Path

	switch req.Method {
	case "GET":
		statusCode, body = Router(path).Get(req)
		// case "POST":
		// 	statusCode, body = Router(p).Post(req)
		// case "PUT":
		// 	statusCode, body = Router(p).Put(req)
		// case "DELETE":
		// 	statusCode, body = Router(p).Delete(req)
	}

	if statusCode != 200 {
		http.Error(w, body, 404)
	} else {
		fmt.Fprintf(w, "%s", body)
	}
}

func main() {
	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
