package main

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"net/url"
	"testing"
	"translation.io/rest"
)

func TestCollections(t *testing.T) {

	type List struct {
		Collections []Collection
	}

	type Single struct {
		Collection Collection
	}

	var (
		l      List
		c      Collection
		status int
		res    rest.APIResponse
		p      Single
		err    error
	)

	t.Log("GET Collections - first count")
	// get first count
	status, res = c.Get(&url.Values{})
	if status != 200 {
		t.Errorf("Could not GET collections, status: '%d'", status)
		t.Log(res.ToJSON())
	}
	err = json.Unmarshal([]byte(res.ToJSON()), &l)
	if err != nil {
		t.Error(err)
	}
	count := len(l.Collections)

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

	t.Log("GET Collections - 2nd count")
	// get second count
	c2 := &Collection{}
	status, res = c2.Get(&url.Values{})
	if status != 200 {
		t.Errorf("Could not GET collections, status: '%d'", status)
		t.Log(res.ToJSON())
	}
	l2 := &List{}
	err = json.Unmarshal([]byte(res.ToJSON()), l2)
	if err != nil {
		t.Error(err)
	}
	newCount := len(l2.Collections)
	if newCount <= count {
		t.Errorf("No new collections were created, count: %d, newCount: %d.", count, newCount)
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
