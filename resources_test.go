package main

import (
	"net/url"
	"testing"
)

func TestCollections(t *testing.T) {
	c := &Collection{}

	v := &url.Values{}
	v.Set("name", "Test Collection")

	t.Log("POST Collection")
	status, res := c.Post(v)

	if status != 200 {
		t.Errorf("Could not POST collection, status: '%d'", status)
		t.Log(res)
	}

}
