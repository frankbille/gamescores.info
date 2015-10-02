package context

import (
	"appengine/datastore"
	"time"
)

// Game is an individual game, played in a league with 2 teams each with at least
// one player
type Game struct {
	DefaultHalResource
	ID       int64          `json:"id"`
	GameDate time.Time      `json:"game-date"`
	Team1    GameTeam       `json:"team1"`
	Team2    GameTeam       `json:"team2"`
	League   *datastore.Key `json:"-"`
	LeagueID int64          `json:"league-id"`
}

// GameTeam represents one of the teams that played and the score they got
type GameTeam struct {
	Players []int64 `json:"players"`
	Score   int     `json:"score"`
}

// Games holds a list of games as well as hyperlinks
type Games struct {
	DefaultHalResource
	Games []Game `json:"games"`
	Total int    `json:"total"`
}
