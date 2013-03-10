package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type Resource interface {
	Get(*url.Values) (int, APIResponse)
	Post(*url.Values) (int, APIResponse)
	Put(*url.Values) (int, APIResponse)
	Delete(*url.Values) (int, APIResponse)
}

type Greeting struct {
	Id string
}

func (g *Greeting) Get(v *url.Values) (int, APIResponse) {
	return 200, &APISuccess{
		"response": "hello" + g.Id,
	}
}

func (g *Greeting) Post(v *url.Values) (int, APIResponse) {
	return 200, &APISuccess{
		"response": "hello" + g.Id,
	}
}

func (g *Greeting) Put(v *url.Values) (int, APIResponse) {
	return 200, &APISuccess{
		"response": "hello" + g.Id,
	}
}

func (g *Greeting) Delete(v *url.Values) (int, APIResponse) {
	return 200, &APISuccess{
		"response": "hello" + g.Id,
	}
}

type NotFound struct{}

func (n *NotFound) Get(v *url.Values) (int, APIResponse) {
	return 404, NotFoundError()
}

func (n *NotFound) Post(v *url.Values) (int, APIResponse) {
	return 404, NotFoundError()
}

func (n *NotFound) Put(v *url.Values) (int, APIResponse) {
	return 404, NotFoundError()
}

func (n *NotFound) Delete(v *url.Values) (int, APIResponse) {
	return 404, NotFoundError()
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

	var response APIResponse

	statusCode := 405 // Method not allowed

	response = APIError{
		Type:    "invalid-method",
		Message: "Sorry, this method is not allowed.",
		Code:    405,
		Param:   []string{},
	}

	req.ParseForm()
	values := &req.Form

	path := req.URL.Path

	switch req.Method {
	case "GET":
		statusCode, response = Router(path).Get(values)
		// case "POST":
		// 	statusCode, body = Router(p).Post(req)
		// case "PUT":
		// 	statusCode, body = Router(p).Put(req)
		// case "DELETE":
		// 	statusCode, body = Router(p).Delete(req)
	}

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
