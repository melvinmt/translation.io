// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"os"
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
	Translations map[string]string
}

type Translation struct {
	Language    string
	Translation string
}

// Implements APIResponse interface
func (c *Collection) ToJSON() string {
	return rest.ParseAPIResponse(c)
}

func (c *Collection) Get(v *url.Values) (int, rest.APIResponse) {

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		return 500, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	if !c.Id.Valid() {

		// GET on /collections (without an ID) is not allowed
		return 405, rest.InvalidMethodError(&[]rest.Rel{
			rest.Rel{"POST": "/collections",
				"Params": "name",
			},
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
				rest.Rel{
					"PUT":    "/collections/" + c.Id.Hex(),
					"Params": "name",
				},
				rest.Rel{"DELETE": "/collections/" + c.Id.Hex()},
				rest.Rel{
					"POST":   "/collections/" + c.Id.Hex() + "/strings",
					"Params": "string",
				},
				rest.Rel{"DELETE": "/collections/" + c.Id.Hex() + "/strings/{StringId}"},
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
			rest.Error: rest.ErrorMsg{
				Type:    "invalid-name",
				Message: "A non-empty name is required.",
				Code:    422,
				Param:   []string{"name"},
			},
		}
	}

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		return 500, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	err = C.Find(bson.M{"name": name}).One(&c)
	if err != nil && err != mgo.ErrNotFound {
		return 500, rest.ServerError()
	}

	if c.Name == "" {
		// Insert new Collection into DB
		c.Id = bson.NewObjectId()
		c.Name = name
		err = C.Insert(c)
		if err != nil {
			return 5003, rest.ServerError()
		}
	}

	// Return Collection
	return 200, &rest.APISuccess{
		"Collection": c,
		"Next": &[]rest.Rel{
			rest.Rel{"GET": "/collections/" + c.Id.Hex()},
			rest.Rel{"PUT": "/collections/" + c.Id.Hex(),
				"Params": "name",
			},
			rest.Rel{"DELETE": "/collections/" + c.Id.Hex()},
			rest.Rel{"POST": "/collections/" + c.Id.Hex() + "/strings",
				"Params": "string",
			},
			rest.Rel{"DELETE": "/collections/" + c.Id.Hex() + "/strings/{StringId}"},
		},
	}
}

func (c *Collection) Put(v *url.Values) (int, rest.APIResponse) {

	// Check Id
	if !c.Id.Valid() {
		return 404, rest.NotFoundError()
	}

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		fmt.Println("PUT /collections/" + c.Id.Hex() + " - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	// Validate Name
	newName := v.Get("name")
	if newName == "" {
		return 422, &rest.APIError{
			rest.Error: rest.ErrorMsg{
				Type:    "invalid-name",
				Message: "A non-empty name is required.",
				Code:    422,
				Param:   []string{"name"},
			},
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
	return 200, &rest.APISuccess{
		"Collection": c,
		"Next": &[]rest.Rel{
			rest.Rel{"GET": "/collections/" + c.Id.Hex()},
			rest.Rel{"DELETE": "/collections/" + c.Id.Hex()},
			rest.Rel{"POST": "/collections/" + c.Id.Hex() + "/strings",
				"Params": "string",
			},
			rest.Rel{"DELETE": "/collections/" + c.Id.Hex() + "/strings/{StringId}"},
		},
	}
}

func (c *Collection) Delete(v *url.Values) (int, rest.APIResponse) {

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		fmt.Println("DELETE /collections/" + c.Id.Hex() + " - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	// Remove Collection
	err = C.RemoveId(c.Id)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("DELETE /collections/" + c.Id.Hex() + " - Delete Collection Error")
		return 500, rest.ServerError()
	}
	return 200, &rest.APISuccess{
		"Success": true,
		"Next": &[]rest.Rel{
			rest.Rel{"POST": "/collections",
				"Params": "name",
			},
		},
	}
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
		rest.Rel{"POST": "/collections/" + c.Collection.Id.Hex() + "/strings",
			"Params": "string",
		},
		rest.Rel{"DELETE": "/collections/" + c.Collection.Id.Hex() + "/strings/{StringId}"},
	})
}

func (c *CollectionStrings) Post(v *url.Values) (int, rest.APIResponse) {

	// Check Id
	if !c.Collection.Id.Valid() {
		return 4041, rest.NotFoundError()
	}

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		fmt.Println("POST /collections/" + c.Collection.Id.Hex() + "/strings - DB Connection Error")
		return 5001, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	// Find Collection
	err = C.FindId(c.Collection.Id).One(&c.Collection)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("POST /collections/" + c.Collection.Id.Hex() + "/strings - Collection Query Error")
		return 5002, rest.ServerError()
	}
	if err == mgo.ErrNotFound {
		return 4042, rest.NotFoundError()
	}

	// Validate string
	str := v.Get("string")
	if str == "" {
		return 422, &rest.APIError{
			rest.Error: rest.ErrorMsg{
				Type:    "invalid-string",
				Message: "A non-empty string is required.",
				Code:    422,
				Param:   []string{"string"},
			},
		}
	}

	// Init new String struct
	s := String{}
	existingString := false

	// Search for similar string in Collection Strings array (makes "POST" idempotent)
	if len(c.Collection.Strings) > 0 {
		for _, Item := range c.Collection.Strings {
			if Item.String == str {
				s.Id = Item.Id
				s.String = str
				existingString = true
				break
			}
		}
	}

	// Search for same String in DB
	S := session.DB(mongoDb).C("strings")
	err = S.Find(bson.M{"string": str}).One(&s)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("POST /collections/" + c.Collection.Id.Hex() + "/strings - Strings Query Error")
		return 5003, rest.ServerError()
	}
	fmt.Println(s.Id.Hex())
	// Create new String
	if s.Id == "" {

		// Set string
		s.Id = bson.NewObjectId()
		s.String = str

		// Translate string into x languages!
		gTranslateUrl := "https://www.googleapis.com/language/translate/v2"
		type GTranslation struct {
			Data struct {
				Translations []struct {
					TranslatedText string
				}
			}
		}

		// Create channel
		ch := make(chan Translation)
		fmt.Println("here2?")

		// Loop through languages
		it := 0
		for lang := range gLangs {
			fmt.Println("here1?")

			// Create a goroutine closure for every translation and collect the results into the channel
			go func(lang string, ch chan Translation) {
				fmt.Println("here?")
				// Initialize Translation struct t
				var t Translation

				// Prepare Values for GTranslate API call
				v := &url.Values{}
				v.Set("key", os.Getenv("GTRANSLATE_KEY"))
				v.Set("q", s.String)
				v.Set("source", "en")
				v.Set("target", lang)
				v.Set("prettyprint", "false")

				// Make GTranslate API Call and unmarshal json response
				url := gTranslateUrl + "?" + v.Encode()
				r, err := http.Get(url)
				if err != nil {
					ch <- t
					return // abort mission
				}
				fmt.Println("Translate:", url)
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					ch <- t
					return // abort mission
				}
				var g GTranslation
				err = json.Unmarshal(body, &g)
				if err != nil {
					ch <- t
					return // abort mission
				}
				// If a translation is returned, set values of t
				if len(g.Data.Translations) > 0 {
					t = Translation{
						Language:    lang,
						Translation: g.Data.Translations[0].TranslatedText,
					}
				}

				// Return Translation to chan
				ch <- t

			}(lang, ch)
			it++
		}
		s.Translations = make(map[string]string)

		// Wait for the goroutines to finish and save translations into String
		for i := 0; i < it; i++ {
			translation := <-ch
			if translation.Translation != "" {
				s.Translations[translation.Language] = translation.Translation
			}
		}

		// Insert new string into strings DB
		err = S.Insert(s)
		if err != nil {
			fmt.Println("POST /collections/" + c.Collection.Id.Hex() + "/strings - String Insert Error")
			return 5004, rest.ServerError()
		}

	}

	if !existingString {
		// Add String to Collection and Update Collection
		c.Collection.Strings = append(c.Collection.Strings, s)
		err = C.UpdateId(c.Collection.Id, c.Collection)
		if err != nil {
			fmt.Println("POST /collections/" + c.Collection.Id.Hex() + "/strings - Update Collection Error")
			return 5006, rest.ServerError()
		}
	}

	return 200, &rest.APISuccess{
		"String": s,
		"Next": &[]rest.Rel{
			rest.Rel{"DELETE": "/collections/" + c.Collection.Id.Hex() + "/strings/" + s.Id.Hex()},
		},
	}
}

