package context

import (
	datastore "appengine/datastore"
)

const entityLeague string = "League"

type leagueDao struct {
	dao
}

func (dao *leagueDao) getLeagues(start, limit int) ([]League, int, error) {
	var leagues []League

	count, err := dao.getList(entityLeague, start, limit, &leagues)

	return leagues, count, err
}

func (dao *leagueDao) getLeague(leagueID int64) (*League, error) {
	var league League
	key := datastore.NewKey(dao.Context, entityLeague, "", leagueID, nil)

	err := dao.get(key, &league)

	return &league, err
}

func (dao *leagueDao) saveLeague(league League) (*League, error) {
	if league.ID == 0 {
		leagueID, _, _ := datastore.AllocateIDs(dao.Context, entityLeague, nil, 1)
		league.ID = leagueID
	}

	key := datastore.NewKey(dao.Context, entityLeague, "", league.ID, nil)

	l, err := dao.save(key, &league)

	if err != nil {
		return nil, err
	}

	savedLeague := l.(*League)

	return savedLeague, err
}
