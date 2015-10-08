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
	err := datastore.Get(dao.Context, key, dataObject)

	if err != nil {
		return err
	}

	return nil
}

func (dao *dao) getList(entityType string, start, limit int, dataObjects interface{}) (int, error) {
	q := datastore.NewQuery(entityType).
		Offset(start).
		Limit(limit)

	_, err := q.GetAll(dao.Context, dataObjects)

	if err != nil {
		return 0, err
	}

	count, err := datastore.NewQuery(entityType).
		Count(dao.Context)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *dao) getListForAncestor(entityType string, start, limit int, ancestor *datastore.Key, dataObjects interface{}) (int, error) {
	q := datastore.NewQuery(entityType).
		Ancestor(ancestor).
		Offset(start).
		Limit(limit)

	_, err := q.GetAll(dao.Context, dataObjects)

	if err != nil {
		return 0, err
	}

	count, err := datastore.NewQuery(entityType).
		Ancestor(ancestor).
		Count(dao.Context)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *dao) getByFilter(entityType, propertyName, propertyValue string, dataObjects interface{}) error {
	q := datastore.NewQuery(entityType).
		Filter(fmt.Sprintf("%s =", propertyName), propertyValue)

	_, err := q.GetAll(dao.Context, dataObjects)

	if err != nil {
		return err
	}

	return nil
}

func (dao *dao) save(key *datastore.Key, obj interface{}) (interface{}, error) {
	_, err := datastore.Put(dao.Context, key, obj)

	if err != nil {
		return nil, err
	}

	return obj, nil
}
