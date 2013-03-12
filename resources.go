// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt

package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/url"
	"translation.io/rest"
)

type Collection struct {
	Id   bson.ObjectId "_id"
	Name string
}

// Implements APIResponse interface
func (c *Collection) ToJSON() string {
	return rest.ParseAPIResponse(c)
}

func (c *Collection) Get(v *url.Values) (int, rest.APIResponse) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("GET /collections - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	cs := []Collection{}
	if !c.Id.Valid() {
		// Id is not defined, return all collections
		err = C.Find(bson.M{}).All(&cs) // <-- potential performance problem!
		if err != nil {
			fmt.Println("GET /collections - Collection Query Error")
			return 500, rest.ServerError()
		}

		return 200, &rest.APISuccess{"collections": cs}
	} else {
		// return single collection
		err = C.FindId(c.Id).One(&c)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("GET /collections - Collection Query Error")
			return 500, rest.ServerError()
		}
		if err == mgo.ErrNotFound {
			return 404, rest.NotFoundError()
		}
		return 200, &rest.APISuccess{"collection": c}
	}
	return 404, rest.NotFoundError()
}

func (c *Collection) Post(v *url.Values) (int, rest.APIResponse) {
	name := v.Get("name")

	// Validate string
	if name == "" {
		return 422, &rest.APIError{
			Type:    "invalid-name",
			Message: "A non-empty name is required.",
			Code:    422,
			Param:   []string{"name"},
		}
	}

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("POST /collections - DB Connection Error")
		return 5001, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	// check if collection already exists
	err = C.Find(bson.M{"name": name}).One(&c)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("POST /collections - Collection Query Error")
		return 5002, rest.ServerError()
	}
	if c.Name != "" {
		return 200, &rest.APISuccess{"collection": c}
	}

	// insert new collection into DB
	c.Id = bson.NewObjectId()
	c.Name = name
	err = C.Insert(c)
	if err != nil {
		fmt.Println("POST /collections - Collection Insert Error")
		panic(err)
		return 5003, rest.ServerError()
	}

	return 200, &rest.APISuccess{"collection": c}

}

func (c *Collection) Put(v *url.Values) (int, rest.APIResponse) {
	// newName := v.Get("name")

	// session, err := mgo.Dial("127.0.0.1")
	// if err != nil {
	// 	fmt.Println("PUT /collections/" + c.Id.String() + " - DB Connection Error")
	// 	return 500, rest.ServerError()
	// }
	// C := session.DB("transio").C("collections")
	// defer session.Close()

	return 200, c
}

func (c *Collection) Delete(v *url.Values) (int, rest.APIResponse) {
	return 200, c
}

type String struct {
	Id         string
	String     string
	OriginLang string
}

// Implements APIResponse interface
func (s *String) ToJSON() string {
	return rest.ParseAPIResponse(s)
}

func (s *String) Get(v *url.Values) (int, rest.APIResponse) {

	// connect with db
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if s.Id == "" {
		return 404, rest.NotFoundError()
	}

	str := &String{
		Id:         "blabla",
		String:     "you know",
		OriginLang: "en-us",
	}
	return 200, str
}

/*
 * POST /String
 *
 * Params
 * - string:      An unique string that has to be translated. This POST method is 
 *                idempotent, so no worries when calling this method multiple times.
 * - origin_lang: The originating language for the string that is being inserted.
 *                Can only be "en-us" at the moment.
 */
func (s *String) Post(v *url.Values) (int, rest.APIResponse) {
	str := v.Get("string")
	lang := v.Get("origin_lang")

	// Validate string
	if str == "" {
		return 422, &rest.APIError{
			Type:    "invalid-string",
			Message: "A non-empty string is required.",
			Code:    422,
			Param:   []string{"string"},
		}
	}

	// Validate origin lang
	if lang != "en-us" {
		return 422, &rest.APIError{
			Type:    "invalid-origin-lang",
			Message: "Origin language can only be 'en-us' at the moment.",
			Code:    422,
			Param:   []string{"origin_lang"},
		}
	}

	// Search for similar string (makes "POST" idempotent)

	// Post new string
	type StringDoc struct {
		Id bson.ObjectId "_id"
	}

	return 200, s
}

func (s *String) Put(v *url.Values) (int, rest.APIResponse) {
	return 405, &rest.APIError{
		Type:    "invalid-method",
		Message: "This method is not allowed, use POST instead.",
		Code:    405,
		Param:   []string{},
	}
}

func (s *String) Delete(v *url.Values) (int, rest.APIResponse) {
	return 200, s
}

type Translations struct {
	Id          string
	StringId    string
	TargetLang  string
	Translation string
}

// Implements APIResponse interface
func (t *Translations) ToJSON() string {
	return rest.ParseAPIResponse(t)
}

func (t *Translations) Get(v *url.Values) (int, rest.APIResponse) {
	return 200, t
}

func (t *Translations) Post(v *url.Values) (int, rest.APIResponse) {
	return 200, t
}

func (t *Translations) Put(v *url.Values) (int, rest.APIResponse) {
	return 200, t
}

func (t *Translations) Delete(v *url.Values) (int, rest.APIResponse) {
	return 200, t
}
