package context

import (
	appengine "appengine"
	datastore "appengine/datastore"
	fmt "fmt"
)

type dao struct {
	Context appengine.Context
}

func createDao(gaeCtx appengine.Context) dao {
	return dao{
		Context: gaeCtx,
	}
}

func (dao *dao) get(key *datastore.Key, dataObject interface{}) error {
	return datastore.Get(dao.Context, key, dataObject)
}

func (dao *dao) getByIds(keys []*datastore.Key, dataObjects interface{}) error {
	return datastore.GetMulti(dao.Context, keys, dataObjects)
}

func (dao *dao) getList(entityType string, start, limit int, dataObjects interface{}) (int, error) {
	q := datastore.NewQuery(entityType).
		Offset(start).
		Limit(limit)

	_, err := q.GetAll(dao.Context, dataObjects)

	if err != nil {
		return 0, err
	}

	return datastore.NewQuery(entityType).
		Count(dao.Context)
}

func (dao *dao) getListForAncestor(entityType string, start, limit int, ancestor *datastore.Key, orderByFields []string, dataObjects interface{}) (int, error) {
	q := datastore.NewQuery(entityType).
		Ancestor(ancestor).
		Offset(start).
		Limit(limit)

	if orderByFields != nil {
		for _, orderByField := range orderByFields {
			q = q.Order(orderByField)
		}
	}

	_, err := q.GetAll(dao.Context, dataObjects)

	if err != nil {
		return 0, err
	}

	return datastore.NewQuery(entityType).
		Ancestor(ancestor).
		Count(dao.Context)
}

func (dao *dao) getByFilter(entityType, propertyName, propertyValue string, dataObjects interface{}) error {
	q := datastore.NewQuery(entityType).
		Filter(fmt.Sprintf("%s =", propertyName), propertyValue)

	_, err := q.GetAll(dao.Context, dataObjects)

	return err
}

func (dao *dao) save(key *datastore.Key, obj interface{}) (interface{}, error) {
	_, err := datastore.Put(dao.Context, key, obj)

	return obj, err
}
