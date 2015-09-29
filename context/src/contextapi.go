package context

import (
	"appengine"
	gin "github.com/gamescores/gin"
	http "net/http"
	"os"
	"strings"
)

const (
	gaeCtxKey = "GaeCtxKey"
)

func init() {
	r := gin.New()

	r.Use(gaeContext())
	r.Use(resolveUser())

	api := r.Group("/api")

	// User endpoints
	api.GET("/me", getCurrentUser)
	api.GET("/login", startLoginProcess)

	// Player endpoints
	players := api.Group("/players")
	players.GET("", getPlayers)
	players.POST("", createPlayer)
	players.GET("/:playerId", getPlayer)
	players.POST("/:playerId", updatePlayer)

	http.Handle("/", r)
}

func gaeContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := appengine.NewContext(c.Request)

		namespace := ""

		if productionDomain := os.Getenv("PRODUCTION_DOMAIN"); productionDomain != "" {
			lastIndex := strings.LastIndex(c.Request.Host, productionDomain)

			if lastIndex > -1 {
				namespace = strings.Replace(c.Request.Host, productionDomain, "", lastIndex)
			}
		}

		if namespace != "" {
			nameSpacedGaeCtx, err := appengine.Namespace(gaeCtx, namespace)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			gaeCtx = nameSpacedGaeCtx
		}

		c.Set(gaeCtxKey, gaeCtx)
	}
}

func getGaeContext(c *gin.Context) appengine.Context {
	gc := c.MustGet(gaeCtxKey)
	return gc.(appengine.Context)
}
