package tests

import (
	"testing"
	"github.com/kalcok/jc/tools"
	"os"
	"github.com/kalcok/jc"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"gopkg.in/mgo.v2"
)

var (
	sessionDB  string
	mgoSession *mgo.Session
)

func initTestSession() {
	address, found := os.LookupEnv("JC_TEST_ADDRESS")
	if !found {
		address = "localhost"
	}

	db, found := os.LookupEnv("JC_TEST_DB")
	if !found {
		db = "jc_test"
	}
	sessionDB = db

	pass, _ := os.LookupEnv("JC_TEST_PASS")
	user, _ := os.LookupEnv("JC_TEST_USER")

	conf := tools.SessionConf{Database: db, Addrs: []string{address}, Password: pass, Username: user}
	tools.InitSession(&conf)
}

func dropTestDB() {
	session, err := tools.GetSessionClone()
	if err != nil {
		panic(err)
	}
	session.DB("").DropDatabase()
	session.Close()
}

func TestMain(m *testing.M) {
	var err error
	mgoSession, err = mgo.Dial("linux.dev")
	if err != nil {
		panic(err)
	}
	defer mgoSession.Clone()
	initTestSession()
	dropTestDB()
	m.Run()
	dropTestDB()
}

func TestSingleDocumentInit(t *testing.T) {
	doc := ExplicitID{Data: "jc_test", MyID: 1001}
	err := jc.NewDocument(&doc)
	if err != nil {
		t.Error(err)
	}
}

func TestImplicitCollectionName(t *testing.T) {
	err_message := "Failed to implicitly create collection name. Expected '%s', got '%s'"

	type camelCased ImplicitID
	cc_string := "camel_cased"

	type TrueCamelCased ImplicitID
	tcc_string := "true_camel_cased"

	type snake_cased ImplicitID
	sc_string := "snake_cased"

	type plain ImplicitID
	plain_string := "plain"

	cc := camelCased{Data: "camelCased"}
	jc.NewDocument(&cc)
	if cc.CollectionName() != cc_string {
		t.Error(err_message, cc_string, cc.CollectionName())
	}

	tcc := TrueCamelCased{Data: "TrueCamelCased"}
	jc.NewDocument(&tcc)
	if tcc.CollectionName() != tcc_string {
		t.Error(err_message, tcc_string, tcc.CollectionName())
	}

	sc := snake_cased{Data: "snake_cased"}
	jc.NewDocument(&sc)
	if sc.CollectionName() != sc_string {
		t.Error(err_message, sc_string, sc.CollectionName())
	}

	p := plain{Data: "plain"}
	jc.NewDocument(&p)
	if p.CollectionName() != plain_string {
		t.Error(err_message, plain_string, p.CollectionName())
	}
}

func TestDefaultDB(t *testing.T) {
	doc := ImplicitID{Data: "TestDefaultDB"}
	jc.NewDocument(&doc)

	if doc.Database() != sessionDB {
		t.Error("failed to init document default DB")
	}
}

func TestDBChange(t *testing.T) {
	doc := ImplicitID{Data: "TestDBChange"}
	jc.NewDocument(&doc)
	new_db := sessionDB + "_fancy"
	doc.SetDatabase(new_db)

	if doc.Database() != new_db {
		t.Error("failed to change documents DB")
	}
}

func TestGetField(t *testing.T) {
	data := "TestGetField"
	doc := ImplicitID{Data: data}
	jc.NewDocument(&doc)
	field, err := doc.GetField("Data")
	if err != nil {
		t.Error("Failed to get field value by name")
	}
	if field != data {
		t.Error("Unexpected value.")
	}
}

func TestGetFieldProtectUnexported(t *testing.T) {
	type unexportedFieldDoc struct {
		jc.Collection
		Data    string `bson:"data"json:"data"`
		private string
	}

	doc := unexportedFieldDoc{Data: "TestGetFieldProtectUnexported", private: "Unexported"}
	jc.NewDocument(&doc)
	_, err := doc.GetField("private")
	if err == nil {
		t.Error("GetField accessed unexported field")
	}
}

