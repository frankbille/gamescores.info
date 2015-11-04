package dao

import (
	"api/domain"
	"api/utils"
	"appengine/datastore"
	"github.com/gamescores/gin"
)

const (
	entityGameRating         = "GameRating"
	entityLeagueResult       = "LeagueResult"
	entityLeaguePlayerResult = "LeaguePlayerResult"
)

type RatingDao struct {
	dao
}

func CreateRatingDao(c *gin.Context) RatingDao {
	dao := createDao(utils.GetGaeContext(c))
	return RatingDao{dao}
}

func (dao *RatingDao) GetGameRatings(gameIds []int64) ([]domain.GameRating, error) {
	ratings := make([]domain.GameRating, len(gameIds))

	keys := make([]*datastore.Key, len(gameIds))

	for idx, gameId := range gameIds {
		keys[idx] = datastore.NewKey(dao.Context, entityGameRating, "", gameId, nil)
	}

	err := dao.getByIds(keys, ratings)

	return ratings, err
}

func (dao *RatingDao) SaveGameRating(gameRating domain.GameRating) (*domain.GameRating, error) {
	key := datastore.NewKey(dao.Context, entityGameRating, "", gameRating.GameID, nil)

	gr, err := dao.save(key, &gameRating)

	savedGameRating := gr.(*domain.GameRating)

	return savedGameRating, err
}

func (dao *RatingDao) GetLeagueResult(leagueId int64) (*domain.LeagueResult, error) {
	var leagueResult domain.LeagueResult

	key := datastore.NewKey(dao.Context, entityLeagueResult, "", leagueId, nil)
	err := dao.get(key, &leagueResult)

	if err != nil {
		return nil, err
	}

	var leaguePlayerResults []domain.LeaguePlayerResult

	_, err = dao.getListForAncestor(entityLeaguePlayerResult, 0, 0, key, []string{}, &leaguePlayerResults)

	if err != nil {
		return nil, err
	}

	leagueResult.PlayerResults = leaguePlayerResults

	return &leagueResult, nil
}

func (dao *RatingDao) SaveLeagueResult(leagueResult domain.LeagueResult) error {
	key := datastore.NewKey(dao.Context, entityLeagueResult, "", leagueResult.LeagueID, nil)

	_, err := dao.save(key, &leagueResult)

	if err != nil {
		return err
	}

	for _, leaguePlayerResult := range leagueResult.PlayerResults {
		playerKey := datastore.NewKey(dao.Context, entityLeaguePlayerResult, "", leaguePlayerResult.PlayerID, key)
		_, err = dao.save(playerKey, &leaguePlayerResult)

		if err != nil {
			return err
		}

	}

	return nil
}
