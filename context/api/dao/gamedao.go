package dao

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
	"api/utils"
	"api/domain"
)

const EntityGame string = "Game"

type GameDao struct {
	dao
}

func CreateGameDao(c *gin.Context) GameDao {
	dao := createDao(utils.GetGaeContext(c))
	return GameDao{dao}
}

func (dao *GameDao) GetGames(start, limit int, leagueID int64) ([]domain.Game, int, error) {
	var games []domain.Game

	leagueKey := datastore.NewKey(dao.Context, EntityLeague, "", leagueID, nil)
	count, err := dao.getListForAncestor(EntityGame, start, limit, leagueKey, []string{"-GameDate"}, &games)

	return games, count, err
}

func (dao *GameDao) GetGame(leagueID, gameID int64) (*domain.Game, error) {
	var game domain.Game
	leagueKey := datastore.NewKey(dao.Context, EntityLeague, "", leagueID, nil)
	key := datastore.NewKey(dao.Context, EntityGame, "", gameID, leagueKey)

	err := dao.get(key, &game)

	return &game, err
}

func (dao *GameDao) SaveGame(game domain.Game) (*domain.Game, error) {
	if game.ID == 0 {
		gameID, _, _ := datastore.AllocateIDs(dao.Context, EntityGame, nil, 1)
		game.ID = gameID
	}

	leagueKey := datastore.NewKey(dao.Context, EntityLeague, "", game.LeagueID, nil)
	game.League = leagueKey
	key := datastore.NewKey(dao.Context, EntityGame, "", game.ID, leagueKey)

	g, err := dao.save(key, &game)

	if err != nil {
		return nil, err
	}

	savedGame := g.(*domain.Game)

	return savedGame, err
}
