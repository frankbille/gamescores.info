package dao

import (
	datastore "appengine/datastore"
	"src/domain"
	"github.com/gamescores/gin"
	"src/utils"
)

const EntityUser string = "User"

type UserDao struct {
	dao
}

func CreateUserDao(c *gin.Context) UserDao {
	dao := createDao(utils.GetGaeRootContext(c))
	return UserDao{dao}
}

func (dao *UserDao) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User

	key := datastore.NewKey(dao.Context, EntityUser, userID, 0, nil)

	err := dao.get(key, &user)

	if err != nil && err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	return &user, err
}

func (dao *UserDao) SaveUser(user *domain.User) error {
	key := datastore.NewKey(dao.Context, EntityUser, user.UserID, 0, nil)

	_, err := dao.save(key, user)

	return err
}
