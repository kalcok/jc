package jc

import (
	"reflect"
	"unicode"
	"fmt"
	"errors"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/kalcok/jc/tools"
)

type document interface {
	setCollection(string)
	CollectionName() string
	SetDatabase(string)
	Database() string
	GetField(string) (interface{}, error)
	ID() interface{}
	Init(reflect.Value, reflect.Type)
	InitDB() error
	Info()
	NewImplicitID() error
	Save(bool) (*mgo.ChangeInfo, error)
	IsInitialized() bool
}

type Collection struct {
	_collectionName  string                `bson:"-"json:"-"`
	_collectionDB    string                `bson:"-"json:"-"`
	_parent          reflect.Value         `bson:"-"json:"-"`
	_parentType      reflect.Type          `bson:"-"json:"-"`
	_hasExplicitID   bool                  `bson:"-"json:"-"`
	_explicitIDField string                `bson:"-"json:"-"`
	_implicitIDValue bson.ObjectId         `bson:"-"json:"-"`
	_skeleton        []reflect.StructField `bson:"-"json:"-"`
	_initialized     bool                  `bson:"-",json:"-"`
}

func (c *Collection) setCollection(name string) {
	c._collectionName = name
}

func (c *Collection) CollectionName() string {
	return c._collectionName
}

func (c *Collection) SetDatabase(name string) {
	c._collectionDB = name

}

func (c *Collection) Database() string {
	return c._collectionDB
}

func (c *Collection) IsInitialized() bool {
	return c._initialized
}

func (c *Collection) Info() {
	fmt.Printf("Database %s\n", c._collectionDB)
	fmt.Printf("Collection %s\n", c._collectionName)
	fmt.Printf("Parent__ %s\n", c._parent)
}

func (c *Collection) Save(reuseSocket bool) (info *mgo.ChangeInfo, err error) {
	var session *mgo.Session
	var documentID interface{}
	idField := "_id"

	if reuseSocket {
		session, err = tools.GetSessionClone()
	} else {
		session, err = tools.GetSessionCopy()
	}
	if err != nil {
		return info, err
	}
	defer session.Close()

	if c._hasExplicitID {
		documentID = c._parent.Elem().FieldByName(c._explicitIDField).Interface()
	} else if len(c._implicitIDValue) > 0 {
		documentID = c._implicitIDValue
	} else {
		c._implicitIDValue = bson.NewObjectId()
		documentID = c._implicitIDValue
	}

	collection := session.DB(c._collectionDB).C(c._collectionName)
	info, err = collection.Upsert(bson.M{idField: documentID}, c._parent.Interface())

	return info, err
}

func (c *Collection) ID() (id interface{}) {
	if c._hasExplicitID {
		id = c._parent.Elem().FieldByName(c._explicitIDField).Interface()
	} else {
		id = c._implicitIDValue
	}
	return id
}

func (c *Collection) GetField(name string) (result interface{}, err error) {
	if unicode.IsLower(rune(name[0])) {
		err = errors.New(fmt.Sprintf("can't access unexported field '%s'", name))
		return
	}
	missingValue := reflect.Value{}
	resultValue := c._parent.Elem().FieldByName(name)
	if resultValue == missingValue {
		err = errors.New(fmt.Sprintf("field '%s' not found", name))
		return
	}
	result = resultValue.Interface()
	return
}

func (c *Collection) Init(parent reflect.Value, parentType reflect.Type) {

	c._parent = parent
	c._parentType = parentType
	c._hasExplicitID = false
	for i := 0; i < reflect.Indirect(c._parent).NumField(); i++ {
		field := c._parentType.Field(i)

		// Find explicit Collection name
		if field.Type == reflect.TypeOf(Collection{}) {
			explicitName := false
			jc_tag, tag_present := field.Tag.Lookup("jc")
			if tag_present {
				jc_fields := strings.Split(jc_tag, ",")
				if len(jc_fields) > 0 && jc_fields[0] != "" {
					c.setCollection(jc_fields[0])
					explicitName = true
				}
			}
			if !explicitName {
				c.setCollection(camelToSnake(parentType.Name()))
			}
		}

		// Find explicit index field
		bson_tag, tag_present := field.Tag.Lookup("bson")
		if tag_present {
			field_id := strings.Split(bson_tag, ",")
			switch field_id[0] {
			case "_id":
				c._explicitIDField = field.Name
				c._hasExplicitID = true
				break
			case "-":
				continue
			default:
				break
			}
		}
		c._skeleton = append(c._skeleton, field)
	}
	c._initialized = true

}

func (c *Collection) InitDB() error {
	session, err := tools.GetSessionClone()
	if err == nil {
		c.SetDatabase(session.DB("").Name)
		defer session.Close()
	} else {
		err = errors.New("database not initialized")
	}
	return err
}

func (c *Collection) NewImplicitID() (err error) {

	if c._hasExplicitID {
		return errors.New("can't assign new ID to document with Explicit ID")
	}
	c._implicitIDValue = bson.NewObjectId()
	return err
}

func NewDocument(c document) error {
	var err error
	objectType := reflect.TypeOf(c).Elem()
	objectValue := reflect.ValueOf(c)
	c.Init(objectValue, objectType)
	err = c.InitDB()

	return err
}

func camelToSnake(camel string) string {
	var (
		snake_name []rune
		next       rune
	)
	for i, c := range camel {
		if unicode.IsUpper(c) && i != 0 {
			snake_name = append(snake_name, '_')
		}
		next = unicode.ToLower(c)
		snake_name = append(snake_name, next)
	}
	return string(snake_name)
}

func initPrototype(prototype reflect.Value, internalType reflect.Type) reflect.Value {
	var inputs []reflect.Value
	inputs = append(inputs, reflect.ValueOf(prototype))
	inputs = append(inputs, reflect.ValueOf(internalType))
	prototype.MethodByName("Init").Call(inputs)
	prototype.MethodByName("InitDB").Call(nil)

	return reflect.Indirect(prototype)

}
