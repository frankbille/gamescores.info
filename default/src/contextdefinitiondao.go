package defaultapp

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
)

const entityContextDefinition string = "ContextDefinition"

type contextDefinitionDao struct {
	dao
}

func createContextDefinitionDao(c *gin.Context) contextDefinitionDao {
	dao := createDao(getGaeContext(c))
	return contextDefinitionDao{dao}
}

func (dao *contextDefinitionDao) checkIDExists(ID string) (bool, error) {
	key := datastore.NewKey(dao.Context, entityContextDefinition, ID, 0, nil)

	err := datastore.Get(dao.Context, key, &ContextDefinition{})

	if err == nil {
		return true, nil
	} else {
		if err == datastore.ErrNoSuchEntity {
			return false, nil
		} else {
			return false, err
		}
	}
}

func (dao *contextDefinitionDao) saveContext(contextDefinition ContextDefinition) error {
	key := datastore.NewKey(dao.Context, entityContextDefinition, contextDefinition.ID, 0, nil)

	_, err := dao.save(key, &contextDefinition)

	return err
}