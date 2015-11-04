package service

import (
	"api/dao"
	"api/domain"
	"api/rating"
	"api/utils"
	"fmt"
	"github.com/gamescores/gin"
	"net/url"
	"sort"
)

const (
	relGameRatingList domain.RelType = "gameratinglist"
)

type RatingService struct {
	latestPlayerRatings map[int64]map[int64]domain.PlayerRating
	defaultRatingType   rating.RatingType
}

func CreateRatingService() RatingService {
	return RatingService{
		latestPlayerRatings: make(map[int64]map[int64]domain.PlayerRating),
		defaultRatingType:   rating.RATING_ELO,
	}
}

func (rs RatingService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	ratingsRoute := parentRoute.Group("/ratings")

	ratingsRoute.GET("/:leagueId/games", rs.getRatingsGinService)
	ratingsRoute.GET("/:leagueId/result", rs.getLeagueResultGinService)
	ratingsRoute.GET("/:leagueId/recalc", rs.recalculateForLeagueGinService)
}

func (rs *RatingService) SaveRating(c *gin.Context, game *domain.Game) error {
	return rs.recalculateForLeague(c, game.LeagueID)
}

func (rs *RatingService) recalculateForLeagueGinService(c *gin.Context) {
	leagueId := utils.ConvertToInt64(c.Param("leagueId"))

	err := rs.recalculateForLeague(c, leagueId)

	if err != nil {
		utils.AbortWithError(c, err)
		return
	}

	c.String(200, "OK")
}

func (rs *RatingService) recalculateForLeague(c *gin.Context, leagueId int64) error {
	rs.latestPlayerRatings[leagueId] = make(map[int64]domain.PlayerRating)
	latestPlayerRatings := rs.latestPlayerRatings[leagueId]

	gameDao := dao.CreateGameDao(c)

	games, err := gameDao.GetGamesForRating(leagueId)

	if err != nil {
		return err
	}

	ratingDao := dao.CreateRatingDao(c)

	ratingProvider := rating.RatingProviderFactory(rs.defaultRatingType)
	ratingCalculator := ratingProvider.GetRatingCalculator()

	for _, game := range games {
		gameRating := ratingCalculator.CalculateRating(latestPlayerRatings, game)
		ratingDao.SaveGameRating(gameRating)
	}

	return rs.recreateLeagueResult(c, leagueId)
}

func (rs *RatingService) recreateLeagueResult(c *gin.Context, leagueId int64) error {
	latestPlayerRatings := rs.latestPlayerRatings[leagueId]

	leaguePlayerIds := rs.convertMapInt64KeysToArray(latestPlayerRatings)

	leagueResult := domain.LeagueResult{
		LeagueID:      leagueId,
		PlayerResults: []domain.LeaguePlayerResult{},
	}

	for _, playerId := range leaguePlayerIds {
		latestPlayerRating := latestPlayerRatings[playerId]

		leagueResult.PlayerResults = append(leagueResult.PlayerResults, domain.LeaguePlayerResult{
			PlayerID: playerId,
			Rating:   latestPlayerRating.NewRating,
		})
	}

	sort.Sort(leagueResult.PlayerResults)

	for idx := 0; idx < len(leagueResult.PlayerResults); idx++ {
		leaguePlayerResult := &leagueResult.PlayerResults[idx]
		leaguePlayerResult.Position = idx + 1
	}

	ratingDao := dao.CreateRatingDao(c)
	return ratingDao.SaveLeagueResult(leagueResult)
}

func (rs *RatingService) convertMapInt64KeysToArray(dataMap map[int64]domain.PlayerRating) []int64 {
	v := make([]int64, len(dataMap))
	idx := 0
	for i := range dataMap {
		v[idx] = i
		idx++
	}
	return v
}

func (rs *RatingService) convertPlayerArrayToPlayerMap(players []domain.Player) map[int64]domain.Player {
	playerMap := make(map[int64]domain.Player)
	for _, player := range players {
		playerMap[player.ID] = player
	}
	return playerMap
}

func (rs *RatingService) getRatingsGinService(c *gin.Context) {
	leagueId := utils.ConvertToInt64(c.Param("leagueId"))
	gameIdStrings := c.Request.URL.Query()["gameId"]

	gameIds := make([]int64, len(gameIdStrings))
	for idx, gameIdString := range gameIdStrings {
		gameIds[idx] = utils.ConvertToInt64(gameIdString)
	}

	ratings, err := rs.getRatings(c, leagueId, gameIds)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error getting ratings: %v", err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, ratings)
}

func (rs *RatingService) getRatings(c *gin.Context, leagueId int64, gameIds []int64) ([]domain.GameRating, error) {
	ratingDao := dao.CreateRatingDao(c)

	return ratingDao.GetGameRatings(gameIds)
}

func (rs *RatingService) getLeagueResultGinService(c *gin.Context) {
	leagueId := utils.ConvertToInt64(c.Param("leagueId"))

	leagueResult, err := rs.getLeagueResult(c, leagueId)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error getting league result: %v", err)
		c.AbortWithError(500, err)
		return
	}

	leagueResult.AddLink(domain.RelSelf, fmt.Sprintf("/api/ratings/%d/result", leagueId))

	c.JSON(200, leagueResult)
}

func (rs *RatingService) getLeagueResult(c *gin.Context, leagueId int64) (*domain.LeagueResult, error) {
	ratingDao := dao.CreateRatingDao(c)

	leagueResult, err := ratingDao.GetLeagueResult(leagueId)

	if err != nil {
		return nil, err
	}

	playerDao := dao.CreatePlayerDao(c)

	leaguePlayerIds := make([]int64, len(leagueResult.PlayerResults))
	for idx, leaguePlayerResult := range leagueResult.PlayerResults {
		leaguePlayerIds[idx] = leaguePlayerResult.PlayerID
	}

	players, err := playerDao.GetAllPlayersByID(leaguePlayerIds)
	playerMap := rs.convertPlayerArrayToPlayerMap(players)

	utils.GetGaeContext(c).Infof("%#v", playerMap)

	for idx, _ := range leagueResult.PlayerResults {
		leaguePlayerResult := &leagueResult.PlayerResults[idx]
		leaguePlayerResult.Player = playerMap[leaguePlayerResult.PlayerID]
	}

	return leagueResult, nil
}

func addGetGameRatingsByIDLinks(games *domain.Games, leagueId int64, c *gin.Context) {
	gameRatingsListURL, err := url.Parse(fmt.Sprintf("/api/ratings/%d/games", leagueId))

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error parsing URL: %v", err)
		c.AbortWithError(500, err)
		return
	}

	q := gameRatingsListURL.Query()
	for _, game := range games.Games {
		q.Add("gameId", fmt.Sprintf("%d", game.ID))
	}
	gameRatingsListURL.RawQuery = q.Encode()

	games.AddLink(relGameRatingList, gameRatingsListURL.String())
}
