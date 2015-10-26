package service

import (
	"api/dao"
	"api/domain"
	"fmt"
	gin "github.com/gamescores/gin"
	"strconv"
)

type LeagueService struct {
}

func CreateLeagueService() LeagueService {
	return LeagueService{}
}

func (ls LeagueService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	leagues := parentRoute.Group("/leagues")
	leagues.GET("", ls.getLeagues)
	leagues.POST("", mustBeAuthenticated(), ls.createLeague)
	leagues.GET("/:leagueId", ls.getLeague)
	leagues.POST("/:leagueId", mustBeAuthenticated(), ls.updateLeague)
}

func (ls LeagueService) getLeagues(c *gin.Context) {
	var currentPage = getCurrentPage(c)
	var recordsPerPage = 50
	var start = getStartRecord(currentPage, recordsPerPage)

	leagueDao := dao.CreateLeagueDao(c)

	leagueArray, totalLeagueCount, err := leagueDao.GetLeagues(start, recordsPerPage)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if leagueArray == nil {
		leagueArray = []domain.League{}
	}

	for index := range leagueArray {
		addLeagueLinks(&leagueArray[index], c)
	}

	leagues := &domain.Leagues{
		Leagues: leagueArray,
	}

	addPaginationLinks(leagues, "/api/leagues", currentPage, recordsPerPage, totalLeagueCount)

	if isAuthenticated(c) {
		leagues.AddLink(domain.RelCreate, "/api/leagues")
	}

	c.JSON(200, leagues)
}

func (ls LeagueService) getLeague(c *gin.Context) {
	leagueID := getLeagueIDFromURL(c)

	if leagueID <= 0 {
		c.Redirect(304, "/api/leagues")
		return
	}

	leagueDao := dao.CreateLeagueDao(c)

	league, err := leagueDao.GetLeague(leagueID)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	addLeagueLinks(league, c)
	c.JSON(200, league)
}

func (ls LeagueService) createLeague(c *gin.Context) {
	var league domain.League

	c.Bind(&league)

	league.ID = 0
	league.Active = true

	ls.doSaveLeague(league, c)
}

func (ls LeagueService) updateLeague(c *gin.Context) {
	var league domain.League

	c.Bind(&league)

	ls.doSaveLeague(league, c)
}

func (ls LeagueService) doSaveLeague(league domain.League, c *gin.Context) {
	leagueDao := dao.CreateLeagueDao(c)

	savedLeague, err := leagueDao.SaveLeague(league)

	if err != nil {
		c.AbortWithError(500, err)
	}

	addLeagueLinks(savedLeague, c)
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

func addLeagueLinks(league *domain.League, c *gin.Context) {
	selfURL := fmt.Sprintf("/api/leagues/%d", league.ID)

	league.AddLink(domain.RelSelf, selfURL)

	if isAuthenticated(c) {
		league.AddLink(domain.RelUpdate, selfURL)
	}

	addLeagueGameLinks(league, c)
}
