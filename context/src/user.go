package context

const (
	// Standard user. Can add and remove games
	Standard UserRole = "standard"
	// Admin user. Can set up leages, change settings and make give other users
	// admin rights.
	Admin UserRole = "admin"

	relLogin  RelType = "login"
	relLogout RelType = "logout"
)

// UserRole is the type for the roles a user can have
type UserRole string

// User is the authenticated user and it's settings.
type User struct {
	DefaultHalResource
	// The id from the appengine.user.User type
	UserID        string   `json:"-"`
	LoggedIn      bool     `json:"logged-in" datastore:"-"`
	Name          string   `json:"name,omitempty"`
	Email         string   `json:"-"`
	DefaultLeague int64    `json:"default-league,omitempty"`
	ClaimedPlayer int64    `json:"claimed-player,omitempty"`
	Role          UserRole `json:"role"`
}
