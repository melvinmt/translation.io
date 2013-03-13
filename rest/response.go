// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software 
// without restriction, including without limitation the rights to use, copy, modify, 
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit 
// persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR 
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR 
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
// DEALINGS IN THE SOFTWARE.

package rest

import (
	"encoding/json"
	"net/url"
)

// All API Responses (including entities) must implement this interface
type APIResponse interface {
	ToJSON() string
}

type Rel map[string]string

type ErrorMsg struct {
	Type    string
	Message string
	Code    int
	Param   []string
	Allowed *[]Rel
}

// Strict Error APIResponse
type APIError struct {
	Error ErrorMsg
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
		Error: ErrorMsg{
			Type:    "server-error",
			Message: "This response could not be processed at this time.",
			Code:    500,
			Param:   []string{},
		},
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
		Error: ErrorMsg{
			Type:    "not-found",
			Message: "This resource was not found.",
			Code:    404,
			Param:   []string{},
		},
	}
}

// Default Not Found Error (when things can't be found)
func InvalidMethodError(rel *[]Rel) *APIError {
	return &APIError{
		Error: ErrorMsg{
			Type:    "invalid-method",
			Message: "Sorry, this method is not allowed.",
			Code:    405,
			Param:   []string{},
			Allowed: rel,
		},
	}
}
