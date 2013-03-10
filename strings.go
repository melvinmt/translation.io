package main

import (
	"net/url"
)

type Strings struct {
	Id          string
	String      string
	Origin_Lang string
}

// Implements APIResponse interface
func (s *Strings) ToJSON() string {
	return ParseAPIResponse(s)
}

func (s *Strings) Get(v *url.Values) (int, APIResponse) {

	return 200, s
}

/*
 * POST /strings
 *
 * Params
 * - string:      An unique string that has to be translated. This POST method is 
 *                idempotent, so no worries when calling this method multiple times.
 * - origin_lang: The originating language for the string that is being inserted.
 *                Can only be "en-us" at the moment.
 */
func (s *Strings) Post(v *url.Values) (int, APIResponse) {
	str := v.Get("string")
	lang := v.Get("origin_lang")

	// Validate string
	if str == "" {
		return 422, &APIError{
			Type:    "invalid-string",
			Message: "A non-empty string is required.",
			Code:    422,
			Param:   []string{"string"},
		}
	}

	// Validate origin lang
	if lang != "en-us" {
		return 422, &APIError{
			Type:    "invalid-origin-lang",
			Message: "Origin language can only be 'en-us' at the moment.",
			Code:    422,
			Param:   []string{"origin_lang"},
		}
	}

	// Search for similar string (makes "POST" idempotent)

	// Post new string

	return 200, s
}

func (s *Strings) Put(v *url.Values) (int, APIResponse) {
	return 405, &APIError{
		Type:    "invalid-method",
		Message: "This method is not allowed, use POST instead.",
		Code:    405,
		Param:   []string{},
	}
}

func (s *Strings) Delete(v *url.Values) (int, APIResponse) {
	return 200, s
}
