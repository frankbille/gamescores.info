package service

import (
	"api/domain"
	"api/rating"
	"api/utils"
	"github.com/gamescores/gin"
)

type RatingService struct {
	gameRatings             map[int64]domain.GameRating
	leagueResults           map[int64]domain.LeagueResult
	defaultRatingCalculator rating.RatingCalculator
}

func CreateRatingService() RatingService {
	return RatingService{
		gameRatings:             make(map[int64]domain.GameRating),
		leagueResults:           make(map[int64]domain.LeagueResult),
		defaultRatingCalculator: rating.CreateEloRatingCalculator(),
	}
}

func (rs RatingService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	ratingsRoute := parentRoute.Group("/ratings")

	ratingsRoute.GET("/forGames", rs.getRatingsGinService)
	ratingsRoute.GET("/leagueResult/:leagueId", rs.getLeagueResultGinService)
}

func (rs *RatingService) SaveRating(game *domain.Game) {

}

func (rs *RatingService) getRatingsGinService(c *gin.Context) {
	gameIdStrings := c.Request.URL.Query()["gameId"]

	gameIds := make([]int64, len(gameIdStrings))
	for idx, gameIdString := range gameIdStrings {
		gameIds[idx] = utils.ConvertToInt64(gameIdString)
	}

	ratings := rs.getRatings(gameIds)

	c.JSON(200, ratings)
}

func (rs *RatingService) getRatings(gameIds []int64) []domain.GameRating {
	ratings := make([]domain.GameRating, len(gameIds))

	for idx, gameId := range gameIds {
		ratings[idx] = rs.gameRatings[gameId]
	}

	return ratings
}

func (rs *RatingService) getLeagueResultGinService(c *gin.Context) {
	leagueId := utils.ConvertToInt64(c.Param("leagueId"))

	leagueResult := rs.getLeagueResult(leagueId)

	c.JSON(200, leagueResult)
}

func (rs *RatingService) getLeagueResult(leagueId int64) domain.LeagueResult {
	return rs.leagueResults[leagueId]
}
