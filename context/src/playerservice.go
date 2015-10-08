package context

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"net/url"
	"strconv"
)

const (
	relPlayerList RelType = "playerlist"
)

type playerService struct {
}

func createPlayerService() playerService {
	return playerService{}
}

func (ps playerService) CreateRoutes(parentRoute *gin.RouterGroup) {
	players := parentRoute.Group("/players")
	players.GET("", ps.getPlayers)
	players.POST("", mustBeAuthenticated(), ps.createPlayer)
	players.GET("/:playerId", ps.getPlayer)
	players.POST("/:playerId", mustBeAuthenticated(), ps.updatePlayer)
}

func (ps playerService) getPlayers(c *gin.Context) {
	var currentPage = getCurrentPage(c)
	var recordsPerPage = 50
	var start = getStartRecord(currentPage, recordsPerPage)

	playerDao := createPlayerDao(c)

	playerArray, totalPlayerCount, err := playerDao.getPlayers(start, recordsPerPage)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if playerArray == nil {
		playerArray = []Player{}
	}

	for index := range playerArray {
		addPlayerLinks(&playerArray[index], c)
	}

	players := &Players{
		Players: playerArray,
	}

	addPaginationLinks(players, "/api/players", currentPage, recordsPerPage, totalPlayerCount)
	if isAuthenticated(c) {
		players.AddLink(relCreate, "/api/players")
	}

	c.JSON(200, players)
}

func (ps playerService) getPlayer(c *gin.Context) {
	playerID := getPlayerIDFromURL(c)

	if playerID <= 0 {
		c.Redirect(304, "/api/players")
		return
	}

	playerDao := createPlayerDao(c)

	player, err := playerDao.getPlayer(playerID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	addPlayerLinks(player, c)
	c.JSON(200, player)
}

func (ps playerService) createPlayer(c *gin.Context) {
	var player Player

	c.Bind(&player)

	player.ID = 0
	player.Active = true

	ps.doSavePlayer(player, c)
}

func (ps playerService) updatePlayer(c *gin.Context) {
	var player Player

	c.Bind(&player)

	ps.doSavePlayer(player, c)
}

func (ps playerService) doSavePlayer(player Player, c *gin.Context) {
	playerDao := createPlayerDao(c)

	savedPlayer, err := playerDao.savePlayer(player)

	if err != nil {
		c.AbortWithError(500, err)
	}

	addPlayerLinks(savedPlayer, c)
	c.JSON(200, savedPlayer)
}

// Private helper methods
func getPlayerIDFromURL(c *gin.Context) int64 {
	playerIDString := c.Params.ByName("playerId")
	playerID, err := strconv.ParseInt(playerIDString, 10, 64)
	if err != nil {
		return 0
	}
	return playerID
}

func addPlayerLinks(player *Player, c *gin.Context) {
	selfURL := fmt.Sprintf("/api/players/%d", player.ID)

	player.AddLink(relSelf, selfURL)

	if isAuthenticated(c) {
		player.AddLink(relUpdate, selfURL)
	}
}

func addGetPlayerListByIDLinks(games *Games, playerIds []int64, c *gin.Context) {
	playerListURL, err := url.Parse("/api/players")

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	q := playerListURL.Query()
	for _, playerID := range playerIds {

		getGaeContext(c).Infof("PlayerID: %v", fmt.Sprintf("%d", playerID))
		q.Add("id", fmt.Sprintf("%d", playerID))
	}
	playerListURL.RawQuery = q.Encode()

	games.AddLink(relPlayerList, playerListURL.String())
}
