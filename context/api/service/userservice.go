package service

import (
	"api/dao"
	"api/domain"
	"api/utils"
	"appengine/datastore"
	appengineuser "appengine/user"
	gin "github.com/gamescores/gin"
)

const (
	userKey = "user"
)

type UserService struct {
}

func CreateUserService() UserService {
	return UserService{}
}

func (us UserService) CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup) {
	parentRoute.GET("/me", us.getCurrentUser)
	parentRoute.GET("/login", us.startLoginProcess)
}

func (us UserService) getCurrentUser(c *gin.Context) {
	user := getCurrentUserFromGinContext(c)
	c.JSON(200, user)
}

func (us UserService) startLoginProcess(c *gin.Context) {
	gaeCtx := utils.GetGaeRootContext(c)

	loginURL, err := appengineuser.LoginURL(gaeCtx, "")

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, loginURL)
}

func getCurrentUserFromGinContext(c *gin.Context) *domain.User {
	usr := c.MustGet(userKey)
	return usr.(*domain.User)
}

func ResolveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := utils.GetGaeRootContext(c)
		currentGaeUser := appengineuser.Current(gaeCtx)

		var user *domain.User
		if currentGaeUser != nil {
			userDao := dao.CreateUserDao(c)

			var err error
			user, err = userDao.GetUserByID(currentGaeUser.ID)

			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			if user == nil {
				user = &domain.User{
					UserID: currentGaeUser.ID,
					Email:  currentGaeUser.Email,
				}
				userDao.SaveUser(user)
			}

			user.LoggedIn = true
			contextDefinition := GetGameContext(c)

			userKey := datastore.NewKey(gaeCtx, dao.EntityUser, user.UserID, 0, nil)

			if contextDefinition.IsUserOwner(userKey) {
				user.Role = domain.Admin
			} else {
				user.Role = domain.Standard
			}

			logoutURL, _ := appengineuser.LogoutURL(gaeCtx, "")
			user.AddLink(domain.RelLogout, logoutURL)
		} else {
			user = &domain.User{
				LoggedIn: false,
			}

			user.AddLink(domain.RelLogin, "/api/login")
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
		return user.Role == domain.Admin
	}
	return false
}
