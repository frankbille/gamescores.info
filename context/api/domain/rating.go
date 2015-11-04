package domain

type GameRating struct {
	GameID            int64      `json:"gameId"`
	WinningTeamRating TeamRating `json:"winningTeamRating"`
	LoosingTeamRating TeamRating `json:"loosingTeamRating"`
}

type TeamRating struct {
	Rating        float64        `json:"rating"`
	PlayerRatings []PlayerRating `json:"playerRatings"`
}

type PlayerRating struct {
	PlayerID  int64   `json:"playerId"`
	OldRating float64 `json:"oldrating"`
	NewRating float64 `json:"newrating"`
}

type LeagueResult struct {
	DefaultHalResource
	LeagueID      int64               `json:"leagueId"`
	PlayerResults LeaguePlayerResults `json:"players"`
}

type LeaguePlayerResults []LeaguePlayerResult

func (lr LeaguePlayerResults) Len() int {
	return len(lr)
}

func (lr LeaguePlayerResults) Less(i, j int) bool {
	return lr[i].Rating > lr[j].Rating
}

func (lr LeaguePlayerResults) Swap(i, j int) {
	lr[i], lr[j] = lr[j], lr[i]
}

type LeaguePlayerResult struct {
	Player   Player  `json:"player"`
	Position int     `json:"position"`
	Rating   float64 `json:"rating"`
}
