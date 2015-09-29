package context

// League describes a grouping of games.
type League struct {
	DefaultHalResource
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// Leagues contains a list of leagues, including pagination links.
type Leagues struct {
	DefaultHalResource
	Leagues []League `json:"leagues"`
}
