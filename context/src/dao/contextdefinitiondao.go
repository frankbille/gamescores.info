package dao

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
	"appengine/memcache"
	"fmt"
	"src/domain"
	"src/utils"
)

const entityContextDefinition string = "ContextDefinition"

type ContextDefinitionDao struct {
	dao
}

func CreateContextDefinitionDao(c *gin.Context) ContextDefinitionDao {
	dao := createDao(utils.GetGaeRootContext(c))
	return ContextDefinitionDao{dao}
}

func (dao *ContextDefinitionDao) GetContext(namespace string) (*domain.ContextDefinition, error) {
	var contextDefinition domain.ContextDefinition

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
