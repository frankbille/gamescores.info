package context

import (
	appengineuser "appengine/user"
	gin "github.com/gamescores/gin"
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
	parentRoute.GET("/login", us.startLoginProcess)
}

func (us userService) getCurrentUser(c *gin.Context) {
	user := getCurrentUserFromGinContext(c)
	c.JSON(200, user)
}

func (us userService) startLoginProcess(c *gin.Context) {
	gaeCtx := getGaeRootContext(c)

	loginURL, err := appengineuser.LoginURL(gaeCtx, "")

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, loginURL)
}

func getCurrentUserFromGinContext(c *gin.Context) *User {
	usr := c.MustGet(userKey)
	return usr.(*User)
}

func resolveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := getGaeRootContext(c)
		currentGaeUser := appengineuser.Current(gaeCtx)

		var user *User
		if currentGaeUser != nil {
			dao := createDao(gaeCtx)
			userDao := userDao{dao}

			var err error
			user, err = userDao.getUserByID(currentGaeUser.ID)

			if err != nil {
				c.AbortWithError(500, err)
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
		} else {
			user = &User{
				LoggedIn: false,
			}

			user.AddLink(relLogin, "/api/login")
		}

		c.Set(userKey, user)
	}
}

func mustBeAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAuthenticated(c) {
			c.Next()
		} else {
			c.AbortWithStatus(401)
		}
	}
}

func mustBeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAdmin(c) {
			c.Next()
		} else {
			c.AbortWithStatus(401)
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
