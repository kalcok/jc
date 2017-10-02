package wrappers

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type UuidField struct {
	Data uuid.UUID
}

func (u UuidField) GetBSON() (interface{}, error) {
	bytes, err := u.Data.MarshalBinary()
	return bson.Binary{Kind: 3, Data: bytes}, err
}

func (u *UuidField) SetBSON(raw bson.Raw) error {
	var bytes bson.Binary
	err := raw.Unmarshal(&bytes)
	if err == nil {
		err = u.Data.UnmarshalBinary(bytes.Data)
	}
	return err
}

func NewUuid() UuidField {
	return UuidField{Data: uuid.New()}
}
