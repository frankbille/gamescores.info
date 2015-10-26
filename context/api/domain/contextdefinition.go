package domain

import (
	"appengine/datastore"
)

type ContextDefinition struct {
	DefaultHalResource
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Active bool           `json:"active"`
	Owner  *datastore.Key `json:"-"`
}

func (cd ContextDefinition) IsUserOwner(userKey *datastore.Key) bool {
	return userKey.Equal(cd.Owner)
}
