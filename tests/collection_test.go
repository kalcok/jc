package tests

import (
	"testing"
	"github.com/kalcok/jc/tools"
	"os"
	"github.com/kalcok/jc"
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

	pass, _ := os.LookupEnv("JC_TEST_PASS")
	user, _ := os.LookupEnv("JC_TEST_USER")

	conf := tools.SessionConf{Database: db, Addrs: []string{address}, Password: pass, Username: user}
	tools.InitSession(&conf)
}

func dropTestDB() {
	session, err := tools.GetSession()
	if err != nil {
		panic(err)
	}
	session.DB("").DropDatabase()
}

func TestMain(m *testing.M) {
	initTestSession()
	dropTestDB()
	m.Run()
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

	sc := snake_cased{Data: "snakce_cased"}
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
	doc := ExplicitID{Data: "TestInsertExplicitIDCollection", MyID: 666}
	err := jc.NewDocument(&doc)
	if err != nil {
		return
	}

	doc.Save(true)

}

func TestInsertImplicitIDCollection(t *testing.T) {
	doc := ImplicitID{Data: "TestInsertImplicitIDCollection"}
	err := jc.NewDocument(&doc)

	if err != nil {
		return
	}

	_, err = doc.Save(true)
	if err != nil {
		t.Error(err)
	}
}
