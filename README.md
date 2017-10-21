# jc
jc (pronounced as 'juice') is Object Document Manager (ODM) in golang built on top of `mgo` MongoDB driver. Aim is to provide abstraction layer on top of `mgo` without too much of a performance impact.

## Usage
Things you can do with `jc` (so far) :
 * Definind models
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

\#TODO

### Query
\# TODO
