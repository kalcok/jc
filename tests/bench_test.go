package tests

import (
	"testing"
	"github.com/kalcok/jc"
)

// Insert simple document with plain mgo
func BenchmarkMgoSave(b *testing.B) {
	session := mgoSession.Clone()
	collection := "BenchmarkMgoSave"
	data := simpleUserMGO{name: "Foo", surname: "Bar", address: "Foobar av"}

	for i := 0; i < b.N; i++ {
		data.age = i
		session.DB(sessionDB).C(collection).Insert(&data)
	}
}

// Insert documents with jc
// Creates new records with NewImplicitID call for better performance
func BenchmarkJcSave(b *testing.B) {
	data := simpleUserJC{name: "Foo", surname: "Bar", address: "Foobar av"}
	err := jc.NewDocument(&data)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		data.NewImplicitID()
		data.age = i
		data.Save(true)
	}
}

// Insert documents with jc
// For each insert, new jc document is initialized with NewDocument call
// Probably worst case scenario (performance wise)
func BenchmarkJcSaveReinit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := simpleUserJC{name: "Foo", surname: "Bar", address: "Foobar av", age: i}
		err := jc.NewDocument(&data)
		if err != nil {
			panic(err)
		}
		data.NewImplicitID()
		data.Save(true)
	}
}
