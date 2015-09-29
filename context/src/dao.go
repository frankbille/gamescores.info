package context

import (
	appengine "appengine"
	datastore "appengine/datastore"
	memcache "appengine/memcache"
	fmt "fmt"
	reflect "reflect"
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
	cacheKey := dao.createCacheKeyForDatastoreKey(key)

	_, err := memcache.JSON.Get(dao.Context, cacheKey, dataObject)

	if err == memcache.ErrCacheMiss {
		err = datastore.Get(dao.Context, key, dataObject)

		if err != nil {
			return err
		}

		dao.saveToCache(cacheKey, dataObject)

		return nil
	} else if err != nil {
		return err
	} else {
		return nil
	}
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

func (dao *dao) getByFilter(entityType, propertyName, propertyValue string, dataObjects interface{}) error {
	cacheKey := dao.createCacheKey(entityType, propertyName, propertyValue)

	_, err := memcache.JSON.Get(dao.Context, cacheKey, dataObjects)

	if err == memcache.ErrCacheMiss {
		q := datastore.NewQuery(entityType).
			Filter(fmt.Sprintf("%s =", propertyName), propertyValue)

		_, err := q.GetAll(dao.Context, dataObjects)

		if err != nil {
			return err
		}

		dov := reflect.ValueOf(dataObjects)
		if dov.Elem().Len() > 0 {
			dao.saveToCache(cacheKey, dataObjects)
		}

		return nil
	} else if err != nil {
		return err
	} else {
		return nil
	}
}

func (dao *dao) save(key *datastore.Key, obj interface{}) (interface{}, error) {
	_, err := datastore.Put(dao.Context, key, obj)

	if err != nil {
		return nil, err
	}

	dao.saveToCache(dao.createCacheKeyForDatastoreKey(key), obj)

	return obj, nil
}

func (dao *dao) createCacheKey(entityType, propertyName, propertyValue string) string {
	return fmt.Sprintf("%s-%s-%s", entityType, propertyName, propertyValue)
}

func (dao *dao) createCacheKeyForDatastoreKey(key *datastore.Key) string {
	return key.String()
}

func (dao *dao) saveToCache(cacheKey string, dataObject interface{}) {
	cacheItem := &memcache.Item{
		Key:    cacheKey,
		Object: dataObject,
	}
	memcache.JSON.Set(dao.Context, cacheItem)
}
