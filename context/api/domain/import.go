package domain

type ScoreBoardV1Import struct {
	DefaultHalResource
	DbDumpUrl string `json:"dbDumpUrl"`
}

type ScoreBoardV1ImportStatus struct {
	DefaultHalResource
	Importing           bool `json:"importing"`
	ImportedPlayerCount int  `json:"importedPlayerCount"`
	TotalPlayerCount    int  `json:"totalPlayerCount"`
	ImportedLeagueCount int  `json:"importedLeagueCount"`
	TotalLeagueCount    int  `json:"totalLeagueCount"`
	ImportedGameCount   int  `json:"importedGameCount"`
	TotalGameCount      int  `json:"totalGameCount"`
}