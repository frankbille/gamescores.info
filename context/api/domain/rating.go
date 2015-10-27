package domain

type GameRating struct {
	GameID            int64      `json:"gameId"`
	WinningTeamRating TeamRating `json:"winningTeamRating`
	LoosingTeamRating TeamRating `json:"lossingTeamRating`
}

type TeamRating struct {
	Rating        float64        `json:"rating"`
	PlayerRatings []PlayerRating `json:"playerRatings"`
}

type PlayerRating struct {
	PlayerID int64   `json:"playerId"`
	Rating   float64 `json:"rating"`
}

type LeagueResult struct {
	LeagueID     int64                `json:"leagueId"`
	PlayerResult []LeaguePlayerResult `json:"players"`
}

type LeaguePlayerResult struct {
	PlayerID   int64   `json:"playerId"`
	PlayerName string  `json:"playerName"`
	Position   int     `json:"position"`
	Rating     float64 `json:"rating"`
}
