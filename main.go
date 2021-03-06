// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt

package main

import (
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"translation.io/rest"
)

var mongoPath string
var mongoDb string

// The Router method routes requests to the appropriate Resource
func Router(path string) rest.Resource {
	if match, params := rest.MatchRoute("/collections/([a-z0-9]+)/strings/?([a-z0-9]+)?/?", path); match {
		if bson.IsObjectIdHex(params[1]) {
			cs := &CollectionStrings{}
			cs.Collection.Id = bson.ObjectIdHex(params[1])
			if len(params) > 1 && bson.IsObjectIdHex(params[2]) {
				cs.String.Id = bson.ObjectIdHex(params[2])
			}
			return cs
		} else {
			return &rest.NotFound{}
		}
	} else if match, params := rest.MatchRoute("/collections/([a-z0-9]+)", path); match {
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

func DocsHandler(w http.ResponseWriter, req *http.Request) {
	path := "docs" + req.URL.Path
	if path == "docs/" {
		path = "docs/index.html"
	}
	css, _ := regexp.Compile(".css$")
	if css.MatchString(path) {
		w.Header().Set("Content-Type", "text/css")
	}
	js, _ := regexp.Compile(".js$")
	if js.MatchString(path) {
		w.Header().Set("Content-Type", "text/javascript")
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Danger, Will Robinson! Danger! Page not found!", 404)
	}
	fmt.Fprintf(w, "%s", string(body))
}

// Handler parses all HTTP requests and retrieves responses
func APIHandler(w http.ResponseWriter, req *http.Request) {

	// Cast APIResponse interface
	var response rest.APIResponse

	// Parse Request
	req.ParseForm()
	values := &req.Form
	path := req.URL.Path

	// Redirect everything else to Docs
	regex, _ := regexp.Compile("/collections.*")
	if !regex.MatchString(path) {
		DocsHandler(w, req)
		return
	}

	// Default error when method is not found
	statusCode := 405
	response = rest.InvalidMethodError(&[]rest.Rel{
		rest.Rel{
			"POST":   "/collections",
			"Params": "name",
		},
		rest.Rel{"GET": "/collections/{CollectionId}"},
		rest.Rel{"PUT": "/collections/{CollectionId}",
			"Params": "name",
		},
		rest.Rel{"DELETE": "/collections/{CollectionId}"},
		rest.Rel{"POST": "/collections/{CollectionId}/strings",
			"Params": "string",
		},
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

	if os.Getenv("MONGOHQ_URL") != "" {
		mongoPath = os.Getenv("MONGOHQ_URL")
		mongoDb = "app13325198"
	} else {
		mongoPath = "127.0.0.1"
		mongoDb = "transio"
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	var port string
	if os.Getenv("PRODUCTION") == "true" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}
	fmt.Println("translation.io is running on http://localhost:" + port)
	http.HandleFunc("/", APIHandler)
	http.ListenAndServe(":"+port, nil)
}
