package jc

import (
	"fmt"
	"reflect"
	"errors"
	"gopkg.in/mgo.v2"
	"github.com/kalcok/jc/tools"
)

type Query struct {
	Database    string
	collection  string
	filter      interface{}
	limit       int
	skip        int
	singleValue bool
	result      interface{}
	resultType  reflect.Type
}

func NewQuery(result interface{}) (newQuery Query, err error) {
	newQuery.result = result

	internalType := reflect.TypeOf(result).Elem()

	if internalType.Kind() == reflect.Slice {
		internalType = internalType.Elem()
		newQuery.singleValue = false
	} else {
		newQuery.singleValue = true
	}

	prototype := reflect.New(internalType)
	targetInterface := reflect.TypeOf((*document)(nil)).Elem()
	if !prototype.Type().Implements(targetInterface) {
		err = errors.New(
			fmt.Sprintf("Supplied 'result' type does not implement '%s'.", targetInterface.String()))
		return
	}

	newQuery.resultType = internalType
	proto_val := initPrototype(prototype, internalType)

	newQuery.collection = proto_val.FieldByName("_collectionName").String()
	newQuery.Database = proto_val.FieldByName("_collectionDB").String()

	return
}

func (q *Query) Execute(reuseSocket bool) (err error) {
	var session *mgo.Session

	if reuseSocket {
		session, err = tools.GetSessionClone()
	} else {
		session, err = tools.GetSessionCopy()
	}
	if err != nil {
		return
	}
	defer session.Close()

	query := session.DB(q.Database).C(q.collection).Find(q.filter)

	if q.skip > 0 {
		query = query.Skip(q.skip)
	}

	if q.limit > 0 {
		query = query.Limit(q.limit)
	}

	if q.singleValue {
		err = query.One(q.result)
	} else {
		err = query.All(q.result)
	}

	if err != nil {
		return
	}
	q.initResult()
	return
}

func (q *Query) Collection() string {
	return q.collection
}

func (q *Query) Filter(filter interface{}) *Query {
	q.filter = filter
	return q
}

func (q *Query) GetFilter() interface{} {
	return q.filter
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) GetLimit() int {
	return q.limit
}

func (q *Query) Skip(skip int) *Query {
	q.skip = skip
	return q
}

func (q *Query) GetSkip() int {
	return q.skip
}

func (q *Query) initResult() error {
	var err error
	resultv := reflect.ValueOf(q.result)

	if q.singleValue {
		initPrototype(resultv, q.resultType)
	} else {
		slicev := resultv.Elem()
		for i := 0; i < slicev.Len(); i++ {
			initPrototype(slicev.Index(i).Addr(), q.resultType)
		}
	}

	return err
}
