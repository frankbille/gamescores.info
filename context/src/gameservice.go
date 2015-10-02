package context

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"strconv"
)

const (
	relCreateGame RelType = "create-game"
)

type gameService struct {
}

func createGameService() gameService {
	return gameService{}
}

func (gs gameService) CreateRoutes(parentRoute *gin.RouterGroup) {
	games := parentRoute.Group("/leagues/:leagueId/games")
	games.GET("", gs.getGames)
	games.POST("", mustBeAuthenticated(), gs.createGame)
	games.GET("/:gameId", gs.getGame)
	games.POST("/:gameId", mustBeAuthenticated(), gs.updateGame)
}

func (gs gameService) getGames(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(304, "/api/leagues")
		return
	}

	currentPage := getCurrentPage(c)
	recordsPerPage := 100
	start := getStartRecord(currentPage, recordsPerPage)

	gameDao := createGameDao(c)

	gameArray, totalGameCount, err := gameDao.getGames(start, recordsPerPage, leagueID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if gameArray == nil {
		gameArray = []Game{}
	}

	for index := range gameArray {
		gs.addGameLinks(leagueID, &gameArray[index], c)
	}

	games := &Games{
		Games: gameArray,
		Total: totalGameCount,
	}

	gamesURL := fmt.Sprintf("/api/leagues/%d/games", leagueID)
	addPaginationLinks(games, gamesURL, currentPage, recordsPerPage, totalGameCount)
	if isAuthenticated(c) {
		games.AddLink(relCreate, gamesURL)
	}

	c.JSON(200, games)
}

func (gs gameService) getGame(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(304, "/api/leagues")
		return
	}

	gameID := getGameIDFromURL(c)

	if gameID <= 0 {
		c.Redirect(304, fmt.Sprintf("/api/leagues/%d/games", leagueID))
		return
	}

	gameDao := createGameDao(c)

	game, err := gameDao.getGame(leagueID, gameID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	gs.addGameLinks(leagueID, game, c)
	c.JSON(200, game)
}

func (gs gameService) createGame(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	var game Game

	c.Bind(&game)

	game.ID = 0
	game.LeagueID = leagueID

	gs.doSaveGame(game, c)
}

func (gs gameService) updateGame(c *gin.Context) {
	var game Game

	c.Bind(&game)

	gs.doSaveGame(game, c)
}

func (gs gameService) doSaveGame(game Game, c *gin.Context) {
	gameDao := createGameDao(c)

	savedGame, err := gameDao.saveGame(game)

	if err != nil {
		c.AbortWithError(500, err)
	}

	gs.addGameLinks(game.LeagueID, savedGame, c)
	c.JSON(200, savedGame)
}

func (gs gameService) addGameLinks(leagueID int64, game *Game, c *gin.Context) {
	gameURL := fmt.Sprintf("/api/leagues/%d/games/%d", leagueID, game.ID)

	game.AddLink(relSelf, gameURL)

	if isAuthenticated(c) {
		game.AddLink(relUpdate, gameURL)
	}
}

func getGameIDFromURL(c *gin.Context) int64 {
	gameIDString := c.Params.ByName("gameId")
	gameID, err := strconv.ParseInt(gameIDString, 10, 64)
	if err != nil {
		return 0
	}
	return gameID
}
