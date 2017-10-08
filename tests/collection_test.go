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
