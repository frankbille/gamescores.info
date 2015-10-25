package service

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"strconv"
	"src/domain"
	"src/dao"
	"src/utils"
)

const (
	relCreateGame domain.RelType = "creategame"
	relGames      domain.RelType = "games"
)

type gameService struct {
}

func CreateGameService() gameService {
	return gameService{}
}

func (gs gameService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	games := parentRoute.Group("/leagues/:leagueId/games")
	games.GET("", gs.getGames)
	games.POST("", mustBeAuthenticated(), gs.createGame)
	games.GET("/:gameId", gs.getGame)
	games.POST("/:gameId", mustBeAuthenticated(), gs.updateGame)
}

func (gs gameService) getGames(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(302, "/api/leagues")
		return
	}

	currentPage := getCurrentPage(c)
	recordsPerPage := 50
	start := getStartRecord(currentPage, recordsPerPage)

	gameDao := dao.CreateGameDao(c)

	gameArray, totalGameCount, err := gameDao.GetGames(start, recordsPerPage, leagueID)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error loading games: %v", err)
		c.AbortWithError(500, err)
		return
	}

	if gameArray == nil {
		gameArray = []domain.Game{}
	}

	for index := range gameArray {
		gs.addGameLinks(leagueID, &gameArray[index], c)
	}

	games := &domain.Games{
		Games: gameArray,
		Total: totalGameCount,
	}

	gs.addGamesLinks(games, leagueID, currentPage, recordsPerPage, totalGameCount, c)

	c.JSON(200, games)
}

func (gs gameService) getGame(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(302, "/api/leagues")
		return
	}

	gameID := getGameIDFromURL(c)

	if gameID <= 0 {
		c.Redirect(302, fmt.Sprintf("/api/leagues/%d/games", leagueID))
		return
	}

	gameDao := dao.CreateGameDao(c)

	game, err := gameDao.GetGame(leagueID, gameID)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error loading game: %v", err)
		c.AbortWithError(500, err)
		return
	}

	gs.addGameLinks(leagueID, game, c)
	c.JSON(200, game)
}

func (gs gameService) createGame(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	var game domain.Game

	c.Bind(&game)

	game.ID = 0
	game.LeagueID = leagueID

	gs.doSaveGame(game, c)
}

func (gs gameService) updateGame(c *gin.Context) {
	var game domain.Game

	c.Bind(&game)

	gs.doSaveGame(game, c)
}

func (gs gameService) doSaveGame(game domain.Game, c *gin.Context) {
	gameDao := dao.CreateGameDao(c)

	savedGame, err := gameDao.SaveGame(game)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error saving game: %v", err)
		c.AbortWithError(500, err)
	}

	gs.addGameLinks(game.LeagueID, savedGame, c)
	c.JSON(200, savedGame)
}

func (gs gameService) addGameLinks(leagueID int64, game *domain.Game, c *gin.Context) {
	gameURL := fmt.Sprintf("/api/leagues/%d/games/%d", leagueID, game.ID)

	game.AddLink(domain.RelSelf, gameURL)

	if isAuthenticated(c) {
		game.AddLink(domain.RelUpdate, gameURL)
	}
}

func (gs gameService) addGamesLinks(games *domain.Games, leagueID int64, currentPage, recordsPerPage, totalGameCount int, c *gin.Context) {
	gamesURL := fmt.Sprintf("/api/leagues/%d/games", leagueID)
	addPaginationLinks(games, gamesURL, currentPage, recordsPerPage, totalGameCount)
	if isAuthenticated(c) {
		games.AddLink(domain.RelCreate, gamesURL)
	}

	// Create a unique list of player id's from all the games returned
	playerIDSet := utils.NewInt64Set()
	for _, game := range games.Games {
		gs.addPlayerIdsFromGameTeam(playerIDSet, game.Team1)
		gs.addPlayerIdsFromGameTeam(playerIDSet, game.Team2)
	}
	addGetPlayerListByIDLinks(games, playerIDSet.Values(), c)
}

func (gs gameService) addPlayerIdsFromGameTeam(playerIDSet *utils.Int64Set, gameTeam domain.GameTeam) {
	for _, playerID := range gameTeam.Players {
		playerIDSet.Add(playerID)
	}
}

func addLeagueGameLinks(league *domain.League, c *gin.Context) {
	gamesURL := fmt.Sprintf("/api/leagues/%d/games", league.ID)
	league.AddLink(relGames, gamesURL)

	if isAuthenticated(c) {
		league.AddLink(relCreateGame, gamesURL)
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
