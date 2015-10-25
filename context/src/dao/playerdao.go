package dao

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
	"src/domain"
	"src/utils"
)

const EntityPlayer string = "Player"

type PlayerDao struct {
	dao
}

func CreatePlayerDao(c *gin.Context) PlayerDao {
	dao := createDao(utils.GetGaeContext(c))
	return PlayerDao{dao}
}

func (dao *PlayerDao) GetPlayers(start, limit int) ([]domain.Player, int, error) {
	var players []domain.Player

	count, err := dao.getList(EntityPlayer, start, limit, &players)

	return players, count, err
}

func (dao *PlayerDao) GetAllPlayersByID(playerIds []int64) ([]domain.Player, error) {
	players := make([]domain.Player, len(playerIds))

	keys := make([]*datastore.Key, len(playerIds))
	idx := 0
	for _, playerID := range playerIds {
		key := datastore.NewKey(dao.Context, EntityPlayer, "", playerID, nil)
		keys[idx] = key
		idx++
	}

	err := dao.getByIds(keys, players)

	return players, err
}

func (dao *PlayerDao) GetPlayer(playerID int64) (*domain.Player, error) {
	var player domain.Player
	key := datastore.NewKey(dao.Context, EntityPlayer, "", playerID, nil)

	err := dao.get(key, &player)

	return &player, err
}

func (dao *PlayerDao) SavePlayer(player domain.Player) (*domain.Player, error) {
	if player.ID == 0 {
		playerID, _, _ := datastore.AllocateIDs(dao.Context, EntityPlayer, nil, 1)
		player.ID = playerID
	}

	key := datastore.NewKey(dao.Context, EntityPlayer, "", player.ID, nil)

	p, err := dao.save(key, &player)

	if err != nil {
		return nil, err
	}

	savedPlayer := p.(*domain.Player)

	return savedPlayer, err
}
