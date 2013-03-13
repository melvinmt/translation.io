package main

import (
	"encoding/json"
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
	name := "Test Collection " // + bson.NewObjectId().Hex() // random title

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

	newName := "Updated Collection " // + bson.NewObjectId().Hex() // random title

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

	t.Log("POST string to Collection")
	cs := &CollectionStrings{}
	cs.Collection.Id = c.Id
	str := "Welcome my friend! " // + bson.NewObjectId().Hex()

	v = &url.Values{}
	v.Set("string", str)

	status, res = cs.Post(v)
	if status != 200 {
		t.Errorf("Could not POST to CollectionStrings, status: %d", status)
		t.Log(res.ToJSON())
	}

	t.Log("DELETE string from Collection")
	c1 := len(cs.Collection.Strings)
	cs.String.Id = cs.Collection.Strings[0].Id
	cs.Delete(&url.Values{})
	c2 := len(cs.Collection.Strings)
	if c2 >= c1 {
		t.Errorf("String was not deleted from Collection!")
	}

	// t.Log("DELETE Collection")

	// status, res = c.Delete(&url.Values{})
	// if status != 200 {
	// 	t.Errorf("Could not DELETE collection, status: '%d'", status)
	// 	t.Log(res.ToJSON())
	// }

}
