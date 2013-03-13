package main

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"net/url"
	"testing"
	"translation.io/rest"
)

func TestCollections(t *testing.T) {

	type Single struct {
		Collection Collection
	}

	var (
		c      Collection
		status int
		res    rest.APIResponse
		p      Single
		err    error
	)

	t.Log("POST Collection")
	name := "Test Collection " + bson.NewObjectId().String() // random title

	v := &url.Values{}
	v.Set("name", name)

	status, res = c.Post(v)
	if status != 200 {
		t.Errorf("Could not POST collection, status: '%d'", status)
		t.Log(res.ToJSON())
	}

	err = json.Unmarshal([]byte(res.ToJSON()), &p)
	if err != nil {
		t.Error(err)
	}

	c.Id = p.Collection.Id

	t.Log("GET Collection")
	status, res = c.Get(&url.Values{})
	if status != 200 {
		t.Errorf("Could not GET collection, status: '%d'", status)
		t.Log(res.ToJSON())
	}

	t.Log("PUT Collection")

	newName := "Updated Collection " + bson.NewObjectId().String() // random title

	v = &url.Values{}
	v.Set("name", newName)

	status, res = c.Put(v)
	if status != 200 {
		t.Errorf("Could not PUT collection, status: '%d'", status)
		t.Log(res.ToJSON())
	}

	if c.Name != newName {
		t.Errorf("Failed to change name")
	}

	t.Log("DELETE Collection")

	status, res = c.Delete(&url.Values{})
	if status != 200 {
		t.Errorf("Could not DELETE collection, status: '%d'", status)
		t.Log(res.ToJSON())
	}

}

func TestStrings(t *testing.T) {

	cs := &CollectionStrings{}
	cs.Collection.Id = bson.ObjectIdHex("513edd375a8c9b2fed000001")
	str := "Welcome Stranger! " + bson.NewObjectId().String()

	v := &url.Values{}
	v.Set("string", str)

	status, res := cs.Post(v)
	if status != 200 {
		t.Errorf("Could not POST to CollectionStrings, status: %d", status)
		t.Log(res.ToJSON())
	}
}
