package defaultapp

import (
	appengineuser "appengine/user"
	gin "github.com/gamescores/gin"
	"net/http"
)

const (
	userKey = "user"
)

type userService struct {
}

func createUserService() userService {
	return userService{}
}

func (us userService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	parentRoute.GET("/me", us.getCurrentUser)
}

func (us userService) getCurrentUser(c *gin.Context) {
	user := getCurrentUserFromGinContext(c)
	c.JSON(http.StatusOK, user)
}

func getCurrentUserFromGinContext(c *gin.Context) *User {
	usr := c.MustGet(userKey)
	return usr.(*User)
}

func resolveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := getGaeContext(c)
		currentGaeUser := appengineuser.Current(gaeCtx)

		var user *User
		if currentGaeUser != nil {
			dao := createDao(gaeCtx)
			userDao := userDao{dao}

			var err error
			user, err = userDao.getUserByID(currentGaeUser.ID)

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if user == nil {
				user = &User{
					UserID: currentGaeUser.ID,
					Email:  currentGaeUser.Email,
				}
				userDao.saveUser(user)
			}

			user.LoggedIn = true
			if currentGaeUser.Admin {
				user.Role = Admin
			} else {
				user.Role = Standard
			}

			logoutURL, _ := appengineuser.LogoutURL(gaeCtx, "")
			user.AddLink(relLogout, logoutURL)

			user.AddLink(relPrepare, "/api/context/prepare")
		} else {
			user = &User{
				LoggedIn: false,
			}

			loginURL, err := appengineuser.LoginURL(gaeCtx, "")

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			user.AddLink(relLogin, loginURL)
		}

		c.Set(userKey, user)
	}
}

func mustBeAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAuthenticated(c) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func mustBeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAdmin(c) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func isAuthenticated(c *gin.Context) bool {
	user := getCurrentUserFromGinContext(c)
	return user.LoggedIn
}

func isAdmin(c *gin.Context) bool {
	if isAuthenticated(c) {
		user := getCurrentUserFromGinContext(c)
		return user.Role == Admin
	}
	return false
}