func TestGetFieldMissing(t *testing.T) {
	type testDocument struct {
		jc.Collection
		Data string `bson:"data"json:"data"`
	}
	doc := testDocument{Data: "TestGetFieldMissing"}
	jc.NewDocument(&doc)
	_, err := doc.GetField("DefinitelyNotPresent")
	if err == nil {
		t.Error("GetField didn't fail on missing field access")
	}
}

func TestExplicitCollectionName(t *testing.T) {
	err_msg := "Failed to set explicit Collection name. Expected '%s', got '%s'"
	type boringStruct struct {
		jc.Collection `bson:"-"json:"-"jc:"my_awesome_collection"`
		Data string   `bson:"data",json:"data"`
	}

	doc := boringStruct{Data: "Boring data"}
	jc.NewDocument(&doc)

	if doc.CollectionName() != "my_awesome_collection" {
		t.Error(err_msg, "my_awesome_collection", doc.CollectionName())
	}

}

func TestInsertExplicitIDCollection(t *testing.T) {
	id := 666
	doc := ExplicitID{Data: "TestInsertExplicitIDCollection", MyID: 666}
	err := jc.NewDocument(&doc)
	if err != nil {
		return
	}

	// Test Save() call
	_, err = doc.Save(true)
	if err != nil {
		t.Error(err)
	}

	// Test if document is in DB
	result := bson.M{}
	empty_result := bson.M{}
	session, _ := tools.GetSessionClone()
	defer session.Close()
	session.DB(sessionDB).C("explicit_i_d").FindId(id).One(&result)
	if reflect.DeepEqual(result, empty_result) {
		t.Error("Failed to insert document with explicit ID into DB")
	}
}

func TestInsertImplicitIDCollection(t *testing.T) {
	doc := ImplicitID{Data: "TestInsertImplicitIDCollection"}
	err := jc.NewDocument(&doc)

	if err != nil {
		return
	}

	// Test Save() call
	_, err = doc.Save(true)
	if err != nil {
		t.Error(err)
	}
	id := doc.ID()

	// Test if document is in DB
	result := bson.M{}
	empty_result := bson.M{}
	session, _ := tools.GetSessionClone()
	defer session.Close()
	session.DB(sessionDB).C("implicit_i_d").FindId(id).One(&result)
	if reflect.DeepEqual(result, empty_result) {
		t.Error("Failed to insert document with implicit ID into DB")
	}
}

func TestUpsert(t *testing.T) {
	original := "TestUpsert"
	update := "TestUpsertUpdated"
	doc := ImplicitID{Data: original}
	jc.NewDocument(&doc)

	doc.Save(true)
	id := doc.ID()

	doc.Data = update
	doc.Save(true)

	result := bson.M{}
	session, _ := tools.GetSessionClone()
	defer session.Close()
	session.DB(sessionDB).C("implicit_i_d").FindId(id).One(&result)

	if result["data"] != update {
		t.Error("Failed to Upsert document")
	}

}

func TestNewImplicitID(t *testing.T) {
	doc := ImplicitID{Data: "TestNewID"}
	jc.NewDocument(&doc)

	originalID := doc.ID()
	err := doc.NewImplicitID()
	newID := doc.ID()

	if err != nil {
		t.Error("Failed to reinitialize new implicit ID: %s", err)
	}
	if originalID == newID {
		t.Error("Failed to reinitialize new implicit ID: New ID matches old ID")
	}
}

func TestNewImplicitIDError(t *testing.T) {
	doc := ExplicitID{Data: "TestNewImplicitIDError"}
	jc.NewDocument(&doc)

	err := doc.NewImplicitID()

	if err == nil {
		t.Error("NewImplicitID() failed to return error when used on document with explicit ID")
	}
}

func TestDocumentStartsUninitialized(t *testing.T) {
	doc := ImplicitID{Data: "TestDocumentStartsUninitialized"}

	if doc.IsInitialized() {
		t.Error("IsInitialized returns true for uninitialized document")
	}
}

func TestIsInitializedFlag(t *testing.T) {
	doc := ImplicitID{Data: "TestIsInitializedFlag"}

	jc.NewDocument(&doc)

	if !doc.IsInitialized() {
		t.Error("Initialization did not set _initialized flag")
	}
}