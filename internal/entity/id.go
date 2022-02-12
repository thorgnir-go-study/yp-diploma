package entity

import "github.com/gofrs/uuid"

//ID entity ID
type ID = uuid.UUID

var NilID = uuid.Nil

//NewID create a new entity ID
func NewID() (ID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return NilID, err
	}
	return id, nil
}

//StringToID convert a string to an entity ID
func StringToID(s string) (ID, error) {
	id, err := uuid.FromString(s)
	return id, err
}
