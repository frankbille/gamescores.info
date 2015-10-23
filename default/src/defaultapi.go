package defaultapp

import (
	"appengine"
	gin "github.com/gamescores/gin"
	http "net/http"
)

const (
	gaeCtxKey = "GaeCtxKey"
)

type restService interface {
	CreateRoutes(parentRoute *gin.RouterGroup, rootRoute *gin.RouterGroup)
}

func init() {
	r := gin.New()

	root := r.Group("/")

	root.Use(gaeContext())
	root.Use(resolveUser())

	api := root.Group("/api")

	// Create list of services used
	services := []restService{
		createUserService(),
		createContextDefinitionService(),
		//		createPlayerService(),
		//		createLeagueService(),
		//		createGameService(),
		//		createAdminService(),
	}

	// Process the services
	for _, service := range services {
		service.CreateRoutes(api, root)
	}

	http.Handle("/", r)
}

func gaeContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := appengine.NewContext(c.Request)
		c.Set(gaeCtxKey, gaeCtx)
	}
}

func getGaeContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeCtxKey)
	return gc.(appengine.Context)
}
