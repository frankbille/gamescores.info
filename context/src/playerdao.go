package context

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
)

const entityPlayer string = "Player"

type playerDao struct {
	dao
}

func createPlayerDao(c *gin.Context) playerDao {
	dao := createDao(getGaeContext(c))
	return playerDao{dao}
}

func (dao *playerDao) getPlayers(start, limit int) ([]Player, int, error) {
	var players []Player

	count, err := dao.getList(entityPlayer, start, limit, &players)

	return players, count, err
}

func (dao *playerDao) getAllPlayersByID(playerIds []int64) ([]Player, error) {
	players := make([]Player, len(playerIds))

	keys := make([]*datastore.Key, len(playerIds))
	idx := 0
	for _, playerID := range playerIds {
		key := datastore.NewKey(dao.Context, entityPlayer, "", playerID, nil)
		keys[idx] = key
		idx++
	}

	err := dao.getByIds(keys, players)

	return players, err
}

func (dao *playerDao) getPlayer(playerID int64) (*Player, error) {
	var player Player
	key := datastore.NewKey(dao.Context, entityPlayer, "", playerID, nil)

	err := dao.get(key, &player)

	return &player, err
}

func (dao *playerDao) savePlayer(player Player) (*Player, error) {
	if player.ID == 0 {
		playerID, _, _ := datastore.AllocateIDs(dao.Context, entityPlayer, nil, 1)
		player.ID = playerID
	}

	key := datastore.NewKey(dao.Context, entityPlayer, "", player.ID, nil)

	p, err := dao.save(key, &player)

	if err != nil {
		return nil, err
	}

	savedPlayer := p.(*Player)

	return savedPlayer, err
}
