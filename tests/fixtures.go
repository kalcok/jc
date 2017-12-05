package tests

import (
	"github.com/kalcok/jc"
)

// Test Fixtures
type (
	ImplicitID struct {
		jc.Collection `bson:"-"json:"-"`
		Data string   `bson:"data"json:"data"`
	}

	ExplicitID struct {
		jc.Collection `bson:"-"json:"-"`
		MyID int      `bson:"_id"json:"myID"`
		Data string   `bson:"data"json:"data"`
	}
	ExplicitCollection struct {
		jc.Collection `bson:"-"json:"-"jc:"my_collection"`
		Data string   `bson:"data"json:"data"`
	}
	ImplicitCollection struct {
		jc.Collection `bson:"-"json:"-"`
		Data string   `bson:"data"json:"data"`
	}
)

// Benchmark Fixtures
type (
	simpleUserMGO struct {
		name    string
		surname string
		address string
		age     int
	}

	simpleUserJC struct {
		jc.Collection
		name    string
		surname string
		address string
		age     int
	}
)
