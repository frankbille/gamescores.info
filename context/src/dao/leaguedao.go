package dao

import (
	datastore "appengine/datastore"
	gin "github.com/gamescores/gin"
	"src/domain"
	"src/utils"
)

const EntityLeague string = "League"

type LeagueDao struct {
	dao
}

func CreateLeagueDao(c *gin.Context) LeagueDao {
	dao := createDao(utils.GetGaeContext(c))
	return LeagueDao{dao}
}

func (dao *LeagueDao) GetLeagues(start, limit int) ([]domain.League, int, error) {
	var leagues []domain.League

	count, err := dao.getList(EntityLeague, start, limit, &leagues)

	return leagues, count, err
}

func (dao *LeagueDao) GetLeague(leagueID int64) (*domain.League, error) {
	var league domain.League
	key := datastore.NewKey(dao.Context, EntityLeague, "", leagueID, nil)

	err := dao.get(key, &league)

	return &league, err
}

func (dao *LeagueDao) SaveLeague(league domain.League) (*domain.League, error) {
	if league.ID == 0 {
		leagueID, _, _ := datastore.AllocateIDs(dao.Context, EntityLeague, nil, 1)
		league.ID = leagueID
	}

	key := datastore.NewKey(dao.Context, EntityLeague, "", league.ID, nil)

	l, err := dao.save(key, &league)

	if err != nil {
		return nil, err
	}

	savedLeague := l.(*domain.League)

	return savedLeague, err
}
