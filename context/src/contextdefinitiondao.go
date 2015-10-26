package context

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
	"appengine/memcache"
	"fmt"
)

const entityContextDefinition string = "ContextDefinition"

type contextDefinitionDao struct {
	dao
}

func createContextDefinitionDao(c *gin.Context) contextDefinitionDao {
	dao := createDao(getGaeRootContext(c))
	return contextDefinitionDao{dao}
}

func (dao *contextDefinitionDao) getContext(namespace string) (*ContextDefinition, error) {
	var contextDefinition ContextDefinition

	memCacheKey := fmt.Sprintf("NameSpace-%s", namespace)

	_, err := memcache.Gob.Get(dao.Context, memCacheKey, &contextDefinition)

	if err == memcache.ErrCacheMiss {
		key := datastore.NewKey(dao.Context, entityContextDefinition, namespace, 0, nil)

		err = dao.get(key, &contextDefinition)

		if err == nil {
			memcache.Gob.Set(dao.Context, &memcache.Item{
				Key: memCacheKey,
				Object: contextDefinition,
			})
		}
	}

	return &contextDefinition, err
}
