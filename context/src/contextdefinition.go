package context

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

func (cd ContextDefinition) isUserOwner(userKey *datastore.Key) bool {
	return userKey.Equal(cd.Owner)
}
