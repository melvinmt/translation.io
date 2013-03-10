package main

import (
	"encoding/json"
	"net/url"
)

// All API Responses (including entities) must implement this interface
type APIResponse interface {
	ToJSON() string
}

// Strict Error APIResponse
type APIError struct {
	Type    string
	Message string
	Code    int
	Param   []string
}

func (e APIError) ToJSON() string {
	return ParseAPIResponse(e)
}

// Generic Success APIResponse (e.g. when no entities are involved)
type APISuccess map[string]interface{}

func (s APISuccess) ToJSON() string {
	return ParseAPIResponse(s)
}

// Parses APIResponse interfaces
func ParseAPIResponse(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		// Oops, response could not be parsed:
		b2, _ := json.Marshal(ServerError())
		return string(b2)
	}
	return string(b)
}

// Default Server Error (when unexpected things go wrong)
func ServerError() *APIError {
	return &APIError{
		Type:    "server-error",
		Message: "This response could not be processed at this time.",
		Code:    500,
		Param:   []string{},
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

// Default Not Found Error (when things can't be found)
func NotFoundError() *APIError {
	return &APIError{
		Type:    "not-found",
		Message: "This resource was not found.",
		Code:    404,
		Param:   []string{},
	}
}
