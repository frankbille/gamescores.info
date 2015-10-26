package domain

// Player describes the individual players. They are unique per context but
// shared across all leagues and games in that context.
type Player struct {
	DefaultHalResource
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// Players contains a list of players, including pagination links.
type Players struct {
	DefaultHalResource
	Players []Player `json:"players"`
}