func (c *CollectionStrings) Put(v *url.Values) (int, rest.APIResponse) {
	return 405, rest.InvalidMethodError(&[]rest.Rel{
		rest.Rel{"POST": "/collections/" + c.Collection.Id.Hex() + "/strings",
			"Params": "string",
		},
		rest.Rel{"DELETE": "/collections/" + c.Collection.Id.Hex() + "/strings/{StringId}"},
	})
}

func (c *CollectionStrings) Delete(v *url.Values) (int, rest.APIResponse) {

	// Initialize DB
	session, err := mgo.Dial(mongoPath)
	if err != nil {
		fmt.Println("DELETE /collections/" + c.Collection.Id.Hex() + "/strings/" + c.String.Id.Hex() + " - DB Connection Error")
		return 500, rest.ServerError()
	}
	C := session.DB(mongoDb).C("collections")
	defer session.Close()

	// Find collection
	err = C.FindId(c.Collection.Id).One(&c.Collection)
	if err != nil && err != mgo.ErrNotFound {
		return 500, rest.ServerError()
	}
	if err == mgo.ErrNotFound {
		return 404, rest.NotFoundError()
	}

	// Loop through Strings
	for i, String := range c.Collection.Strings {
		if String.Id.Hex() == c.String.Id.Hex() {
			// Remove string from collection
			c.Collection.Strings = append(c.Collection.Strings[:i], c.Collection.Strings[i+1:]...)
		}
	}
	return 200, &rest.APISuccess{
		"Success": true,
		"Next": &[]rest.Rel{
			rest.Rel{"POST": "/collections/" + c.Collection.Id.Hex() + "/strings",
				"Params": "strings",
			},
		},
	}
}
