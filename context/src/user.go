package context

const (
	// Standard user. Can create contexts, and participate in them
	Standard UserRole = "standard"
	// Admin user. Links to a app engine developer account (app engine admin)
	// Used for cross context administration
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
	UserID   string   `json:"-"`
	LoggedIn bool     `json:"loggedIn" datastore:"-"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role" datastore:"-"`
}
