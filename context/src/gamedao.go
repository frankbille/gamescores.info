package context

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
)

const entityGame string = "Game"

type gameDao struct {
	dao
}

func createGameDao(c *gin.Context) gameDao {
	dao := createDao(getGaeContext(c))
	return gameDao{dao}
}

func (dao *gameDao) getGames(start, limit int, leagueID int64) ([]Game, int, error) {
	var games []Game

	leagueKey := datastore.NewKey(dao.Context, entityLeague, "", leagueID, nil)
	count, err := dao.getListForAncestor(entityGame, start, limit, leagueKey, []string{"-GameDate"}, &games)

	return games, count, err
}

func (dao *gameDao) getGame(leagueID, gameID int64) (*Game, error) {
	var game Game
	leagueKey := datastore.NewKey(dao.Context, entityLeague, "", leagueID, nil)
	key := datastore.NewKey(dao.Context, entityGame, "", gameID, leagueKey)

	err := dao.get(key, &game)

	return &game, err
}

func (dao *gameDao) saveGame(game Game) (*Game, error) {
	if game.ID == 0 {
		gameID, _, _ := datastore.AllocateIDs(dao.Context, entityGame, nil, 1)
		game.ID = gameID
	}

	leagueKey := datastore.NewKey(dao.Context, entityLeague, "", game.LeagueID, nil)
	game.League = leagueKey
	key := datastore.NewKey(dao.Context, entityGame, "", game.ID, leagueKey)

	g, err := dao.save(key, &game)

	if err != nil {
		return nil, err
	}

	savedGame := g.(*Game)

	return savedGame, err
}
