package context

import (
	appengineuser "appengine/user"
	gin "github.com/gamescores/gin"
)

const (
	userKey = "user"
)

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
				c.AbortWithError(500, err)
				return
			}

			if user == nil {
				user = &User{
					UserID: currentGaeUser.ID,
					Email:  currentGaeUser.Email,
					Role:   Standard,
				}
				userDao.saveUser(user)
			}

			user.LoggedIn = true

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

func getCurrentUser(c *gin.Context) {
	user := getCurrentUserFromGinContext(c)
	c.JSON(200, user)
}

func startLoginProcess(c *gin.Context) {
	gaeCtx := getGaeContext(c)

	loginURL, err := appengineuser.LoginURL(gaeCtx, "")

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(304, loginURL)
}

func getCurrentUserFromGinContext(c *gin.Context) *User {
	usr := c.MustGet(userKey)
	return usr.(*User)
}
