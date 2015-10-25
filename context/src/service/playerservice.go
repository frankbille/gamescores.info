package service

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"net/url"
	"strconv"
	"src/domain"
	"src/dao"
	"src/utils"
)

const (
	relPlayerList domain.RelType = "playerlist"
)

type playerService struct {
}

func CreatePlayerService() playerService {
	return playerService{}
}

func (ps playerService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	players := parentRoute.Group("/players")
	players.GET("", ps.getPlayers)
	players.POST("", mustBeAuthenticated(), ps.createPlayer)
	players.GET("/:playerId", ps.getPlayer)
	players.POST("/:playerId", mustBeAuthenticated(), ps.updatePlayer)
}

func (ps playerService) getPlayers(c *gin.Context) {
	idList := c.Request.URL.Query()["id"]

	if len(idList) == 0 {
		ps.handleGetAllPlayers(c)
	} else {
		ps.handleGetSpecificPlayers(c, idList)
	}
}

func (ps playerService) handleGetAllPlayers(c *gin.Context) {
	var currentPage = getCurrentPage(c)
	var recordsPerPage = 50
	var start = getStartRecord(currentPage, recordsPerPage)

	playerDao := dao.CreatePlayerDao(c)

	playerArray, totalPlayerCount, err := playerDao.GetPlayers(start, recordsPerPage)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error loading players: %v", err)
		c.AbortWithError(500, err)
		return
	}

	if playerArray == nil {
		playerArray = []domain.Player{}
	}

	for index := range playerArray {
		addPlayerLinks(&playerArray[index], c)
	}

	players := &domain.Players{
		Players: playerArray,
	}

	addPaginationLinks(players, "/api/players", currentPage, recordsPerPage, totalPlayerCount)
	if isAuthenticated(c) {
		players.AddLink(domain.RelCreate, "/api/players")
	}

	c.JSON(200, players)
}

func (ps playerService) handleGetSpecificPlayers(c *gin.Context, idList []string) {
	playerDao := dao.CreatePlayerDao(c)

	playerIds := make([]int64, len(idList))
	idx := 0
	for _, id := range idList {
		playerID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.GetGaeContext(c).Errorf("Not an integer: %v", err)
			c.AbortWithError(500, err)
			return
		}
		playerIds[idx] = playerID
		idx++
	}

	playerArray, err := playerDao.GetAllPlayersByID(playerIds)
	if err != nil {
		utils.GetGaeContext(c).Errorf("Error getting players by id: %v", err)
		c.AbortWithError(500, err)
		return
	}

	for index := range playerArray {
		addPlayerLinks(&playerArray[index], c)
	}

	players := &domain.Players{
		Players: playerArray,
	}

	c.JSON(200, players)
}

func (ps playerService) getPlayer(c *gin.Context) {
	playerID := getPlayerIDFromURL(c)

	if playerID <= 0 {
		c.Redirect(304, "/api/players")
		return
	}

	playerDao := dao.CreatePlayerDao(c)

	player, err := playerDao.GetPlayer(playerID)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error loading player: %v", err)
		c.AbortWithError(500, err)
		return
	}

	addPlayerLinks(player, c)
	c.JSON(200, player)
}

func (ps playerService) createPlayer(c *gin.Context) {
	var player domain.Player

	c.Bind(&player)

	player.ID = 0
	player.Active = true

	ps.doSavePlayer(player, c)
}

func (ps playerService) updatePlayer(c *gin.Context) {
	var player domain.Player

	c.Bind(&player)

	ps.doSavePlayer(player, c)
}

func (ps playerService) doSavePlayer(player domain.Player, c *gin.Context) {
	playerDao := dao.CreatePlayerDao(c)

	savedPlayer, err := playerDao.SavePlayer(player)

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error saving player: %v", err)
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

func addPlayerLinks(player *domain.Player, c *gin.Context) {
	selfURL := fmt.Sprintf("/api/players/%d", player.ID)

	player.AddLink(domain.RelSelf, selfURL)

	if isAuthenticated(c) {
		player.AddLink(domain.RelUpdate, selfURL)
	}
}

func addGetPlayerListByIDLinks(games *domain.Games, playerIds []int64, c *gin.Context) {
	playerListURL, err := url.Parse("/api/players")

	if err != nil {
		utils.GetGaeContext(c).Errorf("Error parsing URL: %v", err)
		c.AbortWithError(500, err)
		return
	}

	q := playerListURL.Query()
	for _, playerID := range playerIds {
		q.Add("id", fmt.Sprintf("%d", playerID))
	}
	playerListURL.RawQuery = q.Encode()

	games.AddLink(relPlayerList, playerListURL.String())
}
