package main

import (
	"fmt"
	"github.com/melvinmt/rest"
	"net/http"
)

// The Router method routes requests to the appropriate Resource
func Router(path string) rest.Resource {
	if ok, params := rest.MatchRoute("/api/v1/strings/?([a-zA-Z0-9]{0,12})?", path); ok {
		if len(params[0]) <= 1 {
			return &Strings{}
		}
		return &Strings{
			Id: params[1],
		}
	}
	return &rest.NotFound{}
}

func main() {
	api := &rest.API{
		Router:  Router, // extends rest.Router()
		Handler: rest.Handler,
	}

	fmt.Println("translation.io is running on http://localhost:8080")
	http.HandleFunc("/api/v1", api.Handler)
	http.ListenAndServe(":8080", nil)
}
