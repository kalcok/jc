# jc
jc (pronounced as 'juice') is Object Document Manager (ODM) in golang built on top of `mgo` MongoDB driver. Aim is to provide abstraction layer on top of `mgo` without too much of a performance impact.

## Usage
Things you can do with `jc` (so far) :
 * Defining models
 * Saving / Updating Documents
 * Creating Queries to fetch documents from DB
 
### Document definition
You can turn any golang structure into model just by including anonymous field `jc.Collection`. 

**Simple model**
```golang
type Person struct {
	jc.Collection 		`bson:"-"json:"-"`
	FirstName string
	LastName  string
}

```
Notice the *json* and *bson* structure tags for the `jc.Collection`. These are very usefull for controll over Un/Marshalling of structures. In this case they say that `jc.Collection` field should be ignored when Saving/Reading from DB (and also when serializing into *json*)

You can create models with either implicit or explicit ID field. In case of implicit ID (as in example above), every document gets automatically generated `mgo.ObjectID`. In case of document with explicit ID, it needs to be supplied before trying to save document into DB.

**Model with explicit ID**
```golang
type Employee struct {
	jc.Collection 		`bson:"-"json:"-"`
	BadgeID   int 		`bson:"_id"`
	FirstName string
	LastName  string
}
```
In this example we used *bson* structure tag to indicate that we want to store value from `BadgeID` field into `_id` field in DB, which is what MongoDB uses to store document indices. 

*bson* tags can be used also for other "non-special" fields,for example, to translate native golang *CamleCase* names into more "standard" (at least for MongoDB) *snake_case* names.

**Custom field names**
```golang
type Person struct {
	jc.Collection 		`bson:"-"json:"-"`
	FirstName string	`bson:"first_name"`
	LastName  string	`bson:"last_name"`
}
```

You can also use struct tags to explicitly define collection name.

**Custom Collection name**
```golang

type NewYorkBranchEmployee struct {
	jc.Collection		`bson:"-"json:"-"jc:"ny_employees"`
	FirstName string
	LastName  string
}
```
### Session Management
Before using any `jc` features, there needs to be initialized session with MongoDB server. Master session is initialized by calling `jc.tools.InitSession()` which takes one argument in form of `jc.toosl.SessionConf` struct. `Sessionconf` is just convenient alias to `mgo.DialInfo`. After you don't need master session anymore, it can be closed with call to `jc.tools.CloseSession()`.

**Eample**
```golang
func main(){
	conf := tools.SessionConf{Addrs: []string{"localhost"}, Database: "jc_test"}
	tools.InitSession(&conf)
	defer tools.CloseSession()
	// program logic
	// ...
}
```
You usually don't need to be concerned with session for the rest of your program after initialization. However if you need direct access to `mgo.Session` object, you can get either clone or copy of master session by calling `jc.tools.GetSessionClone()` or `jc.tools.GetSessionCopy()` respectivelly. Don't forget that these session need to be closed separately by calling `mySession.Close()`
_______________________________________________________________
Note: `tools` package is ripe for renaming
### Document Initiation
Due to some Golang restrictions (mainly, missing constructors) new document instance must be initiated with call to `NewDocument()` which takes pointer to a struct that contains `jc.Collection` as an argument

**Example**
```golang
newPerson := Person{FirstName:"John", LastName:"Foo"}
err := jc.NewDocument(&newPerson)

if err != nil {
	panic(err)
}
```
### Inserting Document (Save)
Once document is properly initialized, it can be inserted into DB by calling Save() method.

**Example**
```golang
changeInfo, err := newPerson.Save(true)

if err != nil {
	panic(err)
}
```
`Save()` takes one boolean argument which decide whether the action will [Clone](https://godoc.org/gopkg.in/mgo.v2#Session.Clone) (`true`) or [Copy](https://godoc.org/gopkg.in/mgo.v2#Session.Copy) (`false`) master session.
Calling `Save()` multiple times on the same document works like Update, to insert new document instance into DB, documents ID must change. On documents that do not have explicit ID field, you can call `NewImplicitID` to generate new ID for the document

### Query
`Query` object is used to pull documents from DB. Query must be instantiated by calling `NewQuery()` and it takes single argument in form of pointer to either single document or slice of documents where eventual result will be saved. Query autmatically recognizes whether its target is single document or slice and adjusts final action to perform either `One()` or `All()` query. (**Note:** Documents passed into `NewQuery()` will be automatically initialized so there is no need to call `NewDocument()` manually)

**Example**
```golang
var person Person
var allPersons []Person

// When executed, this query will pull single record from DB into person object
oneQuery, err := jc.NewQuery(&person)
if err != nil {
	panic(err)
}

// When executed, this query will pull all records that conform given filter and save it into allPersons slice
multiQuery, err := jc.NewQuery(&allPersons)
if err != nil {
	panic(err)
}
```

#### Query restrictions
`jc` supports following restrictions that shape its finall result:
 * `Limit(int)` - Sets maximum number of records pulled from DB
 * `Skip(int)` - Ignore arbitrary number of records form the beggining of result set
 * `Filter(interface{})` - Takes map or struct and applies it as filter on the final result
All the restriction methods also return pointer to the query on which it was called so that they can be chained together.

**Example**
```golang
var tenJohns []Person
multiQuery, err := jc.NewQuery(&tenJohns)
if err != nil {
	panic(err)
}

multiQuery.Skip(5).Limit(10).Filter(Person{FirstName: "John"})
```
 
#### Query execution
To execute query and actually fill target document(s) with data from DB, you must call `Execute()` method. `Execute()` takes one boolean argument which decide whether the action will [Clone](https://godoc.org/gopkg.in/mgo.v2#Session.Clone) (`true`) or [Copy](https://godoc.org/gopkg.in/mgo.v2#Session.Copy) (`false`) master session.
\# TODO
