package context

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"strconv"
)

func getLeagues(c *gin.Context) {
	var currentPage = getCurrentPage(c)
	var recordsPerPage = 50
	var start = getStartRecord(currentPage, recordsPerPage)

	dao := createDao(getGaeContext(c))
	leagueDao := leagueDao{dao}

	leagueArray, totalLeagueCount, err := leagueDao.getLeagues(start, recordsPerPage)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if leagueArray == nil {
		leagueArray = []League{}
	}

	for index := range leagueArray {
		addLeagueLinks(&leagueArray[index])
	}

	leagues := &Leagues{
		Leagues: leagueArray,
	}

	addPaginationLinks(leagues, "/api/leagues", currentPage, recordsPerPage, totalLeagueCount)

	c.JSON(200, leagues)
}

func getLeague(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(304, "/api/leagues")
		return
	}

	dao := createDao(getGaeContext(c))
	leagueDao := leagueDao{dao}

	league, err := leagueDao.getLeague(leagueID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	addLeagueLinks(league)
	c.JSON(200, league)
}

func createLeague(c *gin.Context) {
	var league League

	c.Bind(&league)

	league.ID = 0
	league.Active = true

	doSaveLeague(league, c)
}

func updateLeague(c *gin.Context) {
	var league League

	c.Bind(&league)

	doSaveLeague(league, c)
}

func doSaveLeague(league League, c *gin.Context) {
	dao := createDao(getGaeContext(c))
	leagueDao := leagueDao{dao}

	savedLeague, err := leagueDao.saveLeague(league)

	if err != nil {
		c.AbortWithError(500, err)
	}

	addLeagueLinks(savedLeague)
	c.JSON(200, savedLeague)
}

// Private helper methods
func getLeagueIDFromURL(c *gin.Context) int64 {
	leagueIDString := c.Params.ByName("leagueId")
	leagueID, err := strconv.ParseInt(leagueIDString, 10, 64)
	if err != nil {
		return 0
	}
	return leagueID
}

func addLeagueLinks(league *League) {
	selfURL := fmt.Sprintf("/api/leagues/%d", league.ID)

	league.AddLink(relSelf, selfURL)
}
