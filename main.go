package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type Resource interface {
	Get(*http.Request) (int, string)
}

type Greeting struct {
	Id string
}

func (g *Greeting) Get(r *http.Request) (int, string) {
	return 200, "hello " + g.Id
}

type NotFound struct{}

func (n *NotFound) Get(r *http.Request) (int, string) {
	return 404, "not found"
}

func MatchRoute(r string, p string) (bool, [][]string) {
	regex, err := regexp.Compile(r)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	if regex.MatchString(p) {
		m := regex.FindAllStringSubmatch(p, -1)
		return true, m
	}

	return false, nil
}

func Router(p string) Resource {
	if ok, m := MatchRoute("/greetings/?([0-9]{0,12})?", p); ok {
		if len(m) == 0 || len(m[0]) <= 1 {
			return &Greeting{}
		}
		return &Greeting{
			Id: m[0][1],
		}
	}
	return &NotFound{}
}

func Handler(w http.ResponseWriter, r *http.Request) {

	statusCode := 405 // Method not allowed
	response := "Method not allowed"

	r.ParseForm()

	p := r.URL.Path

	switch r.Method {
	case "GET":
		statusCode, response = Router(p).Get(r)
		// case "POST":
		// 	statusCode, response = Router(p).Post(r)
		// case "PUT":
		// 	statusCode, response = Router(p).Put(r)
		// case "DELETE":
		// 	statusCode, response = Router(p).Delete(r)
	}

	if statusCode != 200 {
		http.Error(w, response, 404)
	} else {
		fmt.Fprintf(w, "%s", response)
	}
}

func main() {
	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
