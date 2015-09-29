package context

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"strconv"
)

func getPlayers(c *gin.Context) {
	var currentPage = getCurrentPage(c)
	var recordsPerPage = 50
	var start = getStartRecord(currentPage, recordsPerPage)

	dao := createDao(getGaeContext(c))
	playerDao := playerDao{dao}

	playerArray, totalPlayerCount, err := playerDao.getPlayers(start, recordsPerPage)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if playerArray == nil {
		playerArray = []Player{}
	}

	for index := range playerArray {
		addPlayerLinks(&playerArray[index])
	}

	players := &Players{
		Players: playerArray,
	}

	addPaginationLinks(players, "/api/players", currentPage, recordsPerPage, totalPlayerCount)

	c.JSON(200, players)
}

func getPlayer(c *gin.Context) {
	playerID := getPlayerIDFromURL(c)

	if playerID <= 0 {
		c.Redirect(304, "/api/players")
		return
	}

	dao := createDao(getGaeContext(c))
	playerDao := playerDao{dao}

	player, err := playerDao.getPlayer(playerID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	addPlayerLinks(player)
	c.JSON(200, player)
}

func createPlayer(c *gin.Context) {
	var player Player

	c.Bind(&player)

	player.ID = 0
	player.Active = true

	doSavePlayer(player, c)
}

func updatePlayer(c *gin.Context) {
	var player Player

	c.Bind(&player)

	doSavePlayer(player, c)
}

func doSavePlayer(player Player, c *gin.Context) {
	dao := createDao(getGaeContext(c))
	playerDao := playerDao{dao}

	savedPlayer, err := playerDao.savePlayer(player)

	if err != nil {
		c.AbortWithError(500, err)
	}

	addPlayerLinks(savedPlayer)
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

func addPlayerLinks(player *Player) {
	selfURL := fmt.Sprintf("/api/players/%d", player.ID)

	player.AddLink(relSelf, selfURL)
}
