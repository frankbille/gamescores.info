package context

import (
	"appengine"
	gin "github.com/gamescores/gin"
	http "net/http"
	"os"
	"strings"
	"fmt"
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
		createPlayerService(),
		createLeagueService(),
		createGameService(),
		createAdminService(),
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

		namespace := ""

		if productionDomain := os.Getenv("PRODUCTION_DOMAIN"); productionDomain != "" {
			if strings.HasPrefix(productionDomain, ".") == false {
				productionDomain = fmt.Sprintf(".%s", productionDomain);
			}

			lastIndex := strings.LastIndex(c.Request.Host, productionDomain)

			if lastIndex > -1 {
				namespace = strings.Replace(c.Request.Host, productionDomain, "", lastIndex)
			}
		} else if devNamespace := os.Getenv("DEV_NAMESPACE"); devNamespace != "" {
			namespace = devNamespace
		}

		if namespace != "" {
			gaeCtx.Debugf("Using namespace: \"%s\"", namespace)
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
