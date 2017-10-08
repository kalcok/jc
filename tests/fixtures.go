package tests

import (
	"github.com/kalcok/jc"
)

type (
	ImplicitID struct {
		jc.Collection `bson:"-",json:"-"`
		Data string   `bson:"data",json:"data"`
	}

	ExplicitID struct {
		jc.Collection `bson:"-",json:"-"`
		MyID int      `bson:"_id",json:"myID"`
		Data string   `bson:"data",json:"data"`
	}
)
