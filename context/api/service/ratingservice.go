package service

import (
	"api/dao"
	"api/domain"
	"api/rating"
	"api/utils"
	"fmt"
	"github.com/gamescores/gin"
	"sort"
)

type RatingService struct {
	gameRatings         map[int64]map[int64]domain.GameRating
	leagueResults       map[int64]domain.LeagueResult
	latestPlayerRatings map[int64]map[int64]domain.PlayerRating
	defaultRatingType   rating.RatingType
}

func CreateRatingService() RatingService {
	return RatingService{
		gameRatings:         make(map[int64]map[int64]domain.GameRating),
		leagueResults:       make(map[int64]domain.LeagueResult),
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
	rs.gameRatings[leagueId] = make(map[int64]domain.GameRating)
	rs.latestPlayerRatings[leagueId] = make(map[int64]domain.PlayerRating)
	gameRatings := rs.gameRatings[leagueId]
	latestPlayerRatings := rs.latestPlayerRatings[leagueId]

	gameDao := dao.CreateGameDao(c)

	games, err := gameDao.GetGamesForRating(leagueId)

	if err != nil {
		return err
	}

	ratingProvider := rating.RatingProviderFactory(rs.defaultRatingType)
	ratingCalculator := ratingProvider.GetRatingCalculator()

	for _, game := range games {
		gameRating := ratingCalculator.CalculateRating(latestPlayerRatings, game)
		gameRatings[game.ID] = gameRating
	}

	return rs.recreateLeagueResult(c, leagueId)
}

func (rs *RatingService) recreateLeagueResult(c *gin.Context, leagueId int64) error {
	latestPlayerRatings := rs.latestPlayerRatings[leagueId]

	leaguePlayerIds := rs.convertMapInt64KeysToArray(latestPlayerRatings)

	playerDao := dao.CreatePlayerDao(c)
	players, err := playerDao.GetAllPlayersByID(leaguePlayerIds)
	playerMap := rs.convertPlayerArrayToPlayerMap(players)

	if err != nil {
		return err
	}

	leagueResult := domain.LeagueResult{
		LeagueID:      leagueId,
		PlayerResults: make([]domain.LeaguePlayerResult, len(leaguePlayerIds)),
	}

	for idx, playerId := range leaguePlayerIds {
		latestPlayerRating := latestPlayerRatings[playerId]
		player := playerMap[playerId]

		addPlayerLinks(&player, c)

		leagueResult.PlayerResults[idx] = domain.LeaguePlayerResult{
			Player: player,
			Rating: latestPlayerRating.Rating,
		}
	}

	sort.Sort(leagueResult.PlayerResults)

	for idx := 0; idx < len(leagueResult.PlayerResults); idx++ {
		leaguePlayerResult := &leagueResult.PlayerResults[idx]
		leaguePlayerResult.Position = idx + 1
	}

	leagueResult.AddLink(domain.RelSelf, fmt.Sprintf("/api/ratings/%d/result", leagueId))

	rs.leagueResults[leagueId] = leagueResult

	return nil
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

	ratings := rs.getRatings(leagueId, gameIds)

	c.JSON(200, ratings)
}

func (rs *RatingService) getRatings(leagueId int64, gameIds []int64) []domain.GameRating {
	ratings := make([]domain.GameRating, len(gameIds))

	for idx, gameId := range gameIds {
		ratings[idx] = rs.gameRatings[leagueId][gameId]
	}

	return ratings
}

func (rs *RatingService) getLeagueResultGinService(c *gin.Context) {
	leagueId := utils.ConvertToInt64(c.Param("leagueId"))

	leagueResult := rs.getLeagueResult(c, leagueId)

	c.JSON(200, leagueResult)
}

func (rs *RatingService) getLeagueResult(c *gin.Context, leagueId int64) domain.LeagueResult {
	leagueResult, found := rs.leagueResults[leagueId]

	if !found {
		rs.recalculateForLeague(c, leagueId)
		leagueResult = rs.leagueResults[leagueId]
	}

	return leagueResult
}
