package tests

import (
	"testing"
	"github.com/kalcok/jc"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

// Creates <num> of records of ImplicitID type
func prepareSimpleRecords(num int, customData string) (ids []interface{}) {
	var doc ImplicitID

	if customData == "" {
		customData = "prepareSimpleRecords"
	}

	for i := 0; i < num; i++ {
		doc = ImplicitID{Data: customData}
		jc.NewDocument(&doc)
		doc.Save(true)
		ids = append(ids, doc.ID())
	}
	return ids
}

func TestQueryInit(t *testing.T) {
	var docs []ImplicitID
	_, err := jc.NewQuery(&docs)

	if err != nil {
		t.Error(fmt.Sprintf("Failed to initialize new query: %s", err))
	}
}

func TestQueryImplicitCollection(t *testing.T) {
	var docs []ImplicitCollection
	expectedCollection := "implicit_collection"

	q, _ := jc.NewQuery(&docs)

	if q.Collection() != expectedCollection {
		t.Error(fmt.Sprintf("Failed to set implicit collection to query. Expected '%s', got '%s'",
			q.Collection(),
			expectedCollection))
	}
}

func TestQueryExplicitCollection(t *testing.T) {
	var docs []ExplicitCollection
	expectedCollection := "my_collection"

	q, _ := jc.NewQuery(&docs)

	if q.Collection() != expectedCollection {
		t.Error(fmt.Sprintf("Failed to set explicit collection to query. Expected '%s', got '%s'",
			q.Collection(),
			expectedCollection))
	}
}

func TestQuerySelectSliceResult(t *testing.T) {
	prepareSimpleRecords(2, "")
	var docs []ImplicitID

	q, _ := jc.NewQuery(&docs)

	err := q.Execute(true)

	if err != nil {
		t.Error(fmt.Sprintf("Failed to select multiple documents from DB. Error: %s", err))
	}
}

func TestQuerySelectNonSliceResult(t *testing.T) {
	data := "single record"
	prepareSimpleRecords(2, data)
	var doc ImplicitID

	q, _ := jc.NewQuery(&doc)

	err := q.Execute(true)

	if err != nil {
		t.Error(fmt.Sprintf("Failed to select single document from DB. Error: %s", err))
	}
}

func TestQuerySelectAll(t *testing.T) {
	dropTestDB()
	recordNumber := 14
	prepareSimpleRecords(recordNumber, "")
	var docs []ImplicitID

	q, _ := jc.NewQuery(&docs)

	q.Execute(true)

	if len(docs) != recordNumber {
		t.Error(fmt.Sprintf("Failed to select All records from DB. Expected %d, got %d",
			recordNumber, len(docs)))
	}
}

func TestQuerySliceIntegrity(t *testing.T) {
	const recordsNumber = 10

	var data [recordsNumber]string
	var results []ExplicitID

	dropTestDB()

	for i := 0; i < recordsNumber; i++ {
		data[i] = fmt.Sprintf("Integrity test %d", i)
		doc := ExplicitID{MyID: i, Data: data[i]}

		jc.NewDocument(&doc)
		doc.Save(true)
	}

	q, _ := jc.NewQuery(&results)

	q.Execute(true)

	for _, record := range results {
		if record.Data != data[record.MyID] {
			t.Error(
				fmt.Sprintf("Integrity failed for record ID '%d'. Expected data '%s', got '%s'",
					record.MyID, data[record.MyID], record.Data))
		}
	}
}

func TestQueryNonSliceIntegrity(t *testing.T) {
	dropTestDB()
	data := "NonSliceIntegrity"
	prepareSimpleRecords(3, data)

	var doc ImplicitID

	q, _ := jc.NewQuery(&doc)
	q.Execute(true)

	if doc.Data != data {
		t.Error(fmt.Sprintf("Integrity failed for single record query. Expected '%s', got '%s'", data, doc.Data))
	}

}

func TestQuerySliceResetOnMultipleSelects(t *testing.T) {
	dropTestDB()
	firstData := "First batch"
	firstBatch := 10
	secondData := "Second batch"
	secondBatch := 5
	prepareSimpleRecords(firstBatch, firstData)

	docs := []ImplicitID{}

	q, _ := jc.NewQuery(&docs)
	q.Execute(true)

	dropTestDB()
	prepareSimpleRecords(secondBatch, secondData)
	q.Execute(true)

	if len(docs) != secondBatch {
		t.Error("Failed to clean up target slice between multiple selects")
	}
}

func TestQueryLimit(t *testing.T) {
	limit := 5
	prepareSimpleRecords(10, "TestQueryLimit")
	docs := []ImplicitID{}

	q, _ := jc.NewQuery(&docs)
	q.Limit(limit).Execute(true)

	if len(docs) != limit {
		t.Error("Failed to impose query limit")
	}
}

func TestQuerySkip(t *testing.T) {
	dropTestDB()
	skip := 1

	sample := ExplicitID{MyID: 1, Data: "TestQuery"}
	jc.NewDocument(&sample)
	sample.Save(true)
	sample.MyID = 2
	sample.Save(true)
	sample.MyID = 3
	sample.Save(true)

	result := []ExplicitID{}
	q, _ := jc.NewQuery(&result)
	q.Skip(skip).Execute(true)

	for _, record := range result {
		if record.MyID == 1 {
			t.Error("Failed to skip over entries in DB")
		}
	}
}

func TestQueryFilter(t *testing.T) {
	expectedData := "This is what we want"
	prepareSimpleRecords(3, "Not what We want")

	expected := ImplicitID{Data: expectedData}
	jc.NewDocument(&expected)

	expected.Save(true)

	result := ImplicitID{}
	q, _ := jc.NewQuery(&result)

	q.Filter(bson.M{"data": expectedData}).Execute(true)

	if result.Data != expectedData {
		t.Error("Failed to select records based on bson filter")
	}
}

func TestQueryFilterStruct(t *testing.T) {
	dropTestDB()
	expectedData := "This is what we want"
	prepareSimpleRecords(3, "Not what we want")
	expected := ImplicitID{Data: expectedData}
	jc.NewDocument(&expected)

	expected.Save(true)

	result := ImplicitID{}

	q, _ := jc.NewQuery(&result)

	q.Filter(ImplicitID{Data: expectedData}).Execute(true)

	if result.Data != expectedData {
		t.Error("Failed to select records based on struct filter")
	}
}