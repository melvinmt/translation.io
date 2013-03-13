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
	Id      bson.ObjectId "_id"
	Name    string
	Strings []String
}

type String struct {
	Id           bson.ObjectId "_id"
	String       string
	Translations []Translation
}

type Translation struct {
	Id          bson.ObjectId "_id"
	StringId    bson.ObjectId
	Language    string
	Translation string
}

// Implements APIResponse interface
func (c *Collection) ToJSON() string {
	return rest.ParseAPIResponse(c)
}

func (c *Collection) Get(v *url.Values) (int, rest.APIResponse) {

	// Initialize DB
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	if !c.Id.Valid() {

		// GET on /collections (without an ID) is not allowed
		return 405, rest.InvalidMethodError(&[]rest.Rel{
			rest.Rel{"POST": "/collections"},
		})

	} else {

		// Return single collection
		err = C.FindId(c.Id).One(&c)
		if err != nil && err != mgo.ErrNotFound {
			return 500, rest.ServerError()
		}
		if err == mgo.ErrNotFound {
			return 404, rest.NotFoundError()
		}
		return 200, &rest.APISuccess{
			"Collection": c,
			"Next": &[]rest.Rel{
				rest.Rel{"GET": "/collections/" + c.Id.String()},
				rest.Rel{"PUT": "/collections/" + c.Id.String()},
				rest.Rel{"DELETE": "/collections/" + c.Id.String()},
			},
		}
	}

	// Return Not Found Error
	return 404, rest.NotFoundError()
}

func (c *Collection) Post(v *url.Values) (int, rest.APIResponse) {

	// Validate Name
	name := v.Get("name")
	if name == "" {
		return 422, &rest.APIError{
			Type:    "invalid-name",
			Message: "A non-empty name is required.",
			Code:    422,
			Param:   []string{"name"},
		}
	}

	// Initialize DB
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	// Check if Collection already exists
	err = C.Find(bson.M{"Name": name}).One(&c)
	if err != nil && err != mgo.ErrNotFound {
		return 500, rest.ServerError()
	}
	if c.Name != "" {
		return 200, &rest.APISuccess{"Collection": c}
	}

	// Insert new Collection into DB
	c.Id = bson.NewObjectId()
	c.Name = name
	err = C.Insert(c)
	if err != nil {
		return 5003, rest.ServerError()
	}

	// Return Collection
	return 200, &rest.APISuccess{"Collection": c}
}

func (c *Collection) Put(v *url.Values) (int, rest.APIResponse) {

	// Check Id
	if !c.Id.Valid() {
		return 404, rest.NotFoundError()
	}

	// Initialize DB
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("PUT /collections/" + c.Id.String() + " - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	// Validate Name
	newName := v.Get("name")
	if newName == "" {
		return 422, &rest.APIError{
			Type:    "invalid-name",
			Message: "A non-empty name is required.",
			Code:    422,
			Param:   []string{"name"},
		}
	}

	// Update Collection
	c.Name = newName
	err = C.UpdateId(c.Id, c)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("POST /collections - Collection Query Error")
		return 500, rest.ServerError()
	}
	if err == mgo.ErrNotFound {
		return 404, rest.NotFoundError()
	}

	// Return Collection
	return 200, &rest.APISuccess{"Collection": c}
}

func (c *Collection) Delete(v *url.Values) (int, rest.APIResponse) {

	// Initialize DB
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("DELETE /collections/" + c.Id.String() + " - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	// Remove Collection
	err = C.RemoveId(c.Id)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("DELETE /collections/" + c.Id.String() + " - Delete Collection Error")
		return 500, rest.ServerError()
	}
	return 200, &rest.APISuccess{"Success": true}
}

type CollectionStrings struct {
	Collection Collection
	String     String
}

// Implements APIResponse interface
func (c *CollectionStrings) ToJSON() string {
	return rest.ParseAPIResponse(c)
}

func (c *CollectionStrings) Get(v *url.Values) (int, rest.APIResponse) {
	return 405, rest.InvalidMethodError(&[]rest.Rel{
		rest.Rel{"POST": "/collections/" + c.Collection.Id.String() + "/strings"},
		rest.Rel{"DELETE": "/collections/" + c.Collection.Id.String() + "/strings/{StringId}"},
	})
}

func (c *CollectionStrings) Post(v *url.Values) (int, rest.APIResponse) {
	// Check Id
	if !c.Collection.Id.Valid() {
		return 404, rest.NotFoundError()
	}

	// Initialize DB
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Println("POST /collections/" + c.Collection.Id.String() + "/strings - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB("transio").C("collections")
	defer session.Close()

	// Find Collection
	err = C.FindId(c.Collection.Id).One(&c.Collection)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("POST /collections/" + c.Collection.Id.String() + "/strings - Collection Query Error")
		return 500, rest.ServerError()
	}
	if err == mgo.ErrNotFound {
		return 404, rest.NotFoundError()
	}

	// Validate string
	str := v.Get("string")
	if str == "" {
		return 422, &rest.APIError{
			Type:    "invalid-string",
			Message: "A non-empty string is required.",
			Code:    422,
			Param:   []string{"string"},
		}
	}

	// Init new String struct
	s := String{}

	// Search for similar string in Collection Strings array (makes "POST" idempotent)
	if len(c.Collection.Strings) > 0 {
		for _, Item := range c.Collection.Strings {
			if Item.String == str {
				s.Id = Item.Id
				s.String = str
				break
			}
		}
	}

	// Search for similar string in DB
	if s.Id == "" {
		S := session.DB("transio").C("strings")

		err = S.Find(bson.M{"String": str}).One(&s)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("POST /collections/" + c.Collection.Id.String() + "/strings - Strings Query Error")
			return 500, rest.ServerError()
		}
		if err == mgo.ErrNotFound {
			return 404, rest.NotFoundError()
		}

		/* TODO: insert new string into strings DB */

		/* TODO: translate string into x languages! */
	}

	// Add String to Collection and Update Collection
	c.Collection.Strings = append(c.Collection.Strings, s)
	err = C.Update(c.Collection.Id, c.Collection)
	if err != nil {
		fmt.Println("POST /collections/" + c.Collection.Id.String() + "/strings - Update Collection Error")
		return 500, rest.ServerError()
	}

	return 200, &rest.APISuccess{"Success": true}
}

func (c *CollectionStrings) Put(v *url.Values) (int, rest.APIResponse) {
	return 405, rest.InvalidMethodError(&[]rest.Rel{
		rest.Rel{"POST": "/collections/" + c.Collection.Id.String() + "/strings"},
		rest.Rel{"DELETE": "/collections/" + c.Collection.Id.String() + "/strings/{StringId}"},
	})
}

func (c *CollectionStrings) Delete(v *url.Values) (int, rest.APIResponse) {
	/* TODO: delete string from collection! */
	return 200, c
}
