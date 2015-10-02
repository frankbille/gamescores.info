package context

import (
	"appengine/datastore"
	"appengine/memcache"
	"fmt"
	gin "github.com/gamescores/gin"
	"math/rand"
	"time"
)

type adminService struct {
}

func createAdminService() adminService {
	return adminService{}
}

func (as adminService) CreateRoutes(parentRoute *gin.RouterGroup) {
	admin := parentRoute.Group("/admin")
	admin.GET("/sample", mustBeAdmin(), as.createSampleData)
}

func (as adminService) createSampleData(c *gin.Context) {
	deleteAll(entityGame, c)
	deleteAll(entityLeague, c)
	deleteAll(entityPlayer, c)
	memcache.Flush(getGaeContext(c))

	playerDao := createPlayerDao(c)
	leagueDao := createLeagueDao(c)
	gameDao := createGameDao(c)

	createdPlayerIds := addPlayers(playerDao, 100)
	createdLeagueIds := addLeagues(leagueDao, 20)
	addGames(gameDao, 2000, time.Now(), createdLeagueIds, createdPlayerIds)

	c.String(200, "OK")
}

func addGames(gameDao gameDao, numGames int, endDate time.Time, createdLeagueIds, createdPlayerIds []int64) {
	date := endDate.AddDate(0, 0, 0-numGames)
	for i := 0; i < numGames; i++ {
		addGame(gameDao, date, createdLeagueIds, createdPlayerIds)
		date = date.AddDate(0, 0, 1)
	}
}

func addGame(gameDao gameDao, gameDate time.Time, createdLeagueIds, createdPlayerIds []int64) {
	leagueID := getRandomID(createdLeagueIds)

	gamePlayers := make([]int64, 4)
	team1 := createTeam(createdPlayerIds, gamePlayers, true)
	team2 := createTeam(createdPlayerIds, gamePlayers, false)

	game := Game{
		GameDate: gameDate,
		Team1:    team1,
		Team2:    team2,
		LeagueID: leagueID,
	}

	gameDao.saveGame(game)
}

func createTeam(createdPlayerIds, gamePlayers []int64, teamWon bool) GameTeam {
	players := []int64{
		getUniqueRandomID(createdPlayerIds, gamePlayers),
		getUniqueRandomID(createdPlayerIds, gamePlayers),
	}

	score := 10
	if teamWon == false {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		score = r.Intn(9)
	}

	return GameTeam{
		Players: players,
		Score:   score,
	}
}

func getUniqueRandomID(createdIds, takenIds []int64) int64 {
	var id int64
	for id = getRandomID(createdIds); contains(takenIds, id); id = getRandomID(createdIds) {
	}
	takenIds[len(takenIds)-1] = id
	return id
}

func getRandomID(createdIds []int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(createdIds))
	return createdIds[index]
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func addLeagues(leagueDao leagueDao, numLeagues int) []int64 {
	createdLeagueIds := make([]int64, numLeagues)
	for i := 0; i < numLeagues; i++ {
		addLeague(leagueDao, fmt.Sprintf("League %d", i+1), i, createdLeagueIds)
	}
	return createdLeagueIds
}

func addLeague(leagueDao leagueDao, name string, index int, createdLeagueIds []int64) {
	savedLeague, _ := leagueDao.saveLeague(League{
		Name:        name,
		Description: "Created sample",
		Active:      true,
	})
	createdLeagueIds[index] = savedLeague.ID
}

func addPlayers(playerDao playerDao, numPlayers int) []int64 {
	createdPlayerIds := make([]int64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		addPlayer(playerDao, fmt.Sprintf("Player %d", i+1), i, createdPlayerIds)
	}
	return createdPlayerIds
}

func addPlayer(playerDao playerDao, name string, index int, createdPlayerIds []int64) {
	savedPlayer, _ := playerDao.savePlayer(Player{
		Name:   name,
		Active: true,
	})
	createdPlayerIds[index] = savedPlayer.ID
}

func deleteAll(kind string, c *gin.Context) {
	q := datastore.NewQuery(kind)
	q = q.KeysOnly()
	keys, _ := q.GetAll(getGaeContext(c), nil)
	datastore.DeleteMulti(getGaeContext(c), keys)
}
