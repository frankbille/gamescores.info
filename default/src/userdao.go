package defaultapp

import (
	datastore "appengine/datastore"
)

const entityUser string = "User"

type userDao struct {
	dao
}

func (dao *userDao) getUserByID(userID string) (*User, error) {
	var user User

	key := datastore.NewKey(dao.Context, entityUser, userID, 0, nil)

	err := dao.get(key, &user)

	if err != nil && err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	return &user, err
}

func (dao *userDao) saveUser(user *User) error {
	key := datastore.NewKey(dao.Context, entityUser, user.UserID, 0, nil)

	_, err := dao.save(key, user)

	return err
}
