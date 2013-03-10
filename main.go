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
	if ok, vars := MatchRoute("/strings/?([a-zA-Z0-9]{0,12})?", path); ok {
		if len(vars[0]) <= 1 {
			return &Strings{}
		}
		return &Strings{
			Id: vars[1],
		}
	}
	return &NotFound{}
}

func Handler(w http.ResponseWriter, req *http.Request) {

	var resp APIResponse

	statusCode := 405 // Method not allowed
	resp = &APIError{
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
		statusCode, resp = Router(path).Get(values)
		// case "POST":
		// 	statusCode, body = Router(p).Post(req)
		// case "PUT":
		// 	statusCode, body = Router(p).Put(req)
		// case "DELETE":
		// 	statusCode, body = Router(p).Delete(req)
	}

	if statusCode != 200 {
		http.Error(w, resp.ToJSON(), 404)
	} else {
		fmt.Fprintf(w, "%s", resp.ToJSON())
	}
}

func main() {
	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
